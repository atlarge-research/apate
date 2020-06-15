# Apate Usage Guide


[[toc]]

## Prerequisites
The following prerequisites are required for running Apate:
* Docker
* Modern Linux Distro

::: details OS Support  
We have only tested Apate on Linux. However it is likely other operating systems will work as long as Docker is present.
:::

::: warning  
Docker is technically not required as you could run the binaries directly, not using the Docker containers is heavily discouraged and not supported.
:::

## Installation
There are two ways of installing Apate. The recommended way is just by downloading and installing the pre-built binary. The other way is to download the source code and compile the binary.

### Pre-built binary <Badge text="recommended"/>
To install the `apate-cli` simply download the binary and put it in a `bin` directory.

```sh
echo "80.60.83.220 TODO_URL" | sudo tee > /etc/hosts
```


To install it globally use this command
<!-- TODO: add url -->
```sh
curl -LO TODO_URL /usr/bin/apate-cli 
```

### From source
To install the `apate-cli` from source first clone the repo:

```bash
git clone https://github.com/atlarge-research/opendc-emulate-kubernetes
```

then move into the dir:
```bash
cd opendc-emulate-kubernetes
```

Finally you can install the `apate-cli` like this:
```bash
go install ./cmd/apate-cli/
```

Now `apate-cli` should be available under `$GOPATH/bin/apate-cli`.

::: tip  
It is also possible to build the controlplane and run it outside of a docker container. However that is out of scope for this guide.
:::

