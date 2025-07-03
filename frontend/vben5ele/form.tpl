<script lang="ts" setup>
    import type { {{.modelName}}Info } from '#/api/{{.folderName}}/model/{{.modelNameLowerCamel}}Model';

    import { ref } from 'vue';

    import { useVbenModal } from '@vben/common-ui';
    import { $t } from '@vben/locales';

    import { ElMessage } from 'element-plus';

    import { useVbenForm } from '#/adapter/form';
    import { create{{.modelName}}, update{{.modelName}} } from '#/api/{{.folderName}}/{{.modelNameLowerCamel}}';

    import { dataFormSchemas } from './schemas';

    defineOptions({
        name: '{{.modelName}}Form',
    });

    const record = ref();
    const isUpdate = ref(false);
    const gridApi = ref();

    async function onSubmit(values: Record<string, any>) {
        const result = isUpdate.value
            ? await update{{.modelName}}(values as {{.modelName}}Info)
            : await create{{.modelName}}(values as {{.modelName}}Info);
        if (result.code === 0) {
            ElMessage.success(result.msg);
            gridApi.value.reload();
        }
    }

    const [Form, formApi] = useVbenForm({
        handleSubmit: onSubmit,
        schema: [...(dataFormSchemas.schema as any)],
        showDefaultActions: false,
        layout: 'vertical',
    });

    const [Modal, modalApi] = useVbenModal({
        fullscreenButton: false,
        onCancel() {
            modalApi.close();
        },
        onConfirm: async () => {
            const validationResult = await formApi.validate();
            if (validationResult.valid) {
                await formApi.submitForm();
                modalApi.close();
            }
        },
        onOpenChange(isOpen: boolean) {
            isUpdate.value = modalApi.getData()?.isUpdate;
            record.value = isOpen ? modalApi.getData()?.record || {} : {};
            gridApi.value = isOpen ? modalApi.getData()?.gridApi : null;
            if (isOpen) {
                formApi.setValues(record.value);
            }
            modalApi.setState({
                title: isUpdate.value ? $t('{{.folderName}}.{{.modelNameLowerCamel}}.edit{{.modelName}}') : $t('{{.folderName}}.{{.modelNameLowerCamel}}.add{{.modelName}}'),
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
