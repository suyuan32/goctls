FROM {{.BaseImage}}

# Define the project name | 定义项目名称
ARG PROJECT={{.ServiceName}}
# Define the config file name | 定义配置文件名
ARG CONFIG_FILE={{.ServiceName}}.yaml
# Define the author | 定义作者
ARG AUTHOR="{{.Author}}"

LABEL org.opencontainers.image.authors=${AUTHOR}

WORKDIR /app
ENV PROJECT=${PROJECT}
ENV CONFIG_FILE=${CONFIG_FILE}
{{if .HasTimezone}}
ENV TZ={{.Timezone}}{{if .Chinese}}
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
{{end}}RUN apk update --no-cache && apk add --no-cache tzdata
{{end}}
COPY ./${PROJECT}_{{.ServiceType}} ./
COPY ./etc/${CONFIG_FILE} ./etc/
{{if .HasPort}}
EXPOSE {{.Port}}
{{end}}
ENTRYPOINT ./${PROJECT}_{{.ServiceType}} -f etc/${CONFIG_FILE}