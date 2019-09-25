# testvectors-runner

## Run Test Vectors

### Run with bmv2 and Docker

Build and run docker container by mounting the test vectors directory
```bash
docker build -t tv_runner -f Dockerfile.test.bmv2 .
docker run -v <PATH_TO_BMV2_TV>:/root/tv/bmv2 --privileged --rm -it --name tv_runner tv_runner
```

Start bmv2 switch by `make switch` and run tests by `make tests`. Or run each test category separately by `make pipeline` first and then `make p4runtime` or `make gnmi` or `make e2e`

Note: login to container from the 2nd shell with the following command if needed
```bash
docker exec -it tv_runner /bin/bash
```

### Run with hardware switches
[TODO]

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
docker build -t tvrunner-dev -f Dockerfile.dev .
docker run --rm -it --name testdev --net container:bmv2 -v <THIS_DIR>:/root/testvectors-runner tvrunner-dev
```
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
docker build -t stratum-bmv2 -f Dockerfile.bmv2 .
docker run --privileged --rm -it --name bmv2 --net=host stratum-bmv2
```
[TODO] Prerequisites for buiding on local machine.
