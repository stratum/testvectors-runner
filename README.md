
# Test Vectors Runner

This project is a reference implementation of a Test Vector runner which executes [Test Vectors](https://github.com/opennetworkinglab/testvectors) based tests for black-box testing of Stratum enabled switches.

Build status (master): [![CircleCI](https://circleci.com/gh/opennetworkinglab/testvectors-runner/tree/master.svg?style=svg&circle-token=73bcc1fad5ddc6b34aede6a16f4b6bedc0630fc2)](https://circleci.com/gh/opennetworkinglab/testvectors-runner/tree/master)

`testvectors-runner` works with various switch targets that expose P4Runtime and gNMI, including [Stratum switches](https://github.com/stratum/stratum). To get started, you'll need Switch Under Test (SUT) and set of corresponding Test Vectors.```

## Start a Stratum Switch

Currently Stratum supports Barefoot Tofino and Broadcom Tomahawk devices, as well as the bmv2 software switch. Check [Stratum Project](https://github.com/stratum/stratum) for details of how to get Stratum running on supported devices.

We also provide a docker image which deploys the bmv2 software switch inside a docker container. To start `stratum-bmv2` switch with two dataplane ports for testing simply run:
```bash
make bmv2
```

> Note: `bmv2` container runs on `host` network and creates two veth pairs on host machine which are used for testing data plane scenarios. 

## Get Test Vectors

Download Test Vector files matching your SUT (tofino/bcm/bmv2) from [Test Vectors repo](https://github.com/opennetworkinglab/testvectors) or create your own Test Vectors.

In addition to Test Vector files, a `target.pb.txt` file and a `port-map.json` file are mandatory for starting testvectors-runner. `target.pb.txt` stores the IP and port that your SUT is using, and `port-map.json` stores a mapping between the switch port number used in Test Vectors and name of the interface on the node where testvectors-runner runs. Check [examples](https://github.com/stratum/testvectors/tree/master/tofino) in Test Vectors repo as well as the [readme](https://github.com/stratum/testvectors/blob/master/README.md) for more details.

## Testing with testvectors-runner

For running with hardware switches, testvectors-runner could be deployed on a server which has both gPRC and data plane connections to the SUT. We'll be supporting testvectors-runner deployment directly on the SUT soon. For running with bmv2 software switch, testvectors-runner needs to be deployed on the same network where the bmv2 container is deployed.

### Use existing testvectors-runner binary docker image
```bash
./tvrunner.sh --target <TARGET_FILE> --port-map <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR>
```
Above command uses [tvrunner](https://hub.docker.com/repository/docker/stratumproject/tvrunner/general) binary docker image, executes testvectors from `tv-dir` on switch running on `target`

### Build and use local testvectors-runner binary docker image
Build testvectors-runner binary image locally using below command:
```bash
docker build -t <IMAGE_NAME> -f build/test/Dockerfile .
```
Run tests with below command:
```bash
./tvrunner.sh --target <TARGET_FILE> --port-map <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR> --image-name <IMAGE_NAME>
```

In both cases, `tvrunner.sh` runs docker container in `host` network. To run docker container in another container's network, use below command:
```bash
./tvrunner.sh --target <TARGET_FILE> --port-map <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR> --network <NETWORK>
```

>Note: For more optional arguments, run *./tvrunner.sh -h*

### Use go run command to run tests
```bash
go run cmd/main/testvectors-runner.go --target <TARGET_FILE> --port-map <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR>
```

### Build go binary, run tests
Build testvectors-runner go binary using below command:
```bash
go build -o build/_output/tvrunner ./cmd/main
```
>Note: Alternatively, you can use *make build* to build the go binary

Use the executed binary to run tests
```bash
build/_output/tvrunner --target <TARGET_FILE> --port-map <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR>
```
>Note: For more optional arguments, run *go run cmd/main/testvectors-runner.go -h* or *build/_output/tvrunner -h*

## Additional Documents
* [Test Vectors Runner Architecture](docs/architecture.md)
