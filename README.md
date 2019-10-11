# Test Vectors Runner

This project is a reference implementation of a Test Vector runner which executes [Test Vectors](https://github.com/opennetworkinglab/testvectors) based tests for black-box testing of Stratum enabled switches.

## Run Test Vectors

testvectors-runner works with various switch targets including bmv2 switches and hardware switches. For running with bmv2 switches we provide a Dockerfile which deploys the runner binary and a bmv2 switch inside a docker container. For running with hardware switches we also provide a script to deploy and run testvectors-runner binary on a server which has both gPRC and data plane connections to the hardware switch under test. In both cases you'll need to point testvectors-runner to the correct Test Vector files either downloaded from [Test Vectors repo](https://github.com/opennetworkinglab/testvectors) or created on your own.

### Run with bmv2 and Docker

Build and run docker container by mounting the test vectors directory
```bash
docker build -t tv_runner --build_arg GIT_USER=<username> --build-arg GIT_PERSONAL_ACCESS_TOKEN=<personal_access_token> -f Dockerfile.test.bmv2 .

docker run -v <PATH_TO_BMV2_TV>:/root/tv/bmv2 --privileged --rm -it --name tv_runner tv_runner
```
> Note 1: replace `<username>` and `<personal_access_token>` with git username and personal access token. This is needed in order to access [Test Vectors repo](https://github.com/opennetworkinglab/testvectors)

> Note 2: replace `<PATH_TO_BMV2_TV>` with your bmv2 Test Vectors path.

Start bmv2 switch by `make switch` and run tests by `make tests`. Or run each test category separately by `make pipeline` first and then `make p4runtime` or `make gnmi` or `make e2e`.

> Note: login to container from another shell with the following command if needed:
> ```bash
> docker exec -it tv_runner /bin/bash
> ```

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
docker build -t stratum-bmv2 -f Dockerfile.bmv2 .

docker run --privileged --rm -it --name bmv2 stratum-bmv2
```
Start bmv2 switch by `make switch`

Build and run docker container by mounting this directory. The development container here runs on bmv2 container's network in order to access the data plane ports for testing.
```bash
docker build -t tvrunner-dev --build_arg GIT_USER=<username> --build-arg GIT_PERSONAL_ACCESS_TOKEN=<personal_access_token> -f build/dev/Dockerfile .

docker run --rm -it --name testdev --net container:bmv2 -v <THIS_DIR>:/root/testvectors-runner tvrunner-dev
```
> Note: replace `<username>` and `<personal_access_token>` with git username and personal access token. This is needed in order to access [Test Vectors repo](https://github.com/opennetworkinglab/testvectors)

Build go binary by running below command:
```bash
go build -o cmd/main/tv_runner cmd/main/testvectors-runner.go
```
Run specific tests by running below command:
```bash
cmd/main/tv_runner -test.v -logLevel=info -tgFile=tests/testdata/bmv2/target.pb.txt -tvFiles=tests/testdata/bmv2/PipelineConfig.pb.txt
```
Run all tests in a specific directory by running below command:
```bash
cmd/main/tv_runner -test.v -logLevel=info -tgFile=tests/testdata/bmv2/target.pb.txt -tvDir=tests/testdata/bmv2/gnmi/
```

### Linux development environment for testvectors-runner
Build and run bmv2 docker container on host network
```bash
docker build -t stratum-bmv2 -f build/bmv2/Dockerfile .

docker run --privileged --rm -it --name bmv2 --net=host stratum-bmv2
```
[TODO] Prerequisites for buiding on local machine.

## Additional Documents
* [Test Vectors Runner Architecture](docs/architecture.md)
