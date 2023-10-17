name: Build docker and publish
run-name: The pipeline for docker build
on: [release]

jobs:
  Build-Project:
    runs-on: ubuntu-latest
    steps:
      - run: git clone {{.repository}}{{if .china}}
      - run: go env -w GOPROXY=https://goproxy.cn,direct{{end}}
      - run: go mod tidy
        working-directory: {{.dir}}
      - run: make build-linux
        working-directory: {{.dir}}
      - run: make docker
        working-directory: {{.dir}}
      - run: make publish-docker
        working-directory: {{.dir}}