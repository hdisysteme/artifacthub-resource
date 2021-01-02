# Artifacthub Resource

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

### in

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
            - name: nexus
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