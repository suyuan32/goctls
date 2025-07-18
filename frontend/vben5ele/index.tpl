<script lang="ts" setup>
  import type { VxeGridListeners, VxeGridProps } from '#/adapter/vxe-table';
  import type { {{.modelName}}Info } from '#/api/{{.folderName}}/model/{{.modelNameLowerCamel}}Model';

  import { h, ref } from 'vue';

  import { Page, useVbenModal, type VbenFormProps } from '@vben/common-ui';
  import { $t } from '@vben/locales';

  import { ElButton, ElMessageBox } from 'element-plus';
  import { isPlainObject } from 'remeda';

  import { useVbenVxeGrid } from '#/adapter/vxe-table';
  import { delete{{.modelName}}, get{{.modelName}}List } from '#/api/{{.folderName}}/{{.modelNameLowerCamel}}';
  import { type ActionItem, TableAction } from '#/components/table/table-action';

  import {{.modelName}}Form from './form.vue';
  import { searchFormSchemas, tableColumns } from './schema';

  // ---------------- form -----------------

  const [FormModal, formModalApi] = useVbenModal({
    connectedComponent: {{.modelName}}Form,
  });

  const showDeleteButton = ref<boolean>(false);

  const gridEvents: VxeGridListeners<any> = {
    checkboxChange(e) {
      showDeleteButton.value = e.$table.getCheckboxRecords().length > 0;
    },
    checkboxAll(e) {
      showDeleteButton.value = e.$table.getCheckboxRecords().length > 0;
    },
  };

  const formOptions: VbenFormProps = {
    // 默认展开
    collapsed: false,
    schema: [...(searchFormSchemas.schema as any)],
  // 控制表单是否显示折叠按钮
  showCollapseButton: true,
          // 按下回车时是否提交表单
          submitOnEnter: false,
  };

  // ------------- table --------------------

  const gridOptions: VxeGridProps<{{.modelName}}Info> = {
    checkboxConfig: {
      highlight: true,
    },
    toolbarConfig: {
      slots: {
        buttons: 'toolbar-buttons',
      },
    },
    columns: [
      ...(tableColumns.columns as any),
  {
    title: $t('common.action'),
            fixed: 'right',
          field: 'action',
          slots: {
  default: ({ row }) =>
            h(TableAction, {
              actions: [
                {
                  type: 'primary',
                  size: 'large',
                  link: true,
                  icon: 'clarity:note-edit-line',
                  onClick: openFormModal.bind(null, row),
                },
                {
                  icon: 'ant-design:delete-outlined',
                  type: 'danger',
                  size: 'large',
                  link: true,
                  popConfirm: {
                    title: $t('common.deleteConfirm'),
                    placement: 'left',
                    confirm: batchDelete.bind(null, [row]),
                  },
                },
              ] as ActionItem[],
            }),
  },
  },
  ],
  height: 'auto',
          keepSource: true,
          pagerConfig: {},
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await get{{.modelName}}List({
          page: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
        return res.data;
      },
    },
  },
  };

  const [Grid, gridApi] = useVbenVxeGrid({
    formOptions,
    gridOptions,
    gridEvents,
  });

  function openFormModal(record: any) {
    if (isPlainObject(record)) {
      formModalApi.setData({
        record,
        isUpdate: true,
        gridApi,
      });
    } else {
      formModalApi.setData({
        record: null,
        isUpdate: false,
        gridApi,
      });
    }
    formModalApi.open();
  }

  function handleBatchDelete() {
    ElMessageBox.confirm($t('common.deleteConfirm'), {
      type: 'warning',
    }).then(() => {
      const ids = gridApi.grid.getCheckboxRecords().map((item: any) => item.id);

      batchDelete(ids);
    });
  }

  async function batchDelete(ids: any[]) {
    const result = await delete{{.modelName}}({
      ids,
    });
    if (result.code === 0) {
      await gridApi.reload();
      showDeleteButton.value = false;
    }
  }
</script>

<template>
  <Page auto-content-height>
    <FormModal />
    <Grid>
      <template #toolbar-buttons>
        <ElButton
                v-show="showDeleteButton"
                type="danger"
                @click="handleBatchDelete"
        >
          {{`{{ $t('common.delete') }}`}}
        </ElButton>
      </template>

      <template #toolbar-tools>
        <ElButton type="primary" @click="openFormModal">
          {{`{{ $t('sys.{{.modelNameLowerCamel}}.add{{.modelName}}') }}`}}
        </ElButton>
      </template>
    </Grid>
  </Page>
</template>
