/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 17:16:03
 * @Description: 打开webTerminal是外部的接口不走openAPI SDK
 */
import { openRequest } from './request';

const API = {
  // 打开webTerminal
  openWebTerminal(pod: string, container: string = 'environment') {
    return openRequest.get(`/web-terminal/terminal`, { pod, container }).then((res: any) => {
      return res;
    });
  },
};

export default API;
