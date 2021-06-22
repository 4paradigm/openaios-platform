/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-02 16:07:45
 * @Description: file content
 */
import React from 'react';
import { Link } from 'umi';
import { CessDropdown, Icon, Menu } from 'cess-ui';
import style from './index.less';
import LogoIcon from '@/assets/images/logo.png';
import keycloakClient from '@/keycloak';

function GlobalHeader() {
  const handleMenuClick = (e: any) => {
    if (e.key === 'logout') {
      if (keycloakClient) {
        keycloakClient.doLogout();
      }
    }
  };
  const menu = (
    <Menu onClick={handleMenuClick}>
      <Menu.Item key="logout">退出登录</Menu.Item>
    </Menu>
  );
  return (
    <div className={style.headerContanter}>
      <header className={style.header}>
        <div>
          <img src={LogoIcon} alt="logon" className={style.logo} />
          <span className={style.title}>
            AIOS 社区版
            {/* <Link to="/home">AIOS 社区版</Link> */}
          </span>
        </div>
        <div>
          <CessDropdown
            className={style.dropdown}
            title={
              <span>
                {keycloakClient.getUsername()}
                <Icon type="caret-bottom" />
              </span>
            }
            type="link"
            overlay={menu}
          />
        </div>
      </header>
    </div>
  );
}

export default GlobalHeader;
