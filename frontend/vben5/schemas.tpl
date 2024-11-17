import type { VxeGridProps } from '#/adapter/vxe-table';

import { h } from 'vue';

import { type VbenFormProps } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { z } from '#/adapter/form';

{{if .hasStatus}}
    import { message, Switch } from 'ant-design-vue';
    import { update{{.modelName}} } from '#/api/{{.folderName}}/{{.modelNameLowerCamel}}';{{end}}

export const tableColumns: VxeGridProps = {
  columns: [
    {
      type: 'checkbox',
      width: 60,
    },
{{.basicData}}
{{if .useBaseInfo}}
    {
      title: $t('common.createTime'),
      field: 'createdAt',
      formatter: 'formatDateTime',
    },
{{end}}
  ],
};

export const searchFormSchemas: VbenFormProps = {
  schema: [{{.searchFormData}}
  ],
};

export const dataFormSchemas: VbenFormProps = {
  schema: [
    {
      fieldName: 'id',
      label: 'ID',
      component: 'Input',
      dependencies: {
        show: false,
        triggerFields: ['id'],
      },
    },{{.formData}}
  ],
};
