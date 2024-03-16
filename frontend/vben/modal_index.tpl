<template>
  <div>
    <BasicTable @register="registerTable">
      <template #tableTitle>
        <Button
          type="primary"
          danger
          preIcon="ant-design:delete-outlined"
          v-if="showDeleteButton"
          @click="handleBatchDelete"
        >
          {{.deleteButtonTitle}}
        </Button>
      </template>
      <template #toolbar>
        <a-button type="primary" @click="handleCreate">
          {{.addButtonTitle}}
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                popConfirm: {
                  title: t('common.deleteConfirm'),
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <BasicModal
      v-bind="$attrs"
      :title="getTitle"
      @register="register"
      :wrapperFooterOffset="50"
      :cancelText="t('common.closeText')"
      :showOkBtn="true"
      forceRender
      @ok="handleSubmit"
    >
      <BasicForm @register="registerForm" />
    </BasicModal>
  </div>
</template>
<script lang="ts">
  import { computed, createVNode, defineComponent, ref, unref } from 'vue';
  import { Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue/lib/icons';
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { Button } from '@/components/Button';
  import { useI18n } from 'vue-i18n';

  import { columns, searchFormSchema, formSchema } from './{{.modelNameLowerCamel}}.data';
  import { get{{.modelName}}List, delete{{.modelName}}, update{{.modelName}}, create{{.modelName}} } from '@/api/{{.folderName}}/{{.modelNameLowerCamel}}';
  import { BasicModal, useModal } from '/@/components/Modal';
  import { BasicForm, useForm } from '/@/components/Form';

  export default defineComponent({
    name: '{{.modelName}}Management',
    components: { BasicTable, BasicModal, TableAction, Button, BasicForm },
    setup() {
      const { t } = useI18n();
      const selectedIds = ref<number[] | string[]>();
      const showDeleteButton = ref<boolean>(false);
      const isUpdate = ref(false);

      const [registerTable, { reload }] = useTable({
        title: t('{{.folderName}}.{{.modelNameLowerCamel}}.{{.modelNameLowerCamel}}List'),
        api: get{{.modelName}}List,
        columns,
        formConfig: {
          labelWidth: 120,
          schemas: searchFormSchema,
        },
        useSearchForm: true,
        showTableSetting: true,
        bordered: true,
        showIndexColumn: false,
        clickToRowSelect: false,
        actionColumn: {
          width: 30,
          title: t('common.action'),
          dataIndex: 'action',
          fixed: undefined,
        },
        rowKey: 'id',
        rowSelection: {
          type: 'checkbox',
          onChange: (selectedRowKeys, _selectedRows) => {
            selectedIds.value = selectedRowKeys as {{if .useUUID}}string[]{{else}}number[]{{end}};
            showDeleteButton.value = selectedRowKeys.length > 0;
          },
        },
      });

      const [registerForm, { resetFields, setFieldsValue, validate }] = useForm({
        labelWidth: 160,
        baseColProps: { span: 24 },
        layout: 'vertical',
        schemas: formSchema,
        showActionButtonGroup: false,
      });

      const [register, { openModal }] = useModal();

      function handleCreate() {
        openModal(true);
        isUpdate.value = false;
        resetFields();
      }

      function handleEdit(record: Recordable) {
        openModal(true);
        isUpdate.value = true;
        setFieldsValue({
          ...record,
        });
      }

      const getTitle = computed(() =>
        !unref(isUpdate) ? t('sys.apis.add{{.modelName}}') : t('sys.apis.edit{{.modelName}}'),
      );

      async function handleSubmit() {
        const values = await validate();
        values['id'] = unref(isUpdate) ? Number(values['id']) : 0;
        let result = unref(isUpdate) ? await update{{.modelName}}(values) : await create{{.modelName}}(values);
        if (result.code === 0) {
          openModal(false);
          reload();
        }
      }

      async function handleDelete(record: Recordable) {
        const result = await delete{{.modelName}}({ ids: [record.id] });
        if (result.code === 0) {
          await reload();
        }
      }

      async function handleBatchDelete() {
        Modal.confirm({
          title: t('common.deleteConfirm'),
          icon: createVNode(ExclamationCircleOutlined),
          async onOk() {
            const result = await delete{{.modelName}}({ ids: selectedIds.value as {{if .useUUID}}string[]{{else}}number[]{{end}} });
            if (result.code === 0) {
              showDeleteButton.value = false;
              await reload();
            }
          },
          onCancel() {
            console.log('Cancel');
          },
        });
      }

      async function handleSuccess() {
        await reload();
      }

      return {
        t,
        registerTable,
        handleCreate,
        handleEdit,
        handleDelete,
        handleSuccess,
        handleBatchDelete,
        showDeleteButton,
        registerForm,
        register,
        handleSubmit,
        getTitle,
      };
    },
  });
</script>
