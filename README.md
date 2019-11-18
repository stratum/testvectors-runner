
# Test Vectors Runner

This project is a reference implementation of a Test Vector runner which executes [Test Vectors](https://github.com/opennetworkinglab/testvectors) based tests for black-box testing of Stratum enabled switches.

Build status (master): [![CircleCI](https://circleci.com/gh/opennetworkinglab/testvectors-runner/tree/master.svg?style=svg&circle-token=73bcc1fad5ddc6b34aede6a16f4b6bedc0630fc2)](https://circleci.com/gh/opennetworkinglab/testvectors-runner/tree/master)

- [Testing Workflow](#testing-workflow)
  * [Testing with testvectors-runner Docker image](#testing-with-testvectors-runner-docker-image)
    + [Testing bmv2 switches](#testing-bmv2-switches)
    + [Testing hardware switches](#testing-hardware-switches)
  * [Testing with testvectors-runner binary](#testing-with-testvectors-runner-binary)
    + [Testing bmv2 switches with testvectors-runner binary](#testing-bmv2-switches-with-testvectors-runner-binary)
    + [Testing hardware switches with testvectors-runner binary](#testing-hardware-switches-with-testvectors-runner-binary)
- [Development Workflow](#development-workflow)
  * [Development in Docker environment](#development-in-docker-environment)
    + [Development with bmv2 switches](#development-with-bmv2-switches)
    + [Development with hardware switches](#development-with-hardware-switches)
  * [Development in Linux environment](#development-in-linux-environment)


## Testing Workflow

This section describes workflows for running testvectors-runner as a tester.

### Testing with testvectors-runner Docker image

testvectors-runner works with various switch targets including bmv2 switches and hardware switches. For running with bmv2 switch we provide a docker image which deploys the bmv2 switch inside a docker container and another docker image for testvectors-runner binary. For running with hardware switches the same testvectors-runner container could also be deployed on a server which has both gPRC and data plane connections to the hardware switch under test. In both cases you'll need to point testvectors-runner to the correct Test Vector files either downloaded from [Test Vectors repo](https://github.com/opennetworkinglab/testvectors) or created on your own.

#### Testing bmv2 switches

Start `stratum-bmv2` switch with two dataplane ports for testing by running:
```bash
make bmv2
```

Then start `tvrunner` container by mounting the test vectors directory:
```bash
make tvrunner-bmv2 TV_DIR=<PATH_TO_BMV2_TV>
```

> Note: replace `<PATH_TO_BMV2_TV>` with your bmv2 Test Vectors path.

> Note: the `tvrunner` container runs on `bmv2` container's network in order to access the data plane ports for testing.

Inside the `tvrunner` container, go to `tools` folder where the Makefile for running integration tests is located and run all test suites by
```bash
make tests
```

Or run each test category separately by `make pipeline` first and then `make p4runtime` or `make gnmi` or `make e2e`.

> Note: restarting `tvrunner` container is needed if `bmv2` container is restarted as both containers need to run in the same network.

#### Testing hardware switches

For now we only support deploying testvectors-runner on a server (hereinafter called the `test node`) connected to the switch under test. The test node should be able to talk to the switch via gRPC as well as have physical connections to the switch ports for data plane verification scenarios. We'll be supporting deployment directly on the switch under test soon.

Download this repo on the test node and start `tvrunner` container by mounting the test vectors directory:
```bash
make tvrunner-hw TV_DIR=<PATH_TO_TV>
```

> Note: replace `<PATH_TO_TV>` with your hardware Test Vectors path.

Then follow the same steps as described in [Testing bmv2 switches](#testing-bmv2-switches) section above to execute the tests against the switch under test.

### Testing with testvectors-runner binary

#### Testing bmv2 switches with testvectors-runner binary

[TODO] Prerequisites for testing on local machine.

#### Testing hardware switches with testvectors-runner binary

Assuming a test node is set up as described in [Testing hardware switches](#testing-hardware-switches) section above and testvectors-runner binary is already downloaded or built, go to `tools` directory and run `deploy.sh <USER@TEST_NODE_IP>` to copy the binary and Makefile to the test node.
> Note: modify the `REMOTE_TVRUNNER_DIR` and `TVRUNNER_BIN` variables in `deploy.sh` as needed.

Once testvectors-runner binary is deployed, login to the test node and use the Makefile located under `REMOTE_TVRUNNER_DIR/tools` to start the tests the same way as the docker environment.
> Note: make sure to download Test Vectors on the test node and point the Makefile to those files by modifying `TV_DIR` variable.

## Development Workflow

This section describes workflows for building and running testvectors-runner as a developer.

### Development in Docker environment

#### Development with bmv2 switches

Start `stratum-bmv2` switch as a container by running:
```bash
make bmv2
```

Then start a container for testvectors-runner development by mounting the test vectors directory:
```bash
make tvrunner-bmv2-dev TV_DIR=<PATH_TO_BMV2_TV>
```

> Note: replace `<PATH_TO_BMV2_TV>` with your bmv2 Test Vectors path.

Inside the `tvrunner` container, build `go` binary by running below command:
```bash
make build
```

Then follow the same steps as described in [Testing bmv2 switches](#testing-bmv2-switches) section above to execute the tests with the new testvectors-runner binary you just built.

> Note: restarting `tvrunner` container is needed if `bmv2` container is restarted as both containers need to run in the same network.

#### Development with hardware switches

Assuming a test node is set up as described in [Testing hardware switches](#testing-hardware-switches) section above.

Download this repo on the test node and start `tvrunner` container by mounting the test vectors directory:
```bash
make tvrunner-hw-dev TV_DIR=<PATH_TO_TV>
```

> Note: replace `<PATH_TO_TV>` with your hardware Test Vectors path.

Then follow the steps in [Development with bmv2 switches](#development-with-bmv2-switches) section to build `go` binary and run tests with the new testvectors-runner binary you just built against the switch under test.

### Development in Linux environment

[TODO] Prerequisites for buiding on local machine.

## Additional Documents
* [Test Vectors Runner Architecture](docs/architecture.md)
