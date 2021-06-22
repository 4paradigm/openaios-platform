import { extend, RequestMethod, RequestOptionsInit } from 'umi-request';
import { notification } from 'antd';
import { ResponseBody } from '@/interfaces';
import keycloakClient from '@/keycloak';

//根请求路径
const ROOT_PATH = '/api';

// ERROR Code
const ERROR_CODE_MAP: { [key: string]: string } = {
  400: '发出的请求有错误，请检查请求参数',
  401: '用户没有权限',
  403: '用户得到授权，但是访问是被禁止的',
  404: '请求的地址不存在',
  406: '请求的格式不可得',
  410: '请求的资源被永久删除，且不会再得到的',
  422: '当创建一个对象时，发生一个验证错误',
  500: '服务器发生错误，请检查服务器',
  502: '网关错误',
  503: '服务不可用，服务器暂时过载或维护',
  504: '网关超时',
};

function handle401() {
  // 接口401需要记录下当前页面URL，以便登录后快速跳回原网页
  let refererURL = window.location.href;
  localStorage.setItem('REFERER_URL', refererURL);
}

const umiRequest = (isOpen: boolean) => {
  const _request = extend({
    //请求 url 前缀
    prefix: isOpen ? '' : ROOT_PATH,
    //请求 url 后缀
    suffix: '',
    //请求超时时间阈值
    timeout: 5000,
    //是否使用缓存 Get请求,开发时可以设置为false
    useCache: false,
    //缓存时长
    ttl: 6000,
    //最大缓存数
    maxCache: 0,
    //返回数据格式
    responseType: 'json',
    //是否获取源response
    getResponse: false, // 使用umi-request封装的response
    //默认请求是否带上cookie
    credentials: 'include',
  });
  //发起请求前拦截
  _request.interceptors.request.use((url, options) => {
    if (keycloakClient.isLoggedIn()) {
      const cb = () => {
        return Promise.resolve({
          url,
          options: {
            ...options,
            params: {
              ...options.params,
            },
            interceptors: true,
            headers: { Authorization: `Bearer ${keycloakClient.getToken()}` },
          },
        });
      };
      return keycloakClient.updateToken(cb) as any;
    } else {
      return {
        url: `${url}`,
        options: {
          ...options,
          params: {
            ...options.params,
          },
          interceptors: true,
        },
      };
    }
  });

  // 处理服务器异常、登录认证
  _request.interceptors.response.use(
    async (response, options): Promise<any> => {
      if (response.status === 200) {
        // 后端post成功之后，没有返回值，前端return一个成功字符串
        if (response.url.includes('storage/download')) {
          return response;
        }
        return response
          .clone()
          .json()
          .then((res) => {
            return res;
          })
          .catch((reason) => {
            return 'success';
          });
      } else {
        // 处理HTTP Code错误
        if (ERROR_CODE_MAP[response.status]) {
          if (response.status === 401) {
            handle401();
            return;
          } else if (response.status === 400) {
            const data: ResponseBody<any> = await response.clone().json();
            // 对没有token接口返回400特殊处理
            if (data.message === 'missing key in request header') {
              notification.error({
                message: ERROR_CODE_MAP['401'],
              });
              // TODO
              // history.push('/login');
              return;
            }
            notification.error({
              message: data.message || ERROR_CODE_MAP[response.status] || '请求出错，请重试',
              description: '',
            });
            return response.json().then(() => null) as any;
          } else {
            notification.error({
              message: ERROR_CODE_MAP[response.status],
            });
          }
          // 所有HTTP4X，5X状态的response，都重置res数据
          return response.json().then(() => null) as any;
        } else {
          // 其他错误或特定需求状态码
          console.log(`请求错误,错误码：${response.status};错误关键字：${response.statusText}`);
        }
      }

      // 保底返回数据
      return response;
    },
  );
  return _request;
};

class RequestWrapper {
  private request: any = null;
  constructor(umiRequest: RequestMethod) {
    this.request = umiRequest;
  }

  public get(url: string, params?: object | URLSearchParams) {
    return this.request.get(url, { params }).catch((e: any) => console.log('ServerError: ', e));
  }

  public post(url: string, data?: RequestOptionsInit['data']) {
    return this.request.post(url, { data }).catch((e: any) => console.log('ServerError: ', e));
  }

  public postFile(url: string, data?: RequestOptionsInit['data'], timeout?: number, signal?: any) {
    // signal参数用来取消请求
    let formdata = new FormData();
    formdata.append('file', data);
    return this.request
      .post(url, { data: formdata, timeout, signal })
      .catch((e: any) => console.log('ServerError', e));
  }

  public getFile(url: string, data?: RequestOptionsInit['data']) {
    return this.request
      .get(url, { responseType: 'blob' })
      .catch((e: any) => console.log('ServerError', e));
  }

  public put(url: string, data?: RequestOptionsInit['data']) {
    return this.request.put(url, { data }).catch((e: any) => console.log('ServerError: ', e));
  }

  public patch(url: string, data?: RequestOptionsInit['data']) {
    return this.request.patch(url, { data }).catch((e: any) => console.log('ServerError: ', e));
  }

  public delete(url: string, data?: RequestOptionsInit['data']) {
    return this.request.delete(url, { data }).catch((e: any) => console.log('ServerError: ', e));
  }

  public original() {
    return this.request;
  }
}

const request = new RequestWrapper(umiRequest(false));
const openRequest = new RequestWrapper(umiRequest(true));

export default request;
export { umiRequest, openRequest, ERROR_CODE_MAP };
