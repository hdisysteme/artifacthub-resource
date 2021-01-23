# Artifacthub Resource

[![Go Report Card](https://goreportcard.com/badge/github.com/PG2000/artifacthub-resource)](https://goreportcard.com/report/github.com/PG2000/artifacthub-resource)

Tracks and get new versions of Helm Charts which are registered 
at https://artifacthub.io/

## Source Configuration

| Parameter         | Required  | Example       | Description                           |
| ------------------|----------:|--------------:|--------------------------------------:|
| repository_name   | yes       | oteemo-charts | the repository name of the package    |
| package_name      | yes       | sonarqube     | the package name                      |
| api_key           | no        | <api-key>     | an api key                            |

Notes:

- if no api key is given it is possible that you will run into a request limit. 
You can obtain an api key from artifacthub.io by creating an account.
  

## Resource Actions

### check

Produces new versions for a helm chart ordered by the created_at date. 

A version is represented as follows:

- version: The Helm Chart Version
- created_at: Time of when the helm chart version was published

### in

Gets the requested version of the helm chart. 

The metadata information are available in your task destination.

- /app_version: The given app version of the helm chart
- /charts_url: The charts url of the helm chart 
- /chart_download_url: The url to download the helm chart
- /name: The name of the helm chart
- /organization_name: The organization name
- /repository_name: The repository name
- /version: The helm chart version

### out

no behavior implemented

## Example Pipeline

```yaml
---
resource_types:
- name: artifacthub
  type: docker-image
  source:
    repository: pg2000/artifacthub-resource
    tag: latest

resources:
- name: sonarqube
  type: artifacthub
  source:
    repository_name: oteemo-charts
    package_name: sonarqube

jobs:
  - name: job
    public: true
    plan:
      - get: sonarqube
        trigger: true
      - task: simple-task
        config:
          inputs:
            - name: sonarqube
          platform: linux
          image_resource:
            type: registry-image
            source: { repository: busybox }
          run:
            path: /bin/sh
            args:
            - -c
            - |
              cd sonarqube
              ls -lah
              cat app_version
              cat chart_download_url
              cat name
              cat version

```
