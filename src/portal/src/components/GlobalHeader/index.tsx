/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-07-06 16:37:22
 * @Description: file content
 */
import React from 'react';
import { LogoutOutlined } from '@ant-design/icons';
import { Link, ICommonState } from 'umi';
import { useSelector } from 'react-redux';
import { CessDropdown, Icon, Menu, Button } from 'cess-ui';
import style from './index.less';
import LogoIcon from '@/assets/images/logo.png';
import keycloakClient from '@/keycloak';

function GlobalHeader() {
  const { isMobile }: ICommonState = useSelector((state: any) => state.common);
  const handleMenuClick = (e: any) => {
    if (e.key === 'logout') {
      if (keycloakClient) {
        keycloakClient.doLogout();
      }
    } else if (e.key === 'docs') {
      // 跳转至 查看文档页面
      const docUrl = window.location.origin + '/docs/';
      window.location.href = docUrl;
    }
  };
  const menu = (
    <Menu onClick={handleMenuClick}>
      <Menu.Item key="docs">使用文档</Menu.Item>
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
          {isMobile ? (
            <>
              {keycloakClient.getUsername()}
              <Button
                className={style.logoutBtn}
                onClick={() => {
                  if (keycloakClient) {
                    keycloakClient.doLogout();
                  }
                }}
                type="link"
              >
                <LogoutOutlined />
              </Button>
            </>
          ) : (
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
          )}
        </div>
      </header>
    </div>
  );
}

export default GlobalHeader;
