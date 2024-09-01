import "../base.api"

type (
    // The data of {{.modelEnglishName}} information | {{.modelChineseName}}信息
    {{.modelName}}Info {
        {{if .HasCreated}}Base{{if .useUUID}}UU{{end}}ID{{.IdType}}Info{{else}}Id        *{{if .useUUID}}string{{end}}{{.IdTypeLower}}    `json:"id,optional"`{{end}}{{.infoData}}
    }

    // The response data of {{.modelEnglishName}} list | {{.modelChineseName}}信息列表数据
    {{.modelName}}ListResp {
        BaseDataInfo

        // The {{.modelEnglishName}} list data | {{.modelChineseName}}信息列表数据
        Data {{.modelName}}ListInfo `json:"data"`
    }

    // The {{.modelEnglishName}} list data | {{.modelChineseName}}信息列表数据
    {{.modelName}}ListInfo {
        BaseListInfo

        // The {{.modelEnglishName}} list data | {{.modelChineseName}}信息列表数据
        Data  []{{.modelName}}Info  `json:"data"`
    }

    // Get {{.modelEnglishName}} list request params | {{.modelChineseName}}信息列表请求参数
    {{.modelName}}ListReq {
        PageInfo{{.listData}}
    }

    // The {{.modelEnglishName}} information response | {{.modelChineseName}}信息返回体
    {{.modelName}}InfoResp {
        BaseDataInfo

        // {{.modelEnglishName}} information | {{.modelChineseName}}信息数据
        Data {{.modelName}}Info `json:"data"`
    }
)

@server(
    jwt: Auth
    group: {{.groupName}}
    middleware: Authority{{if .hasRoutePrefix}}
    prefix: {{.routePrefix}}{{end}}
)

service {{.apiServiceName}} {
    // Create {{.modelEnglishName}} information | 创建{{.modelChineseName}}信息
    @handler create{{.modelName}}
    post /{{.modelNameSnake}}/create ({{.modelName}}Info) returns (BaseMsgResp)

    // Update {{.modelEnglishName}} information | 更新{{.modelChineseName}}信息
    @handler update{{.modelName}}
    post /{{.modelNameSnake}}/update ({{.modelName}}Info) returns (BaseMsgResp)

    // Delete {{.modelEnglishName}} information | 删除{{.modelChineseName}}信息
    @handler delete{{.modelName}}
    post /{{.modelNameSnake}}/delete ({{if .useUUID}}UU{{end}}IDs{{.IdType}}Req) returns (BaseMsgResp)

    // Get {{.modelEnglishName}} list | 获取{{.modelChineseName}}信息列表
    @handler get{{.modelName}}List
    post /{{.modelNameSnake}}/list ({{.modelName}}ListReq) returns ({{.modelName}}ListResp)

    // Get {{.modelEnglishName}} by ID | 通过ID获取{{.modelChineseName}}信息
    @handler get{{.modelName}}ById
    post /{{.modelNameSnake}} ({{if .useUUID}}UU{{end}}ID{{.IdType}}Req) returns ({{.modelName}}InfoResp)
}
