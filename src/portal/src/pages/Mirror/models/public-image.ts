/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-12 11:44:20
 * @Description: file content
 */
import { Model } from 'dva';
import API from '@/services/api';
import { IAction } from '@/interfaces';
import { ImagesApiInstance } from '@/openApi';

export const PublicImageAction = {
  GET_LIST: 'publicImage/getList',
  GET_TOTAL: 'publicImage/getTotal',
};

export interface IMirrorPublicState {
  dataSource: any;
  total: number;
  currentPage: number;
}

const defaultState: IMirrorPublicState = {
  dataSource: [],
  total: 0,
  currentPage: 1,
};

const publicImage: Model = {
  namespace: 'publicImage',
  state: defaultState,
  effects: {
    *getList({ payload }, { call, put }) {
      const { data } = yield call(
        ImagesApiInstance.publicImagesGet.bind(ImagesApiInstance, '', payload, 10),
      );
      let currentPage = payload;
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            dataSource: data,
            currentPage: currentPage,
          },
        });
      }
    },
    *getTotal({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.publicImagesInfoGet.bind(ImagesApiInstance));
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            total: data.image_count,
          },
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

export default publicImage;
