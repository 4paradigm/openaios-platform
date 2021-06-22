import { ChartCategory } from '@/openApi/api';
import { IAPPChartCategory } from '@/interfaces';

/*
 * @Author: liyuying
 * @Date: 2021-05-20 17:29:30
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-22 14:02:43
 * @Description: file content
 */
export enum APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE {
  STRING = 'string',
  MULTILINE = 'multiline',
  BOOLEAN = 'boolean',
  INT = 'int',
  ENUM = 'enum',
  COMPUTE_UNIT = 'ComputeUnit',
}

export const APPLICATION_TEMPLATE_QUESTION_ITEM_REQURE_MESSAGE = {
  [APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.STRING]: '请输入',
  [APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.MULTILINE]: '请输入',
  [APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.INT]: '请输入',
  [APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.BOOLEAN]: '请选择',
  [APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.ENUM]: '请选择',
  [APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.COMPUTE_UNIT]: '请选择',
};
export enum APPLICATION_INSTANCE_CREATE_COLLAPSE {
  DETAIL_DISCRIBTION = '详细介绍',
  CONFIOG_OPTIONS = '应用配置',
  PREVIEW = 'Charts预览',
}

export enum APPLICATION_BASIC_CONFIG_VAR {
  NAME = 'name',
}
/**
 * 不在answers中
 */
export const APPLICATION_BASIC_CONFIG_VAR_MAP: any = {
  [APPLICATION_BASIC_CONFIG_VAR.NAME]: true,
};
/**
 * 应用实例创建中编辑answer的模式
 */
export enum EDIT_APPLICATION_INSTANCE_ANSWER_MODE {
  FORM = 'FORM',
  YAML = 'YAML',
}

/**
 * 实例详情的tab
 */
export enum APPLICATION_INSTANCE_DETAIL_TABS {
  NOTES = 'Notes',
  PODS = 'Pods',
  SERVICE = 'Service',
}

/**
 * 创建我的应用的步骤
 */
export const APPLICATION_CHART_CREATE_STEPS: { [key: string]: { step: number; title: string } } = {
  UPLOAD: {
    step: 1,
    title: '应用上传',
  },
  APP_README: {
    step: 2,
    title: '简介编辑',
  },
  QUESTION_CONFIG: {
    step: 3,
    title: '表单配置',
  },
  PREVIEW: {
    step: 4,
    title: '预览与提交',
  },
};

/**
 * 创建应用按钮的标记
 */
export const CREATE_MY_APP_KEY = 'CREATE_MY_APP';
export const CHART_CATEGORY_PUBLIC = 'Public';

export const CHART_CATEGORY_LIST: IAPPChartCategory[] = [
  {
    category: CHART_CATEGORY_PUBLIC,
    categoryName: '内置应用',
    subCategorys: [
      {
        category: ChartCategory.PublicOfficial,
        categoryName: '官方应用',
      },
      {
        category: ChartCategory.PublicPractical,
        categoryName: '实战应用',
      },
      {
        category: ChartCategory.PublicCommunity,
        categoryName: '社区应用',
      },
    ],
  },
  // 隐藏自定义应用的功能
  // {
  //   category: ChartCategory.Private,
  //   categoryName: '我的应用',
  // },
];

export enum POD_STATE {
  PENDING = 'Pending',
  RUNNING = 'Running',
  SUCCEEDED = 'Succeeded',
  FAILED = 'Failed',
  UNKNOWN = 'Unknown',
}

export const POD_STATE_TAG: any = {
  [POD_STATE.PENDING]: 'default',
  [POD_STATE.RUNNING]: 'processing',
  [POD_STATE.SUCCEEDED]: 'success',
  [POD_STATE.FAILED]: 'error',
  [POD_STATE.UNKNOWN]: 'warning',
};

export const LOG_DEAFAULT_TAIL_LINE: number = 10000;
export const LOG_MIN_TAIL_LINE: number = 1;
export const LOG_MAX_TAIL_LINE: number = 100000;
