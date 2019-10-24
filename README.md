# Test Vectors Runner

This project is a reference implementation of a Test Vector runner which executes [Test Vectors](https://github.com/opennetworkinglab/testvectors) based tests for black-box testing of Stratum enabled switches.

## Run Test Vectors

testvectors-runner works with various switch targets including bmv2 switches and hardware switches. For running with bmv2 switch we provide a docker image  which deploys the bmv2 switch inside a docker container and another docker image for testvectors-runner binary. For running with hardware switches we also provide a script to deploy and run testvectors-runner binary on a server which has both gPRC and data plane connections to the hardware switch under test. In both cases you'll need to point testvectors-runner to the correct Test Vector files either downloaded from [Test Vectors repo](https://github.com/opennetworkinglab/testvectors) or created on your own.

### Run with bmv2 and Docker

Start `stratum-bmv2` switch with two dataplane ports for testing.
```bash
make bmv2
```
Start `tvrunner` container by mounting the test vectors directory.
```bash
make tv-runner TV_DIR=<PATH_TO_BMV2_TV>
```
> Note: replace `<PATH_TO_BMV2_TV>` with your bmv2 Test Vectors path.

Run tests by `make tests`. Or run each test category separately by `make pipeline` first and then `make p4runtime` or `make gnmi` or `make e2e`.

### Run with hardware switches

For now we only support deploying testvectors-runner on a server connected to the switch under test. The server should be able to talk to the switch via gRPC as well as have physical connections to the switch ports for data plane verification scenarios. We'll be supporting deployment directly on the switch under test in the future.

A `port-map.json` file is required to map switch port IDs to corresponding interface names on the server. Modify the `port-map.json` file under `tools/<PLATFORM>/` to match your environment setup.

Assuming testvectors-runner binary is already downloaded or built, go to `tools` directory and run `deploy.sh <USER@SERVER_IP>` to copy the binary and other files to the server.
> Note: modify the `PLATFORM`, `REMOTE_TV_RUNNER_DIR` and `TV_RUNNER_BIN` variables in `deploy.sh` as needed.

Once testvectors-runner binary is deployed, login to the server and use the Makefile located under `REMOTE_TV_RUNNER_DIR` to start the tests the same way as the docker environment.
> Note: make sure to download Test Vectors to the same server and point the Makefile to those files by modifying `TV_DIR` variable.

## Development Environment

### Docker development environment for testvectors-runner
Build and run bmv2 docker container
```bash
make bmv2
```
Build and run docker container by mounting this directory. The development container here runs on bmv2 container's network in order to access the data plane ports for testing.
```bash
make tv-runner-dev
```
Build go binary by running below command:
```bash
make build
```
Run specific tests by running below command:
```bash
build/_output/tv_runner -test.v -logLevel=info -tgFile=<TARGET_FILE_PATH> -tvFiles=/<TEST_VECTOR_FILE_PATH>
```
Run all tests in a specific directory by running below command:
```bash
build/_output/tv_runner -test.v -logLevel=info -tgFile=<TARGET_FILE_PATH> -tvDir=/<TEST_VECTOR_DIRECTORY_PATH>
```

### Linux development environment for testvectors-runner
[TODO] Prerequisites for buiding on local machine.

## Additional Documents
* [Test Vectors Runner Architecture](docs/architecture.md)
