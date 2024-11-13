
    {
      title: $t('common.status'),
      field: 'status',
      slots: {
        default: (e) =>
          h(Switch, {
            checked: e.row.status === 1,
            onClick: () => {
              const newStatus = e.row.status === 1 ? 2 : 1;
              update{{.modelName}}({ id: e.row.id, status: newStatus }).then(() => {
                e.row.status = newStatus;
              });
            },
          }),
      },
    },