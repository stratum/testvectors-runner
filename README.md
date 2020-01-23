
# Test Vectors Runner

This project is a reference implementation of a Test Vector runner which executes [Test Vectors](https://github.com/opennetworkinglab/testvectors) based tests for black-box testing of Stratum enabled switches.

Build status (master): [![CircleCI](https://circleci.com/gh/stratum/testvectors-runner.svg?style=svg)](https://circleci.com/gh/stratum/testvectors-runner)

`testvectors-runner` works with various switch targets that expose P4Runtime and gNMI, including [Stratum switches](https://github.com/stratum/stratum). To get started, you'll need Switch Under Test (SUT) and set of corresponding Test Vectors.

## Start a Stratum Switch


To start Stratum's behavioral model software switch (`stratum-bmv2`) in a Docker container for testing, run:
```bash
make bmv2
```

> Note: The `bmv2` container runs on the `host` network and creates two veth pairs on host machine which are used for testing data plane scenarios. 

To start Stratum on hardware switches, including devices with Barefoot Tofino and Broadcom Tomahawk, visit the [Stratum Project repo](https://github.com/stratum/stratum) for details of how to get Stratum running on supported devices.

## Get Test Vectors

Download Test Vector files matching your SUT (tofino/bcm/bmv2) from [Test Vectors repo](https://github.com/opennetworkinglab/testvectors) or create your own Test Vectors.

In addition to Test Vector files, a `target.pb.txt` file and a `portmap.pb.txt` file are mandatory for starting testvectors-runner. `target.pb.txt` stores the IP and port that your SUT is using, and `portmap.pb.txt` stores information related to specific switch ports used in the Test Vectors. Check [examples](https://github.com/stratum/testvectors/tree/master/tofino) in Test Vectors repo as well as the [readme](https://github.com/stratum/testvectors/blob/master/README.md) for more details.

## Test with testvectors-runner

For running with hardware switches, testvectors-runner could be deployed on a server which has both gPRC and data plane connections to the SUT. For running with `stratum-bmv2`, testvectors-runner needs to be deployed on the same network where the bmv2 container is deployed.

When loopback mode is enabled on hardware switches, it's also supported to deploy testvectors-runner directly on the switch. See the loopback section below for more details.

### Use existing testvectors-runner binary docker image
```bash
./tvrunner.sh --target <TARGET_FILE> --portmap <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR>
```
Above command uses [tvrunner](https://hub.docker.com/repository/docker/stratumproject/tvrunner/general) binary docker image, executes testvectors from `tv-dir` on switch running on `target`. In addition to `--tv-dir` argument, you can also use `--tv-name <TEST_NAME_REGEX>` to run tests matching provided regular expression from `tv-dir`.

For example, assuming bmv2 container is deployed by `make bmv2` command and Test Vectors repo is downloaded to `~/testvectors`, first push a pipeline configuration to the bmv2 switch before running any tests:
```bash
./tvrunner.sh --target ~/testvectors/bmv2/target.pb.txt --portmap ~/testvectors/bmv2/portmap.pb.txt --tv-dir ~/testvectors/bmv2 --tv-name PipelineConfig
```

Above command finds and executes Test Vector with name `PipelineConfig.pb.txt` under `~/testvectors/bmv2`. Then run `p4runtime` test suite by:
```bash
./tvrunner.sh --target ~/testvectors/bmv2/target.pb.txt --portmap ~/testvectors/bmv2/portmap.pb.txt --tv-dir ~/testvectors/bmv2/p4runtime
```

### Build and use local testvectors-runner binary docker image
Build testvectors-runner binary image locally using below command:
```bash
docker build -t <IMAGE_NAME> -f build/test/Dockerfile .
```
Run tests with below command:
```bash
./tvrunner.sh --target <TARGET_FILE> --portmap <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR> --image-name <IMAGE_NAME>
```

In both cases, `tvrunner.sh` runs docker container in `host` network. To run docker container in another container's network, use below command:
```bash
./tvrunner.sh --target <TARGET_FILE> --portmap <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR> --network <NETWORK>
```

>Note: For more optional arguments, run *./tvrunner.sh -h*

### Use go run command to run tests
```bash
go run cmd/main/testvectors-runner.go --target <TARGET_FILE> --portmap <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR>
```

### Build go binary, run tests
Build testvectors-runner go binary using below command:
```bash
make build
```

Use the executed binary to run tests
```bash
./tvrunner --target <TARGET_FILE> --portmap <PORT_MAP_FILE> --tv-dir <TESTVECTORS_DIR>
```
>Note: For more optional arguments, run *go run cmd/main/testvectors-runner.go -h* or *./tvrunner -h*

### Loopback mode

To run tests in loopback mode just add `--dp-mode loopback` to the commands. It applies to all the options above. Take a Tofino switch as an example. First push a pipeline configuration by:
```bash
./tvrunner.sh --target ~/testvectors/tofino/target.pb.txt --portmap ~/testvectors/tofino/portmap.pb.txt --tv-dir ~/testvectors/tofino --tv-name PipelineConfig --dp-mode loopback
```

As part of loopback mode setup, extra `Insert*` Test Vectors need to be executed before running any tests (see more details [here](docs/loopback.md)).
```bash
./tvrunner.sh --target ~/testvectors/tofino/target.pb.txt --portmap ~/testvectors/tofino/portmap.pb.txt --tv-dir ~/testvectors/tofino --tv-name Insert.* --dp-mode loopback
```

Then run `p4runtime` test suite by:
```bash
./tvrunner.sh --target ~/testvectors/tofino/target.pb.txt --portmap ~/testvectors/tofino/portmap.pb.txt --tv-dir ~/testvectors/tofino/p4runtime --dp-mode loopback
```

After all tests are done, run the `Delete*` Test Vectors to clean up.
```bash
./tvrunner.sh --target ~/testvectors/tofino/target.pb.txt --portmap ~/testvectors/tofino/portmap.pb.txt --tv-dir ~/testvectors/tofino --tv-name Delete.* --dp-mode loopback
```

## Additional Documents
* [Test Vectors Runner Architecture](docs/architecture.md)
