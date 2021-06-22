/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 15:55:38
 * @Description: file content
 */
import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ApplicationsApiInstance } from '@/openApi';
import {
  ApplicationInstancePodList,
  ApplicationInstancePod,
  ApplicationInstanceServiceList,
  ApplicationInstanceService,
} from '@/openApi/api';

export const AppServiceAction = {
  UPDATE_STATUS: 'appService/updateState',
  GET_DATA: 'appService/getData',
};

export interface IAppServiceState {
  serviceList: ApplicationInstanceService[];
}

const defaultState: IAppServiceState = {
  serviceList: [],
};

const appService: Model = {
  namespace: 'appService',
  state: defaultState,
  effects: {
    *getData({ payload }, { call, put }) {
      const { data }: { data: ApplicationInstancePodList } = yield call(
        ApplicationsApiInstance.getApplicationServices.bind(ApplicationsApiInstance, payload),
      );
      yield put({
        type: 'updateState',
        payload: {
          serviceList: data.item || [],
        },
      });
    },
  },
  reducers: {
    updateState(state, { payload }: IAction) {
      return {
        ...state,
        ...payload,
      };
    },
  },
};

export default appService;
