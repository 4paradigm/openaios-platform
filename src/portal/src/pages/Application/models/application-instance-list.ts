/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 15:30:16
 * @Description: file content
 */
import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ApplicationInstanceMetadata } from '@/openApi/api';
import { ApplicationsApiInstance } from '@/openApi';

export const ApplicationInstanceListAction = {
  UPDATE_STATUS: 'applicationInstanceList/updateStatus',
  GET_LIST: 'applicationInstanceList/getList',
  DELETE: 'applicationInstanceList/deleteApplication',
};

export interface IApplicationInstanceListState {
  dataSource: ApplicationInstanceMetadata[];
  total: number;
  currentPage: number;
}

const defaultState: IApplicationInstanceListState = {
  dataSource: [],
  total: 0,
  currentPage: 1,
};

const applicationInstanceList: Model = {
  namespace: 'applicationInstanceList',
  state: defaultState,
  effects: {
    *getList({ payload }, { call, put }) {
      const { data } = yield call(
        ApplicationsApiInstance.getApplicationList.bind(
          ApplicationsApiInstance,
          (payload - 1) * 10,
          10,
        ),
      );
      let currentPage = payload;
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            dataSource: data.item,
            total: data.total,
            currentPage: currentPage,
          },
        });
      }
    },
    *deleteApplication({ payload }, { call, put, select }) {
      const { data } = yield call(
        ApplicationsApiInstance.deleteApplication.bind(ApplicationsApiInstance, payload),
      );
      if (data) {
        let { currentPage, total } = yield select(
          (state: { applicationInstanceList: IApplicationInstanceListState }) =>
            state.applicationInstanceList,
        );
        if (total % 10 === 1) {
          // 如果删除的是当页的第一个也是唯一一个环境时，删除之后要把页码减1
          let nextPage = currentPage - 1;
          currentPage = nextPage ? nextPage : 1;
          yield put({
            type: 'updateState',
            payload: {
              currentPage: currentPage,
            },
          });
        }
        yield put({
          type: 'getList',
          payload: currentPage,
        });
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
  },
};

export default applicationInstanceList;
