---
description: Run Quorum Key Manager from Docker image
sidebar_position: 1
---

# Run Quorum Key Manager from Docker image

Use the provided Docker image to run Quorum Key Manager (QKM) in a Docker container without installing QKM.

## Prerequisites

- [Docker](https://docs.docker.com/install/) and [Docker Compose](https://docs.docker.com/compose/install/)
- MacOS or Linux

  :::caution

  The Docker image does not run on Windows.

  :::

### Run Quorum Key Manager

Download the latest QKM [Docker compose file](https://github.com/ConsenSys/quorum-key-manager/blob/main/docker-compose.yml).

Specify a path to your [manifest file or folder](../HowTo/Use-Manifest-File/Overview.md) in an environment variable:

```bash
export HOST_MANIFEST_PATH=<PATH-TO-MANIFEST-FILE>
```

Start QKM using Docker Compose:

```bash
docker-compose -f docker-compose.yml up key-manager
```
