/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-22 14:09:23
 * @Description: file conten
 */
import { IAction, IModel } from '@/interfaces';
import { UserApiInstance } from '@/openApi';
import { notification } from 'cess-ui';

export const CommonActions = {
  INIT_USER: 'common/initUser',
  UPDATE_STATE: 'common/updateState',
};

export interface ICommonState {
  loading: boolean;
  isMobile: boolean;
}

const defaultState: ICommonState = {
  loading: true,
  isMobile: false,
};

const common: IModel<ICommonState> = {
  namespace: 'common',
  state: defaultState,
  effects: {
    *initUser({ payload }, { call, put }) {
      const { data } = yield call(UserApiInstance.userInitPost.bind(UserApiInstance));
      if (data && data === 'success') {
        yield put({
          type: 'updateState',
          payload: {
            loading: false,
          },
        });
      } else {
        notification.error({ message: '当前账号初始化失败，请稍后重试' });
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

export default common;
