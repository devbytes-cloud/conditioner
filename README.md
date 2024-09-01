# Conditioner

<img alt="conditioner.png" height="250" src="logo/conditioner.png" width="250"/>

This `kubectl` plugin allows you to add, update, or remove conditions on Kubernetes nodes. It's a handy tool for cluster administrators to manage node status conditions effectively.

## Features

- **Add a new condition** to a node with specific details.
- **Update an existing condition** on a node, including status, reason, and message.
- **Remove a condition** from a node.

## Prerequisites

- Kubernetes cluster
- `kubectl` installed and configured to communicate with your cluster

## Installation

### Manual 
To install the plugin, download the binary and make it executable:

```bash
curl -LO https://github.com/devbytse-cloud/conditioner/releases/download/{{ .Tag }}/{{ .ArtifactName }}
chmod +x ./kubectl-conditioner
mv ./kubectl-conditioner /usr/local/bin
```

### Homebrew

To install conditioner via homebrew you can use our tap.

```bash
brew tap devbytes-cloud/tap
brew install conditioner
```

or

```bash
brew install devbytes-cloud/tap/conditioner
```

### Docker

The docker image uses bitnami kubectl as its base. This means you have full access to kubectl and conditioner.
To yous the conditioner you just need to supply `conditioner` as your first input into the docker container

```bash
docker run devbytescloud/conditioner conditioner
```

The images can be found https://hub.docker.com/repository/docker/devbytescloud/conditioner

## Installation via Krew

Krew is a plugin manager for `kubectl`, which allows you to install `conditioner` easily. Follow these steps to install `conditioner` using Krew:

1. Ensure you have Krew installed. If not, follow the [Krew installation instructions](https://krew.sigs.k8s.io/docs/user-guide/setup/install/).

2. Run the following command to install `conditioner`:

```shell
kubectl krew install conditioner
```

## Usage

The general syntax for using the plugin is as follows:

```
kubectl conditioner [NODE_NAME] [FLAGS]
```

```shell
kubectl conditioner -h
The 'condition' command allows you to add, update, or remove status conditions on nodes.
You need to provide the node name as an argument and use flags to specify the details of the condition.
The '--type' flag is required and it specifies the type of condition you wish to interact with.
The '--status' flag sets the status for the specific status condition and it can be 'true', 'false', or left blank for 'unknown'.
The '--reason' flag sets the reason for the specific status condition.
The '--message' flag sets the message for the specific status condition.
If you wish to remove the condition from the node entirely, use the '--remove' flag.

Usage:
  conditioner [node name] [flags]

Examples:

# Add a new condition to a node
kubectl conditioner my-node --type Ready --status true --reason KubeletReady --message "kubelet is posting ready status"

# Update an existing condition on a node
kubectl conditioner my-node --type DiskPressure --status false --reason KubeletHasNoDiskPressure --message "kubelet has sufficient disk space available"

# Remove a condition from a node
kubectl conditioner my-node --type NetworkUnavailable --remove


Flags:
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/Users/ddymko/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -h, --help                           help for condition
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
      --message string                 Message for the specific status condition
  -n, --namespace string               If present, the namespace scope for this CLI request
  -r, --reason string                  Reason for the specific status condition
  -x, --remove                         If you wish to remove the condition from the node entirely
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --status string                  Status for the specific status condition [true, false]
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --type string                    (required): type of condition you wish to interact with
      --user string                    The name of the kubeconfig user to use
```

## Configuration

The application uses a configuration file named `.conditioner.json` located in the user's home directory. This file is automatically created if it does not exist.

The configuration file is in JSON format and contains the following fields:

- `prepend-whoami`: A boolean value that indicates whether to prepend the user's identity to the output. Default value is `false`.
- `allow-list`: An array of strings that represents a list of allowed conditions that can be used with conditioner. Default value is an empty array `[]`.

Here is an example of a configuration file:

```json
{
  "prepend-whoami": true,
  "allow-list": ["allowed-condition-1"]
}
```

This json configuration will only allowed `allowed-condition-1` to be used.
```sh
☁  ~  conditioner np-vm-02 --type random-condition --status false --reason conditionerExample --message "readme example"
Error: condition random-condition is not in allow-list [allowed-condition-1]

☁  ~  conditioner np-vm-02 --type allowed-condition-1 --status false --reason conditionerExample --message "readme example"
condition status allowed-condition-1 has been added on node np-vm-02
````

The `prepend-whoami` will append the current user to the `message`

```sh
allowed-condition-1   False   Sun, 01 Sep 2024 07:21:47 -0400   Sun, 01 Sep 2024 07:21:47 -0400   conditionerExample           ddymko: readme example
```

### Examples

- **Add a new condition** to a node:

  ```
  kubectl conditioner my-node --type Ready --status true --reason KubeletReady --message "kubelet is posting ready status"
  ```

- **Update an existing condition** on a node:

  ```
  kubectl conditioner my-node --type DiskPressure --status false --reason KubeletHasNoDiskPressure --message "kubelet has sufficient disk space available"
  ```

- **Remove a condition** from a node:

  ```
  kubectl conditioner my-node --type NetworkUnavailable --remove
  ```

### Flags

- `--type` (required): The type of condition (e.g., Ready, DiskPressure).
- `--status`: The status of the condition (`true`, `false`, or leave blank for `unknown`).
- `--reason`: A machine-readable, camel-case reason for the condition's last transition.
- `--message`: A human-readable message indicating details about the last transition.
- `--remove`: If set, the specified condition will be removed from the node.

## Building From Source

To build the plugin from source, you'll need Go installed. Clone the repository and run:

```bash
go build -o kubectl-conditioner ./cmd
```

This command will create a binary named `kubectl-conditioner` in your current directory.

## Contributing

Contributions are welcome! Please feel free to submit issues, pull requests, or suggest features.

## License

This project is open source and available under the [Apache License](LICENSE).
