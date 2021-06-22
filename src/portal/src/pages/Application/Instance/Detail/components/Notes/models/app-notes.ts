/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 16:22:29
 * @Description: file content
 */
import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ApplicationsApiInstance } from '@/openApi';
import { IAppServiceState } from 'umi';

export const AppNotesAction = {
  UPDATE_STATUS: 'appNotes/updateState',
  GET_DATA: 'appNotes/getData',
};

export interface IAppNotesState {
  notes: string;
}

const defaultState: IAppNotesState = {
  notes: '',
};

const appNotes: Model = {
  namespace: 'appNotes',
  state: defaultState,
  effects: {
    *getData({ payload }, { call, put, select }) {
      let { serviceList }: IAppServiceState = yield select((state: any) => state.appService);
      const { data } = yield call(
        ApplicationsApiInstance.getApplicationNotes.bind(ApplicationsApiInstance, payload),
      );
      if (data && data.notes) {
        const hasNotes = data.notes.replaceAll('\n', '');
        let showNotes = '';
        if (hasNotes) {
          showNotes = data.notes;
          // key 为[[serviceName.podName]],value 为node_port
          serviceList.forEach((service) => {
            if (service.ports) {
              service.ports.forEach((pod) => {
                const key = `[[${service.name}.${pod.name}]]`;
                showNotes = showNotes.replaceAll(key, pod.node_port || '');
              });
            }
          });
        }
        yield put({
          type: 'updateState',
          payload: {
            notes: showNotes,
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

export default appNotes;
