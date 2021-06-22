/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 14:58:45
 * @Description: file content
 */
import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ApplicationsApiInstance } from '@/openApi';
import { ApplicationInstancePodList, ApplicationInstancePod } from '@/openApi/api';
import API from '@/services/api';
import { message } from 'cess-ui';
import { POD_STATE } from '@/constant/application';
import { delay } from '@/utils';

export const AppPodsAction = {
  UPDATE_STATUS: 'appPods/updateState',
  GET_DATA: 'appPods/getData',
  OPEN_WEB_TERMINAL: 'appPods/openWebTerminal',
};

export interface IAppPodsState {
  podsList: ApplicationInstancePod[];
}

const defaultState: IAppPodsState = {
  podsList: [],
};
// [
//     {
//       name: 'openpitrix-hyperpitrix-deployment-cc98b',
//       events: [
//         {
//           type: 'Warning',
//           reason: 'FailedMount',
//           age: '54m（x1492 over 9d）',
//           from: 'Kubelet',
//           message: 'Unable to attach',
//         },
//       ],
//       state: '正常',
//       containers: [
//         {
//           name: 'hyperpitrix1',
//           image: 'openpitrix/openpitrix:v0.5.0',
//           state: '运行中',
//           ports: [
//             { container_port: '9100', protocol: 'tcp' },
//             { container_port: '9200', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//           ],
//         },
//         {
//           name:
//             'hyperpitrix2hyperpitrix2hyperpitrix2hyperpitrix2hyperpitrix2hyperpitrix2hyperpitrix2hyperpitrix2',
//           image:
//             'openpitrix/openpitrix:v0.4.0openpitrix:v0.4.0openpitrix:v0.4.0openpitrix:v0.4.0openpitrix:v0.4.0openpitrix:v0.4.0openpitrix:v0.4.0openpitrix:v0.4.0',
//           state: '运行中',
//           ports: [
//             { container_port: '9100', protocol: 'tcp' },
//             { container_port: '9200', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//           ],
//         },
//       ],
//       create_tm: '2021-05-30',
//     },
//     {
//       name: 'openpitrix-hyperpitrix-deployment-cc98a',
//       events: [],
//       state: '正常',
//       containers: [
//         {
//           name: 'hyperpitrix3',
//           image: 'openpitrix/openpitrix:v0.5.0',
//           state: '运行中',
//           ports: [
//             { container_port: '9100', protocol: 'tcp' },
//             { container_port: '9200', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//           ],
//         },
//         {
//           name: 'test4',
//           image: 'openpitrix/openpitrix:v0.4.0',
//           state: '运行中',
//           ports: [
//             { container_port: '9100', protocol: 'tcp' },
//             { container_port: '9200', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//             { container_port: '9300', protocol: 'tcp' },
//           ],
//         },
//       ],
//       create_tm: '2021-05-31',
//     },
//   ],
const appPods: Model = {
  namespace: 'appPods',
  state: defaultState,
  effects: {
    *getData({ payload }, { call, put }) {
      const { data }: { data: ApplicationInstancePodList } = yield call(
        ApplicationsApiInstance.getApplicationPods.bind(ApplicationsApiInstance, payload),
      );
      const instance_name = payload;
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            podsList: data.item || [],
          },
        });
        let hasIniting: boolean = false;
        (data.item || []).forEach((element: any) => {
          if (element.state === POD_STATE.PENDING) {
            hasIniting = true;
          }
        });
        // 有正在等待的需要刷新
        if (hasIniting) {
          yield call(delay, 3000);
          yield put({
            type: 'getData',
            payload: instance_name,
          });
        }
      }
    },
    *openWebTerminal({ payload }, { call, put }) {
      const data = yield call(API.openWebTerminal, payload.podName, payload.containerName);
      if (data && data.url) {
        window.open(data.url, '_blank');
      } else {
        message.error('打开Terminal失败！');
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

export default appPods;
