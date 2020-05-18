

# Diagram

![Apate's architecture](../documents/Architecture_flow.svg)

# Explanation

## Red (CLI workflow)

When using the CLI to start Apate, the following steps will be taken to start emulation:

1. First you can use the cli to start the Control plane.
2. Now the CLI will be waiting and polling the Control plane for status. Now You should use Kubectl step 1
3. The status will be queried from the Control plane store.
4. With all apatelets registered in the store, a heartbeat connection is set up to communicate their status.
5. When all Apatelets report their status to be healthy, the scenario can be started.

## Black (Kubectl workflow)

// TODO: Is this how you say this? --v
1. Using `Kubectl create` a crd can be spawned
2. Now kubernetes will call the node informer.
3. Now the informer will spawn the requested nodes. (TODO: is this right?)
4. The informer will also contact the Apatelets
5. This will generate a kubeconfig  
6. xxx
7. The cluster will report the spawned nodes to the store. This will be sent to the cli if a cli is polling for status.