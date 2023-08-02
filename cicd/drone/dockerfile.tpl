FROM alpine

RUN apk update && apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY ./app /app/
COPY ./etc/{{.EtcYaml}}.yaml /app/etc/config.yaml
workdir /app/
ENTRYPOINT ["./app"]