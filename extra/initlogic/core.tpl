    // {{.modelNameUpper}}

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("/{{.modelNameSnake}}/create").
        SetDescription("apiDesc.create{{.modelName}}").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("/{{.modelNameSnake}}/update").
        SetDescription("apiDesc.update{{.modelName}}").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("/{{.modelNameSnake}}/delete").
        SetDescription("apiDesc.delete{{.modelName}}").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("/{{.modelNameSnake}}/list").
        SetDescription("apiDesc.get{{.modelName}}List").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("/{{.modelNameSnake}}").
        SetDescription("apiDesc.get{{.modelName}}ById").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    