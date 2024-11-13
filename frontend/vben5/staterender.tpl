
    {
      title: $t('common.status'),
      field: 'state',
      slots: {
        default: (e) =>
          h(Switch, {
            checked: e.row.state,
            onClick: () => {
              const newStatus = !e.row.state;
              update{{.modelName}}({ id: e.row.id, state: newStatus }).then(
                () => {
                  e.row.state = newStatus;
                },
              );
            },
          }),
      },
    },