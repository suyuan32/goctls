    // {{.modelNameUpper}}

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("{{.routePrefix}}/{{.modelNameSnake}}/create").
        SetDescription("apiDesc.create{{.modelName}}").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("{{.routePrefix}}/{{.modelNameSnake}}/update").
        SetDescription("apiDesc.update{{.modelName}}").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("{{.routePrefix}}/{{.modelNameSnake}}/delete").
        SetDescription("apiDesc.delete{{.modelName}}").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("{{.routePrefix}}/{{.modelNameSnake}}/list").
        SetDescription("apiDesc.get{{.modelName}}List").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    apis = append(apis, l.svcCtx.DB.API.Create().
        SetServiceName("{{.serviceName}}").
        SetPath("{{.routePrefix}}/{{.modelNameSnake}}").
        SetDescription("apiDesc.get{{.modelName}}ById").
        SetAPIGroup("{{.modelNameSnake}}").
        SetMethod("POST"),
    )

    