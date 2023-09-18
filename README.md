# Kratos Agent

[docs](./)

## Overview

Kratos Agent is a service that runs on each node in the cluster and is responsible for managing the node's lifecycle. It is responsible for the following:

- Registering the node with the Kratos Control Plane
- Reporting the node's status to the Kratos Control Plane

### APIs

- [x] /agent/clusters (GET)
- [x] /agent/services (GET)
- [x] /agent/services/group (GET)

### Usage

```bash
$ docker pull omalloc/kratos-agent:latest
$ docker run --rm --name kratos-agent -v ./configs:/configs omalloc/kratos-agent:latest
```
