import { type BaseListResp } from '../../model/baseModel';

/**
 *  @description: {{.modelName}} info response
 */
export interface {{.modelName}}Info {
{{.infoData}}}

/**
 *  @description: {{.modelName}} list response
 */

export type {{.modelName}}ListResp = BaseListResp<{{.modelName}}Info>;
