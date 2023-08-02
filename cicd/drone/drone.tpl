kind: pipeline
type: docker
name: {{.DroneName}}-{{.ServiceType}}
steps:
  - name: build-go
    image: golang:1.20.3
    depends_on: [clone]
    volumes:
      - name: go_cache
        path: /go/pkg/mod
    commands:
      - go env -w CGO_ENABLED=0
      - go env -w GOPROXY=https://goproxy.cn,direct
      - go env -w GOPRIVATE={{.GitGoPrivate}}
      - go mod tidy && go build -trimpath -ldflags "-s -w" -o app {{.ServiceName}}.go

  - name: build-{{.ServiceType}}
    image: plugins/docker:20
    environment:
      DRONE_REPO_BRANCH: {{.GitBranch}}
    depends_on: [build-go]
    settings:
      dockerfile: Dockerfile
      registry: {{.Registry}}
      repo: {{.Repo}}
      auto_tag: true
      insecure: true
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password
trigger:
  ref:
    - refs/tags/*
    - refs/heads/master

volumes:
  - name: go_cache
    host:
      path: /root/.go/cache
