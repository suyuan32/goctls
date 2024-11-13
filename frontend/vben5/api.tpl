import { type BaseDataResp, type BaseListReq, type BaseResp, type Base{{if .useUUID}}UU{{end}}IDsReq, Base{{if .useUUID}}UU{{end}}IDReq } from '#/api/model/baseModel';
import { requestClient } from '#/api/request';
import { {{.modelName}}Info, {{.modelName}}ListResp } from './model/{{.modelNameLowerCamel}}Model';

enum Api {
  Create{{.modelName}} = '/{{.prefix}}/{{.modelNameSnake}}/create',
  Update{{.modelName}} = '/{{.prefix}}/{{.modelNameSnake}}/update',
  Get{{.modelName}}List = '/{{.prefix}}/{{.modelNameSnake}}/list',
  Delete{{.modelName}} = '/{{.prefix}}/{{.modelNameSnake}}/delete',
  Get{{.modelName}}ById = '/{{.prefix}}/{{.modelNameSnake}}',
}

/**
 * @description: Get {{.modelNameSpace}} list
 */

export const get{{.modelName}}List = (params: BaseListReq) => {
  return requestClient.post<BaseDataResp<{{.modelName}}ListResp>>(Api.Get{{.modelName}}List, params);
};

/**
 *  @description: Create a new {{.modelNameSpace}}
 */
export const create{{.modelName}} = (params: {{.modelName}}Info) => {
  return requestClient.post<BaseResp>(Api.Create{{.modelName}}, params);
};

/**
 *  @description: Update the {{.modelNameSpace}}
 */
export const update{{.modelName}} = (params: {{.modelName}}Info) => {
  return requestClient.post<BaseResp>(Api.Update{{.modelName}}, params);
};

/**
 *  @description: Delete {{.modelNameSpace}}s
 */
export const delete{{.modelName}} = (params: Base{{if .useUUID}}UU{{end}}IDsReq) => {
  return requestClient.post<BaseResp>(Api.Delete{{.modelName}}, params: params);
};

/**
 *  @description: Get {{.modelNameSpace}} By ID
 */
export const get{{.modelName}}ById = (params: Base{{if .useUUID}}UU{{end}}IDReq) => {
  return requestClient.post<BaseDataResp<{{.modelName}}Info>>(Api.Get{{.modelName}}ById, params);
};
