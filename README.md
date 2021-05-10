# kubectl-credentials-broker

[![CI](https://github.com/takumakume/kubectl-credentials-broker/actions/workflows/ci.yml/badge.svg)](https://github.com/takumakume/kubectl-credentials-broker/actions/workflows/ci.yml)

**This software is:**

- Cli tool to work as kubectl plugin.
- To work as `client-go credential plugin` .
- It is possible to execute arbitrary commands before kubectl execution.Next, the specified client-certificate / key and token file is read and authentication is performed based on the specifications of `client-go credential plugin`.
- It can update the client-certificate / key and token by executing any command.

![image](docs/credentials-broker.jpeg)
