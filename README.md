# OSPROXY

![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/sebastocorp/osproxy)
![GitHub](https://img.shields.io/github/license/sebastocorp/osproxy)

Object Storage Proxy is a little command to serve object storage objects from diferent buckets to your internal services.

## Motivation

TODO

## Flags

As every configuration parameter can be defined in the config file, there are only few flags that can be defined.
They are described in the following table:

| Name | Command | Default | Description |
|:---  |:---     |:---     |:---         |
| `--config`    | `serve` | `osproxy.yaml` | Path to the YAML config file |
| `--log-level` | `serve` |    `info`      | Verbosity level for logs |

> Output is thrown always in JSON as it is more suitable for automations

```console
osproxy run \
    --log-level=info
    --config="./hitman.yaml"
```

### Configuration

Current configuration version: `v1alpha1`

#### Transfer Service Parameters

Configuration to the transfer service to call in case of not found the object

| Name   | Default | Description |
|:---    |:---     |:---         |
| `transferService.host`     | `""` |  |
| `transferService.port`     | `""` |  |
| `transferService.endpoint` | `""` |  |

#### Object Storage Config Parameters

Configuration to the backend object storage service

| Name   | Default | Description |
|:---    |:---     |:---         |
| `osConfig.s3.endpoint`        | `""` |  |
| `osConfig.s3.accessKeyId`     | `""` |  |
| `osConfig.s3.secretAccessKey` | `""` |  |

#### Relation Parameters

Configuration to the backend object storage service

| Name   | Default | Description |
|:---    |:---     |:---         |
| `relation.type`                                      | `""` |  |
| `relation.buckets`                                   | `""` |  |
| `relation.buckets.<value>.frontend.bucketName`       | `""` |  |
| `relation.buckets.<value>.frontend.removePathPrefix` | `""` |  |
| `relation.buckets.<value>.frontend.addPathPrefix`    | `""` |  |
| `relation.buckets.<value>.backend.bucketName`        | `""` |  |
| `relation.buckets.<value>.backend.removePathPrefix`  | `""` |  |
| `relation.buckets.<value>.backend.addPathPrefix`     | `""` |  |

## How to deploy

This project is designed specially for Kubernetes, but also provides binary files and Docker images to make it easy to be deployed however wanted

### Binaries

Binary files for most popular platforms will be added to the [releases](https://github.com/sebastocorp/osproxy/releases)

### Kubernetes

You can deploy `osproxy` in Kubernetes using Helm as follows:

```console
helm repo add hitman https://sebastocorp.github.io/osproxy/

helm upgrade --install --wait osproxy \
  --namespace osproxy \
  --create-namespace sebastocorp/osproxy
```

> More information and Helm packages [here](https://sebastocorp.github.io/osproxy/)

### Docker

Docker images can be found in GitHub's [packages](https://github.com/sebastocorp/hitman/pkgs/container/osproxy) related to this repository

> Do you need it in a different container registry? I think this is not needed, but if I'm wrong, please, let's discuss
> it in the best place for that: an issue

## How to contribute

We are open to external collaborations for this project: improvements, bugfixes, whatever.

For doing it, open an issue to discuss the need of the changes, then:

- Fork the repository
- Make your changes to the code
- Open a PR and wait for review

The code will be reviewed and tested (always)

> We are developers and hate bad code. For that reason we ask you the highest quality
> on each line of code to improve this project on each iteration.

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
