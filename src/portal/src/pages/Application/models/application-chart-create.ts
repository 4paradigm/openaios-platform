/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:17:26
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-15 16:01:39
 * @Description: 创建我的应用
 */
import { Model } from 'dva';
import { history } from 'umi';
import { IAction } from '@/interfaces';
import { APPLICATION_CHART_CREATE_STEPS } from '@/constant/application';
import { Chart, ChartCategory } from '@/openApi/api';
import { AppstoreApiInstance } from '@/openApi';
import { message } from 'cess-ui';
import { APPLICATION_CHART } from '@/router/url';

export const ApplicationChartCreateAction = {
  UPDATE_STATE: 'applicationChartCreate/updateState',
  INIT_DATA: 'applicationChartCreate/initData',
  CREATE_APPLICATION: 'applicationChartCreate/createApplicationChart',
};

export interface IApplicationChartCreateState {
  loading: boolean;
  modalVisible: boolean;
  currentStep: number;
  appReadMeDoc: string;
  chartData: Chart;
  questionYaml: string;
}

const defaultState: IApplicationChartCreateState = {
  loading: false,
  modalVisible: false,
  currentStep: APPLICATION_CHART_CREATE_STEPS.UPLOAD.step,
  appReadMeDoc: '',
  questionYaml: '',
  chartData: {
    metadata: {
      name: '',
      description: '',
      version: '',
      url: '',
      icon_link: '',
    },
    files: {},
  },
};

const applicationChartCreate: Model = {
  namespace: 'applicationChartCreate',
  state: defaultState,
  effects: {
    *createApplicationChart({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          loading: true,
        },
      });
      const { data } = yield call(
        AppstoreApiInstance.uploadChart.bind(
          AppstoreApiInstance,
          ChartCategory.Private,
          payload.file,
          {
            timeout: 5000000,
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
        history.push(APPLICATION_CHART);
        message.success('我的应用创建成功');
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

export default applicationChartCreate;
