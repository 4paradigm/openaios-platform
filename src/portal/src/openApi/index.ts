/*
 * @Author: liyuying
 * @Date: 2021-05-11 17:21:31
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-09 15:40:24
 * @Description: 导出api实例
 */
import axios from 'axios';
import {
  ApplicationsApi,
  AppstoreApi,
  CompetitionApi,
  ComputingResourceApi,
  EnvironmentsApi,
  FinishedApi,
  ImagesApi,
  LogsApi,
  ReleasesApi,
  StorageApi,
  UsersApi,
} from '@/openApi/api/';

import keycloakClient from '@/keycloak';
import { notification } from 'cess-ui';
import { ResponseBody } from '@/interfaces';
import { ERROR_CODE_MAP } from '@/services/request';

const _axios = axios.create();
//@ts-ignore
_axios.interceptors.request.use((config: any) => {
  if (keycloakClient.isLoggedIn()) {
    const cb = () => {
      config.headers.Authorization = `Bearer ${keycloakClient.getToken()}`;
      if (config.url.indexOf('/storage/download') > -1) {
        config.responseType = 'blob';
      }
      return Promise.resolve(config);
    };
    return keycloakClient.updateToken(cb);
  } else {
    return config;
  }
});
// 处理服务器异常、登录认证
_axios.interceptors.response.use(
  (response) => Promise.resolve(response),
  async (error) => {
    if (error && error.response && error.response.status) {
      const response = error.response;
      // 处理HTTP Code错误
      if (ERROR_CODE_MAP[response.status]) {
        if (response.status === 400) {
          const data: ResponseBody<any> = response.data;
          // 对没有token接口返回400特殊处理
          if (data.message === 'missing key in request header') {
            notification.error({
              message: ERROR_CODE_MAP['401'],
            });
            return;
          }
          notification.error({
            message: data.message || ERROR_CODE_MAP[response.status] || '请求出错，请重试',
            description: '',
          });
          return { data: null } as any;
        } else {
          notification.error({
            message: ERROR_CODE_MAP[response.status],
          });
        }
        // 所有HTTP4X，5X状态的response，都重置res数据
        return { data: null } as any;
      } else {
        // 其他错误或特定需求状态码
        console.log(`请求错误,错误码：${response.status};错误关键字：${response.statusText}`);
      }
    }
  },
);
const BASE_PATH = '/api'.replace(/\/+$/, '');

export const ApplicationsApiInstance = new ApplicationsApi(undefined, BASE_PATH, _axios);
export const AppstoreApiInstance = new AppstoreApi(undefined, BASE_PATH, _axios);
export const CompetitionApiApiInstance = new CompetitionApi(undefined, BASE_PATH, _axios);
export const ComputingResourceApiInstance = new ComputingResourceApi(undefined, BASE_PATH, _axios);
export const EnvironmentsApiInstance = new EnvironmentsApi(undefined, BASE_PATH, _axios);
export const FinishedApiInstance = new FinishedApi(undefined, BASE_PATH, _axios);
export const ImagesApiInstance = new ImagesApi(undefined, BASE_PATH, _axios);
export const LogsApiInstance = new LogsApi(undefined, BASE_PATH, _axios);
export const ReleasesApiInstance = new ReleasesApi(undefined, BASE_PATH, _axios);
export const StorageApiInstance = new StorageApi(undefined, BASE_PATH, _axios);
export const UserApiInstance = new UsersApi(undefined, BASE_PATH, _axios);
