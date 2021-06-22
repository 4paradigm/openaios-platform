/*
 * @Author: liyuying
 * @Date: 2021-05-25 15:15:58
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-15 15:31:41
 * @Description: 创建应用
 */
import { Model } from 'dva';
import { history } from 'umi';
import { IAction, IApplicationInstance } from '@/interfaces';
import JsYmal from 'js-yaml';
import { EDIT_APPLICATION_INSTANCE_ANSWER_MODE } from '@/constant/application';
import { Chart } from '@/openApi/api';
import { ApplicationsApiInstance, AppstoreApiInstance } from '@/openApi';
import { transferObjectFromFlatToLevles } from '@/utils';
import { questionYamlToQuestionConfig } from '../utils';
import { decode } from 'js-base64';
import { APPLICATION_INSTANCE } from '@/router/url';
import { message } from 'cess-ui';

export const ApplicationInstanceCreateAction = {
  UPDATE_STATE: 'applicationInstanceCreate/updateState',
  INIT_DATA: 'applicationInstanceCreate/initData',
  GET_TEMPLATE: 'applicationInstanceCreate/getApplicationTemplate',
  CREATE_INSTANCE: 'applicationInstanceCreate/createApplicationInstance',
};

export interface IApplicationInstanceCreateState {
  loading: boolean;
  editMode: EDIT_APPLICATION_INSTANCE_ANSWER_MODE;
  applicationInstance: IApplicationInstance;
  applicationTemplate: Chart;
  pageLoading: boolean;
}

const defaultState: IApplicationInstanceCreateState = {
  loading: false,
  pageLoading: true,
  editMode: EDIT_APPLICATION_INSTANCE_ANSWER_MODE.FORM,
  applicationTemplate: {
    metadata: {
      name: '',
      description: '',
      version: '',
      icon_link: '',
    },
    version_list: [],
    files: {},
  },
  applicationInstance: {
    chart_name: '',
    name: '',
    version: '',
    questionsGroupList: [],
    answers: {},
    answersYaml: '',
  },
};

const applicationInstanceCreate: Model = {
  namespace: 'applicationInstanceCreate',
  state: defaultState,
  effects: {
    *initData({ payload }, { call, put }) {
      /********************* 处理question.yaml *********************** */
      const appTemplate = payload;
      let questionYaml = '';
      if (!appTemplate.files) {
        appTemplate.files = {};
      }
      try {
        if (appTemplate.files) {
          for (const name in appTemplate.files) {
            if (name === 'questions.yaml' || name === 'questions.yml') {
              questionYaml = decode(appTemplate.files[name] || '');
            }
          }
        }
        const { questionsGroupList = [], answers = {} } = questionYamlToQuestionConfig(
          questionYaml,
        );
        // 添加answer.yaml
        appTemplate.files['answers.yaml'] = JsYmal.dump(transferObjectFromFlatToLevles(answers));
        // 判断默认编辑模式
        let editMode = '';
        if (questionsGroupList.length > 0) {
          editMode = EDIT_APPLICATION_INSTANCE_ANSWER_MODE.FORM;
        } else {
          editMode = EDIT_APPLICATION_INSTANCE_ANSWER_MODE.YAML;
        }

        yield put({
          type: 'updateState',
          payload: {
            editMode,
            pageLoading: false,
            applicationTemplate: appTemplate,
            applicationInstance: {
              name: '',
              chart_name: appTemplate.metadata?.name,
              version: appTemplate.metadata?.version,
              questionsGroupList: questionsGroupList,
              answers: {},
              answersYaml: '',
            },
          },
        });
      } catch (e) {
        console.log(e);
      }
    },
    *getApplicationTemplate({ payload }, { call, put }) {
      const { data } = yield call(
        AppstoreApiInstance.getAppstoreChart.bind(
          ApplicationsApiInstance,
          payload.category,
          payload.name,
          payload.version,
        ),
      );
      if (data) {
        yield put({
          type: 'initData',
          payload: data,
        });
      }
    },
    *createApplicationInstance({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          loading: true,
        },
      });
      const { data } = yield call(
        ApplicationsApiInstance.createApplication.bind(
          ApplicationsApiInstance,
          payload.appInstanceName,
          {
            url: payload.url,
            answers: payload.answers,
          },
        ),
      );
      yield put({
        type: 'updateState',
        payload: {
          loading: false,
        },
      });
      if (data) {
        history.push(APPLICATION_INSTANCE);
        message.success('应用实例创建成功');
      }
    },
  },
  reducers: {
    updateState(state, { payload }: IAction) {
      return {
        ...state,
        ...payload,
      };
    },
    initData(state, { payload }: IAction) {
      return {
        ...defaultState,
      };
    },
  },
};

export default applicationInstanceCreate;
