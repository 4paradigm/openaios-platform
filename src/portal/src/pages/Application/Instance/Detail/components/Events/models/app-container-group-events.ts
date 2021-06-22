/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-03 11:34:44
 * @Description: file content
 */
import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ApplicationInstanceEvent } from '@/openApi/api';

export const AppContainerGroupEventsAction = {
  UPDATE_STATUS: 'appContainerGroupEvents/updateState',
};

export interface IAppContainerGroupEventsState {
  modalVisible: boolean;
  eventList: ApplicationInstanceEvent[];
}

const defaultState: IAppContainerGroupEventsState = {
  modalVisible: false,
  eventList: [],
};

const appContainerGroupEvents: Model = {
  namespace: 'appContainerGroupEvents',
  state: defaultState,
  reducers: {
    updateState(state, { payload }: IAction) {
      return {
        ...state,
        ...payload,
      };
    },
  },
};

export default appContainerGroupEvents;
