
# Test Vectors Runner

This project is a reference implementation of a Test Vector runner which executes [Test Vectors](https://github.com/opennetworkinglab/testvectors) based tests for black-box testing of Stratum enabled switches.

Build status (master): [![CircleCI](https://circleci.com/gh/opennetworkinglab/testvectors-runner/tree/master.svg?style=svg&circle-token=73bcc1fad5ddc6b34aede6a16f4b6bedc0630fc2)](https://circleci.com/gh/opennetworkinglab/testvectors-runner/tree/master)

testvectors-runner works with various switch targets including hardware switches with Tofino/Tomahawk and bmv2 software switches. To get started, you'll first need to get a hardware or software switch running Stratum as Switch Under Test (SUT) and have corresponding Test Vectors downloaded or created on your own. Then you could either directly run a one-line command which downloads and runs a pre-built testvectors-runner docker image and executes specified Test Vectors against SUT, or make changes to the source code and build and run your own version of testvectors-runner. Check the following sections for detailed steps.

## Get a Stratum Enabled Switch

Currently Stratum supports Barefoot Tofino and Broadcom Tomahawk devices, as well as the bmv2 software switch. Check [Stratum Project](https://github.com/stratum/stratum) for details of how to get Stratum running on supported devices.

We also provide a docker image which deploys the bmv2 software switch inside a docker container. To start `stratum-bmv2` switch with two dataplane ports for testing simply run:
```bash
make bmv2
```

> Note: `bmv2` container runs on `host` network and creates two veth pairs on host machine which are used for testing data plane scenarios. 

## Get Test Vectors

Download Test Vector files matching your SUT (tofino/bcm/bmv2) from [Test Vectors repo](https://github.com/opennetworkinglab/testvectors) or create your own Test Vectors.

In addition to Test Vector files, a `target.pb.txt` file and a `port-map.json` file are mandatory for starting testvectors-runner. `target.pb.txt` stores the IP and port that your SUT is using, and `port-map.json` stores a mapping between the switch port number used in Test Vectors and name of the interface on the test node where testvectors-runner runs. Check [examples](https://github.com/stratum/testvectors/tree/master/tofino) in Test Vectors repo as well as the [readme](https://github.com/opennetworkinglab/testvectors) for more details.

## Testing with testvectors-runner

For running with bmv2 software switch, testvectors-runner needs to be deployed on the same node where the bmv2 container is deployed. For running with hardware switches, testvectors-runner could be deployed on a server which has both gPRC and data plane connections to the hardware SUT. We'll be supporting testvectors-runner deployment directly on the SUT soon.

### Use existing tvrunner binary docker image
