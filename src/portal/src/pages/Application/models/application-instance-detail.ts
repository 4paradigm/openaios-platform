/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-15 15:31:02
 * @Description: file content
 */
import { history } from 'umi';
import { Model } from 'dva';
import { IAction, IApplicationInstance } from '@/interfaces';
import { ApplicationsApiInstance } from '@/openApi';
import { ApplicationInstanceMetadata } from '@/openApi/api';
import { APPLICATION_INSTANCE } from '@/router/url';

export const ApplicationInstanceDetailAction = {
  UPDATE_STATUS: 'applicationInstanceDetail/updateStatus',
  GET_DATA: 'applicationInstanceDetail/getData',
  DELETE: 'applicationInstanceDetail/deleteApplication',
};

export interface IApplicationInstanceDetailState {
  appInstance: ApplicationInstanceMetadata;
  loading: boolean;
}

const defaultState: IApplicationInstanceDetailState = {
  appInstance: {} as any,
  loading: false,
};

const applicationInstanceDetail: Model = {
  namespace: 'applicationInstanceDetail',
  state: defaultState,
  effects: {
    *getData({ payload }, { call, put }) {
      const { data } = yield call(
        ApplicationsApiInstance.getApplicationMetadata.bind(
          ApplicationsApiInstance,
          payload.instance_name,
        ),
      );
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            appInstance: data,
          },
        });
      }
    },
    *deleteApplication({ payload }, { call, put, select }) {
      yield put({
        type: 'updateState',
        payload: {
          loading: true,
        },
      });
      const { data } = yield call(
        ApplicationsApiInstance.deleteApplication.bind(ApplicationsApiInstance, payload),
      );
      if (data) {
        history.push(APPLICATION_INSTANCE);
      }
      yield put({
        type: 'updateState',
        payload: {
          loading: false,
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

export default applicationInstanceDetail;
