# Test Vectors Runner Architecture
![architecture](images/architecture.png)
Test Vector Runner reads from one or multiple [Test Vector](https://github.com/stratum/testvectors) files and compiles them with an orchestrator by processing various types of Actions and Expectations. And based on the Action or Expectation type, the orchestrator calls corresponding framework modules to build and send/receive either gRPC messages or data plane packets. We also provide libraries which provide common functions for logging, test setup and teardown.
