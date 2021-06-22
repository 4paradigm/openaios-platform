/*
 * @Author: liyuying
 * @Date: 2021-04-27 11:14:07
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-10 11:19:32
 * @Description: file content
 */

import Keycloak, { KeycloakInstance } from 'keycloak-js';
import { notification } from 'cess-ui';
const isProduction = process.env.NODE_ENV === 'production';
// @ts-ignore
const _kc: KeycloakInstance = new Keycloak(
  isProduction ? '/keycloak-config.json' : '/keycloak-config.json.example',
);

/**
 * Initializes Keycloak instance and calls the provided callback function if successfully authenticated.
 *
 * @param onAuthenticatedCallback
 */
const initKeycloak = (onAuthenticatedCallback: any) => {
  _kc
    .init({ onLoad: 'login-required' })
    .then((authenticated: any) => {
      if (authenticated) {
        onAuthenticatedCallback();
      } else {
        doLogin();
      }
    })
    .catch((e) => {
      if (e && e.error === 'access_denied') {
        notification.error({
          message: '您是否未接受服务条款，请重新登录并同意服务条款',
        });
      } else {
        notification.error({
          message: '权限失效',
        });
      }
      setTimeout(() => {
        doLogin();
      }, 2000);
    });
};

const doLogin = _kc.login;

const doLogout = _kc.logout;

const getToken = () => _kc.idToken;

const isLoggedIn = () => !!_kc.idToken;

const updateToken = (successCallback: any) =>
  _kc.updateToken(5).then(successCallback).catch(doLogin);

const getUsername = () => _kc.tokenParsed && (_kc.tokenParsed as any).preferred_username;

const getUserId = () => _kc.subject;

const hasRole = (roles: any) => roles.some((role: any) => _kc.hasRealmRole(role));

const keycloakClient = {
  initKeycloak,
  doLogin,
  doLogout,
  isLoggedIn,
  getToken,
  updateToken,
  getUsername,
  getUserId,
  hasRole,
};

export default keycloakClient;
