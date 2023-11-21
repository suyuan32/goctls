    if err := l.svcCtx.DB.Schema.Create(l.ctx, schema.WithForeignKeys(false)); err != nil {
        logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
        return nil, errorx.NewInternalError(err.Error())
    }