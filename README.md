# kubectl-credentials-broker

[![CI](https://github.com/takumakume/kubectl-credentials-broker/actions/workflows/ci.yml/badge.svg)](https://github.com/takumakume/kubectl-credentials-broker/actions/workflows/ci.yml)

**This software is:**

- Cli tool to work as kubectl plugin.
- Use [`kubernetes client-go credential plugin`](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#client-go-credential-plugins).
- It is possible to execute arbitrary commands before kubectl execution.Next, the specified client-certificate / key and token file is read and authentication is performed based on the specifications of `client-go credential plugin`.
- It can update the client-certificate / key and token by executing any command.

![image](docs/credentials-broker.jpeg)

# Install

## OSX

```sh
$ brew install takumakume/tap/kubectl-credentials-broker
```

# Usage

## `credentials-broker` command

This command is expected to be executed via client-go credentials pluigin.
There is nothing even to run alone.

```sh
$ kubectl credentials-broker --help
credentials-broker

Usage:
  credentials-broker [flags]
  credentials-broker [command]

Available Commands:
  help        Help about any command
  kubeconfig  kubeconfig

Flags:
      --before-exec-command string       A command line to run before responding to the credential plugin. For example, it can be used to update certificate and token files. (optional)
      --client-certificate-path string   PEM-encoded client certificate file path. Can contain CA certificate. If this flag is specified, --client-key-path is also required. (optional)
      --client-key-path string           PEM-encoded client key file path. (optional)
  -h, --help                             help for credentials-broker
      --token-path string                Token file path. (optional)
  -v, --version                          version for credentials-broker

Use "credentials-broker [command] --help" for more information about a command.
```

To use this, you need to set the following settings in your kubeconfig. and set current-context.

The `exec` section can be easily added with the `kubectl credentials-broker kubeconfig set` command. This will be explained later.

```yaml
users:
- name: user1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - credentials-broker
      - --before-exec-command
      - /path/to/update.sh
      - --client-certificate-path
      - /path/to/server.cert
      - --client-key-path
      - /path/to/server.key
      - --token-path
      - /path/to/token
      command: kubectl
      env:
      - name: FOO
        value: var
```

For example, make `/path/to/update.sh` as follows.

```sh
#!/bin/bash

openssl x509 -checkend 1 -noout -in /path/to/server.cert
if [ $? -ne 0 ]; then
    # Generate client certificate file if expired ...
    #  /path/to/server.cert (Can contain CA certificate)
    #  /path/to/server.key
fi

if [ ! -e "/path/to/token" ]; then
    # Generate token if not exists ...
    #  /path/to/token
fi
```

In this example, the certificate is reissued by some means, such as when the client certificate has expired. And, generate the token if the token file does not exist.

Make this script so that it can be executed every time kubectl is executed and before making a request to kube-apiserver.

Normally you don't run it directly, but if you do, you'll get the following results:

```sh
$ kubectl credentials-broker --before-exec-command /path/to/update.sh \
  --client-certificate-path /path/to/server.cert \
  --client-key-path /path/to/server.key \
  --token-path /path/to/token | jq .
```

```json
{
  "kind": "ExecCredential",
  "apiVersion": "client.authentication.k8s.io/v1beta1",
  "spec": {},
  "status": {
    "token": "XXXXXXXX",
    "clientCertificateData": "-----BEGIN CERTIFICATE-----\n ...",
    "clientKeyData": "-----BEGIN RSA PRIVATE KEY-----\n ..."
  }
}
```

You can use either `token` or `clientCertificateData` + `clientKeyData`.

Now supports the following APIs:

- `client.authentication.k8s.io/v1alpha1`
- `client.authentication.k8s.io/v1beta1`

**Tips**

If `token` is already defined as shown below, credentials plugin will not be kicked and must be removed.
`kubectl credentials-broker kubeconfig set` command removes `token`.

```yaml
users:
- name: user1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - credentials-broker # credentials plugin does not run
      - ...
      command: kubectl
    token: XXXXXXXX # kubectl using this token
```

## `credentials-broker kubeconfig set` command

このコマンドはcurrent-contextのuserにexec commandを追加します。

```
$ kubectl credentials-broker kubeconfig set
set

Usage:
  credentials-broker kubeconfig set [flags]

Flags:
      --before-exec-command string       A command line to run before responding to the credential plugin. For example, it can be used to update certificate and token files. (optional)
      --client-certificate-path string   PEM-encoded client certificate file path. Can contain CA certificate. If this flag is specified, --client-key-path is also required. (optional)
      --client-key-path string           PEM-encoded client key file path. (optional)
      --env stringToString               Environment variables to set when running the plugin. (optional) ex. 'HOGE=huga,FOO=bar' (default [])
      --exec-api-version string          API version to use when decoding the ExecCredentials resource (Default: client.authentication.k8s.io/v1beta1) (default "client.authentication.k8s.io/v1beta1")
  -f, --force                            Do not confirm overwriting of kubeconfig (Default: false)
  -h, --help                             help for set
      --token-path string                Token file path. (optional)
```

This command is used when you want to make the following kubeconfig.

Kubeconfig:
```yaml
users:
- name: user1 # current-context user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - credentials-broker
      - --before-exec-command
      - /path/to/update.sh
      - --token-path
      - /path/to/token
      command: kubectl
      env:
      - name: FOO
        value: var
      - name: HOGE
        value: huga
```

Command:
```sh
$ kubectl credentials-broker kubeconfig set \
  --before-exec-command /path/to/update.sh \
  --token-path=/tmp/token \
  --env=FOO=var,HOGE=huga
---

SHOW DIFF kubeconfig.old kubeconfig.new

---
continue? (y/N) (yes/no) [no]: yes
update successful
```

Do not confirm with `--force|-f` flag.
