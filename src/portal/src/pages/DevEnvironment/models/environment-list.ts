/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-07-05 20:47:21
 * @Description: file content
 */
import { Model } from 'dva';
import API from '@/services/api';
import { EnvironmentsApiInstance } from '@/openApi';
import { IAction } from '@/interfaces';
import { delay } from '@/utils';
import { message } from 'cess-ui';
import { ENVIRONMENT_STATUS } from '@/constant/environment';
import { EnvironmentRuntimeInfo, ApplicationInstanceEvent } from '@/openApi/api';

export const EnvironmentListAction = {
  UPDATE_STATUS: 'environmentList/updateState',
  GET_LIST: 'environmentList/getList',
  DELETE_ENVIRONMENT: 'environmentList/deleteEnvironment',
  OPEN_TERMINAL: 'environmentList/openWebTerminal',
};

export interface IEnvironmentListState {
  dataSource: EnvironmentRuntimeInfo[];
  total: number;
  currentPage: number;
  eventVisible: boolean;
  eventList: ApplicationInstanceEvent[];
}

const defaultState: IEnvironmentListState = {
  dataSource: [],
  total: 0,
  currentPage: 1,
  eventVisible: false,
  eventList: [],
};

const environmentList: Model = {
  namespace: 'environmentList',
  state: defaultState,
  effects: {
    *getList({ payload }, { call, put }) {
      const { data } = yield call(
        EnvironmentsApiInstance.getEnvironmentList.bind(
          EnvironmentsApiInstance,
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
        let hasIniting: boolean = false;
        data.item.forEach((element: any) => {
          if (
            element.state === ENVIRONMENT_STATUS.Pending ||
            element.state === ENVIRONMENT_STATUS.Unknown
          ) {
            hasIniting = true;
          }
        });
        // 有正在等待的需要刷新
        if (hasIniting) {
          yield call(delay, 3000);
          yield put({
            type: 'getList',
            payload: currentPage,
          });
        }
      }
    },
    *openWebTerminal({ payload }, { call, put }) {
      const data = yield call(API.openWebTerminal, payload);
      if (data && data.url) {
        window.open(data.url, '_blank');
      } else {
        message.error('打开Terminal失败！');
      }
    },
    *deleteEnvironment({ payload }, { call, put, select }) {
      const { data } = yield call(
        EnvironmentsApiInstance.deleteEnvironment.bind(EnvironmentsApiInstance, payload),
      );
      if (data) {
        let { currentPage, total } = yield select(
          (state: { environmentList: IEnvironmentListState }) => state.environmentList,
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

export default environmentList;
