
  {
    title: t('common.status'),
    dataIndex: 'state',
    width: 50,
    customRender: ({ record }) => {
      if (!Reflect.has(record, 'pendingStatus')) {
        record.pendingStatus = false;
      }
      return h(Switch, {
        checked: record.state === true,
        checkedChildren: t('common.on'),
        unCheckedChildren: t('common.off'),
        loading: record.pendingStatus,
        onChange(checked, _) {
          record.pendingStatus = true;
          const newState = checked ? true : false;
          update{{.modelName}}({ id: record.id, state: newState })
            .then(() => {
              record.state = newState;
            })
            .finally(() => {
              record.pendingStatus = false;
            });
        },
      });
    },
  },