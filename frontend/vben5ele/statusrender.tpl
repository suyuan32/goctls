
    {
      title: $t('common.status'),
      field: 'status',
      slots: {
        default: (e) =>
          h(ElSwitch, {
            activeValue: 1,
            inactiveValue: 2,
            modelValue: e.row.status,
            onChange: () => {
              const newStatus = e.row.status === 1 ? 2 : 1;
              update{{.modelName}}({ id: e.row.id, status: newStatus }).then(() => {
                e.row.state = newStatus;
                });
            },
          }),
      },
    },