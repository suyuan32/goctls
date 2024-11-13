<script lang="ts" setup>
  import type { {{.modelName}}Info } from '#/api/sys/model/{{.modelNameLowerCamel}}Model';

  import { ref } from 'vue';

  import { useVbenModal } from '@vben/common-ui';
  import { $t } from '@vben/locales';

  import { message } from 'ant-design-vue';

  import { useVbenForm } from '#/adapter/form';
  import { update{{.modelName}} } from '#/api/sys/{{.modelNameLowerCamel}}';

  import { dataFormSchemas } from './schema';

  defineOptions({
    name: '{{.modelName}}Form',
  });

  const record = ref();
  const isUpdate = ref(false);
  const gridApi = ref();

  async function onSubmit(values: Record<string, any>) {
    const result = await updateUser(values as {{.modelName}}Info);
    if (result.code === 0) {
      message.success(result.msg);
    }
  }

  const [Form, formApi] = useVbenForm({
    handleSubmit: onSubmit,
    schema: [...(dataFormSchemas.schema as any)],
  showDefaultActions: false,
  });

  const [Modal, modalApi] = useVbenModal({
    fullscreenButton: false,
    onCancel() {
      modalApi.close();
    },
    onConfirm: async () => {
      await formApi.submitForm();
      modalApi.close();
    },
    onOpenChange(isOpen: boolean) {
      isUpdate.value = modalApi.getData()?.isUpdate;
      record.value = isOpen ? modalApi.getData()?.record || {} : {};
      gridApi.value = isOpen ? modalApi.getData()?.gridApi : null;
      if (isOpen) {
        formApi.setValues(record.value);
      }
      modalApi.setState({
        title: isUpdate.value ? $t('sys.{{.modelNameLowerCamel}}.edit{{.modelName}}') : $t('sys.user.add{{.modelName}}'),
      });
    },
  });

  defineExpose(modalApi);
</script>
<template>
  <Modal>
    <Form />
  </Modal>
</template>
