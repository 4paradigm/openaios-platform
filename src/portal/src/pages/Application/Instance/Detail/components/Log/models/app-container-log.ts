/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-18 14:31:23
 * @Description: file content
 */
import { Model } from 'dva';
import { notification } from 'cess-ui';
import { IAction } from '@/interfaces';
import keycloakClient from '@/keycloak';
import { LOG_DEAFAULT_TAIL_LINE } from '@/constant/application';

export const AppContainerLogAction = {
  UPDATE_STATUS: 'appContainerLog/updateState',
  GET_DATA: 'appContainerLog/getData',
};

export interface IAppContainerLogState {
  modalVisible: boolean;
  podName: string | number;
  containerName: string | number;
  logContent: string;
  xmlHttp: XMLHttpRequest;
  tailLines: number;
}

const defaultState: IAppContainerLogState = {
  modalVisible: false,
  containerName: '',
  podName: '',
  logContent: '',
  xmlHttp: null as any,
  tailLines: LOG_DEAFAULT_TAIL_LINE,
};

const appContainerLog: Model = {
  namespace: 'appContainerLog',
  state: defaultState,
  effects: {
    *getData({ payload }, { call, put, select }) {
      const { containerName, podName, tailLines, xmlHttp } = yield select(
        (state: { appContainerLog: IAppContainerLogState }) => state.appContainerLog,
      );
      if (xmlHttp) {
        xmlHttp.abort();
        yield put({
          type: 'updateState',
          payload: { xmlHttp: null },
        });
      }
      try {
        const currentXmlHttp = new XMLHttpRequest();
        currentXmlHttp.responseType = 'text';
        currentXmlHttp.onreadystatechange = () => {
          console.log('readyState:', currentXmlHttp.readyState);
        };
        currentXmlHttp.onprogress = (ev: any) => {
          if (payload.callBack) {
            payload.callBack(currentXmlHttp.responseText);
          }
        };
        const sendHttpRequest = () => {
          currentXmlHttp.open(
            'GET',
            `/api/log/pod/${podName}?container_name=${containerName}&tail_lines=${tailLines}`,
            true,
          );
          currentXmlHttp.setRequestHeader('Authorization', `Bearer ${keycloakClient.getToken()}`);
          currentXmlHttp.send(null);
        };
        keycloakClient.updateToken(sendHttpRequest);
        yield put({
          type: 'updateState',
          payload: {
            xmlHttp: currentXmlHttp,
          },
        });
      } catch (error) {
        notification.error({ message: '当前浏览器不支持流式日志' });
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

export default appContainerLog;
