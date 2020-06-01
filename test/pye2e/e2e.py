import time
import unittest
import subprocess
import os
import uuid
import tempfile
import re
from typing import List

path = os.path.dirname(os.path.realpath(__file__))

kubectl_pattern = re.compile(r"(?P<nodename>[a-zA-Z0-9\-]+)\s+(?P<status>\w+)\s+(?:\w+)\s+(?P<uptime>[\w\d]+)\s+(?:\w+)")
describe_pattern = re.compile(r"Capacity:.*?\n\s+cpu:\s+(?P<cpu>.*?)\n\s+ephemeral-storage:\s+(?P<eph>.*?)\n\s+memory:\s+(?P<memory>.*?)\n\s+pods:.*?\n\s+storage:\s+(?P<storage>.*)\n")


class E2ETest(unittest.TestCase):
    def call_apate_cli(self, *args: str, stdin: str = "") -> str:
        read, write = os.pipe()
        os.write(write, stdin.encode())
        os.close(write)

        output = subprocess.run(["go", "run", f"{path}/../../cmd/apate/main.go", *args], stdout=subprocess.PIPE, stdin=read)

        self.assertEqual(output.returncode, 0)
        return output.stdout.decode("utf-8")

    def call_kubectl(self, *args: str, autokubeconfig=True) -> str:
        command = ["kubectl"]
        if autokubeconfig:
            kubeconfig = self.call_apate_cli("kubeconfig")
            filename = tempfile.gettempdir() + "/apate-e2e-kubeconfig" + str(uuid.uuid4())
            with open(filename, "w") as f:
                f.write(kubeconfig)
            command.extend(["--kubeconfig", filename])

        output = subprocess.run([*command, *args], stdout=subprocess.PIPE)

        self.assertEqual(output.returncode, 0)
        return output.stdout.decode("utf-8")

    def stop_apate(self):
        print("stopping apate-cp container")
        subprocess.run("docker stop apate-cp".split())

        while "apate-cp" in subprocess.run(["docker", "ps"], stdout=subprocess.PIPE).stdout.decode():
            print("waiting for apate-cp to stop")

    def tearDown(self) -> None:
        self.stop_apate()

    def setUp(self) -> None:
        self.stop_apate()

        print("starting control plane")
        self.call_apate_cli("create")

    def kubectl_create(self, config):
        filename = tempfile.gettempdir() + "/apate-e2e-config" + str(uuid.uuid4())
        with open(filename, "w") as f:
            f.write(config)

        self.call_kubectl("create", "-f", filename)

    def test_node(self):
        nodes = self.call_kubectl(*"get nodes".split())
        matches = kubectl_pattern.findall(nodes.split("\n", 1)[1])

        # Test if the apate control plane was started
        self.assertEqual(len(matches), 1)
        self.assertEqual(matches[0][0], "apate-control-plane")

        self.kubectl_create("""
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: e2e-deployment
spec:
    replicas: 2
    resources:
        memory: 5G
        cpu: 10
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
        """)

        nodes = self.call_kubectl(*"get nodes".split())
        matches: List[str] = kubectl_pattern.findall(nodes.split("\n", 1)[1])

        # after spawning 2 apatelets, there should be 3 nodes in total
        self.assertEqual(len(matches), 3)

        # assert two of these 3 nodes are apatelets
        self.assertEqual(len(
            [i for i in matches if i[0].startswith("apatelet-")]
        ), 2)

    def test_node_failure(self):
        nodes = self.call_kubectl(*"get nodes".split())
        print(nodes)
        matches = kubectl_pattern.findall(nodes.split("\n", 1)[1])

        # Test if the apate control plane was started
        self.assertEqual(len(matches), 1)
        self.assertEqual(matches[0][0], "apate-control-plane")

        self.kubectl_create("""
    apiVersion: apate.opendc.org/v1
    kind: NodeConfiguration
    metadata:
        name: test-deployment1
    spec:
        replicas: 1
        resources:
            memory: 5G
            cpu: 1000
            storage: 5T
            ephemeral_storage: 120G
            max_pods: 150
        tasks:
            - timestamp: 10s
              state:
                  node_failed: true
            """)

        for i in range(10):
            nodes = self.call_kubectl(*"get nodes".split())
            matches: List[str] = kubectl_pattern.findall(nodes.split("\n", 1)[1])

            # after spawning 1 apatelet, there should be 2 nodes in total
            self.assertEqual(len(matches), 2)

            # assert two of these 1 node is an apatelet
            apatelet = [i for i in matches if i[0].startswith("apatelet-")]
            self.assertEqual(len(apatelet), 1)

            apatelet = apatelet[0]

            # wait for it to be running
            if apatelet[1] == "Ready":
                break
            else:
                print("waiting for Apatelet to be ready")
                # wait 5 seconds and try again
                time.sleep(5)
        else:
            self.fail("Node never became ready")


        # assert that the resources of this apatelet match the spec
        apateletname = apatelet[0]

        description = self.call_kubectl("describe", "nodes", apateletname)
        m = describe_pattern.search(description)
        self.assertIsNotNone(m)
        # assert cpu is 1000
        self.assertEqual(m.group("cpu").strip(), "1e3")
        # assert ephemeral storage is 120G
        self.assertEqual(m.group("eph").strip(), str(120 * 1024 * 1024 * 1024))
        # assert memory is 5G
        self.assertEqual(m.group("memory").strip(), str(5 * 1024 * 1024 * 1024))
        # assert storage is 5T
        self.assertEqual(m.group("storage").strip(), str(5 * 1024 * 1024 * 1024 * 1024))

        print("running scenario")

        # start emulating
        self.call_apate_cli("run", stdin="\n")

        # now our one node should fail soon
        for i in range(10):
            # get running nodes
            nodes = self.call_kubectl(*"get nodes".split())
            print(nodes)
            matches: List[str] = kubectl_pattern.findall(nodes.split("\n", 1)[1])

            # find the apatelet
            apatelet = [i for i in matches if i[0].startswith("apatelet-")]
            self.assertEqual(len(apatelet), 1)
            apatelet = apatelet[0]

            # what's it's status?
            if apatelet[1] == "Ready":
                # wait while it's ready
                time.sleep(10)
                print("waiting while Apatelet is ready for it to fail")
            elif apatelet[1] == "NotReady":
                break
            else:
                self.fail("status was Neither Ready or NotReady")
        else:
            self.fail("node did not become NotReady")

        # let tearDown stop the cluster