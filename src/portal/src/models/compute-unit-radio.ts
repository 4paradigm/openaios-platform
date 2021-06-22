/*
 * @Author: liyuying
 * @Date: 2021-05-27 19:22:23
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 18:56:06
 * @Description: file content
 */
import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ComputingResourceApiInstance } from '@/openApi';
import { ComputeUnitSpec } from '@/openApi/api';

export const ComputeUnitRadioAction = {
  GET_DATA: 'computeUnitRadio/getComputeUnit',
};
export interface IComputeUnitRadioState {
  computeUnitList: ComputeUnitSpec[];
}

const defaultState: IComputeUnitRadioState = {
  computeUnitList: [],
};

const computeUnitRadio: Model = {
  namespace: 'computeUnitRadio',
  state: defaultState,
  effects: {
    *getComputeUnit({ payload }, { call, put }) {
      const { data } = yield call(
        ComputingResourceApiInstance.getComputingUnitSpecs.bind(ComputingResourceApiInstance),
      );
      if (data) {
        const list = (data || []).sort((item1: ComputeUnitSpec, item2: ComputeUnitSpec) => {
          return (item1.name || '').localeCompare(item2.name || '');
        });
        yield put({
          type: 'updateState',
          payload: {
            computeUnitList: list,
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

export default computeUnitRadio;
