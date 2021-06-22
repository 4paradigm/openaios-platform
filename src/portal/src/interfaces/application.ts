import { APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE } from '@/constant/application';
import { Chart, ChartMetadata } from '@/openApi/api';
import { DataNode } from 'rc-tree/lib/interface';
/*
 * @Author: liyuying
 * @Date: 2021-05-20 17:30:19
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-09 15:36:43
 * @Description: file content
 */
export interface IQuestionItem {
  /* 字段key */
  variable: string;
  /* 字段展示名称 */
  label: string;
  /* 字段描述 */
  description?: string;
  /* 字段类型 */
  type: APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE;
  /* 字段是否必填 */
  required?: boolean;
  /* 字段的默认值 */
  default?: string;
  /* 配置分组，不同组展示divider */
  group?: string;
  /* 字符串类型的最大长度及最小长度 */
  min_length?: number;
  max_length?: number;
  /* 数字类型的取值范围 */
  min?: number;
  max?: number;
  /* 下拉选择的待选值 */
  options?: string[];
  /* 判读值合法的正则 */
  valid_chars?: string;
  /* 判读值不合法的正则 */
  invalid_chars?: string;
  /* 子问题组 */
  subquestions: IQuestionItem[];
  /* 判断是否需要展示当前的字段，如show_if: "serviceType=Nodeport" */
  show_if?: string;
  /* 判断是否需要子问题组，如show_subquestion_if: "true" */
  show_subquestion_if?: string;
  /* ------------------------前端展示的处理 start--------------- */
  /* 前端处理后增加的输入规则 */
  rules?: any[];
  /* 是否展示的条件 处理show_if得到的*/
  showConditions?: {
    key: string;
    value: any;
  }[];
  /* ------------------------前端展示的处理 end--------------- */
}
export interface IQuestionsGroup {
  groupName: string;
  questionList: IQuestionItem[];
}
/**
 * 应用实例
 */
export interface IApplicationInstance {
  /* 模版名称 */
  chart_name: string | number;
  /* 应用实例名称 */
  name: string;
  /* 应用版本 */
  version: string;
  /* 配置列表 */
  questionsGroupList?: IQuestionsGroup[];
  /* 配置的结果列表 */
  answers?: { [key: string]: any };
  /* 切换yaml配置时编辑的yaml */
  answersYaml?: string;
  /* 创建时间 */
  createTime?: string | number;
}

/**
 * 应用商店的对象
 */
export type IChartList = {
  [key: string]: ChartMetadata[];
};

export interface IAPPChartCategory {
  category: string;
  categoryName: string;
  subCategorys?: IAPPChartCategory[];
}
/**
 * 预览yaml的结构
 */
export interface IFileTree extends DataNode {
  parrentKey: string;
  label: string;
}
