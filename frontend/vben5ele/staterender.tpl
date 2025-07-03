
    {
      title: $t('common.status'),
      field: 'state',
      slots: {
        default: (e) =>
          h(ElSwitch, {
            activeValue: true,
            inactiveValue: false,
            modelValue: e.row.state,
            onChange: () => {
              const newStatus = !e.row.state;
              update{{.modelName}}({ id: e.row.id, state: newStatus }).then(() => {
                e.row.state = newStatus;
              });
            },
          }),
      },
    },