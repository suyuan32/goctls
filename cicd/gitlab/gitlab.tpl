variables:
  REPO: docker.io

stages:
  - info
  - build
  - publish

info-job:
  stage: info
  script:
    - echo "Initialize the environment ..."
    - export DOCKER_USERNAME=$DOCKER_USERNAME
    - export DOCKER_PASSWORD=$DOCKER_PASSWORD
    - export REPO=$REPO

build-job:
  stage: build-docker
  script:
    - echo "Compiling the code and build the docker image ..."
    - go mod tidy
    - make build-linux
    - make docker
    - echo "Compilation and build are done."

deploy-job:
  stage: publish
  environment: production
  script:
    - echo "Publish docker images ..."
    - make publish-docker
    - echo "Docker images successfully published."