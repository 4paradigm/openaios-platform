/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-22 14:47:07
 * @Description: file content
 */
import React, { useEffect, useState } from 'react';
import GlobalHeader from '@/components/GlobalHeader';
import styles from './index.less';
import '@/assets/styles/layout.less';
import { IAction } from '@/interfaces';
import zhCN from 'antd/lib/locale/zh_CN';
import { ConfigProvider, notification } from 'cess-ui';
import Loading from '@/components/Loading';
import keycloakClient from '@/keycloak';
import homeSvg from '@/assets/images/home.svg';
import { Link, useDispatch, useSelector, ICommonState, CommonActions } from 'umi';

interface Props {
  dispatch: (a: IAction) => IAction;
}

const BasicLayout: React.FC<Props> = (props) => {
  const dispatch = useDispatch();
  const { loading }: ICommonState = useSelector((state: any) => state.common);
  const [hasLogin, setHasLogin] = useState(false);

  useEffect(() => {
    let isMobile = false;
    // 判断移动端
    if (
      navigator.userAgent.match(
        /(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i,
      )
    ) {
      isMobile = true;
    }
    dispatch({
      type: CommonActions.UPDATE_STATE,
      payload: { isMobile },
    });
    if (isMobile) {
      document.body.style.zoom = window.innerWidth / (1200 + 36) + '';
    }
  }, []);
  useEffect(() => {
    if (!keycloakClient.isLoggedIn()) {
      const timer = setTimeout(() => {
        setHasLogin(true);
        notification.error({ message: '权限验证，请求超时！请稍后刷新重试' });
      }, 6 * 1000);
      keycloakClient.initKeycloak(() => {
        setHasLogin(true);
        clearTimeout(timer);
      });
    } else {
      setHasLogin(true);
    }
  }, []);

  useEffect(() => {
    // 初始化用户，每次进入之前先调用初始化用户，成功后再进行其他接口的调用
    if (hasLogin) {
      dispatch({ type: CommonActions.INIT_USER });
    }
  }, [hasLogin]);

  return (
    <ConfigProvider locale={zhCN}>
      <div className={styles.normal}>
        {loading ? (
          <Loading />
        ) : (
          <div className={styles.layoutContainer}>
            <div className={styles.layoutSide}>
              <Link to="/home">
                <img src={homeSvg} className={styles.layoutHome} alt="home"></img>
              </Link>
            </div>
            <div className={styles.layoutMain}>
              <GlobalHeader />
              <div>{props.children}</div>
            </div>
          </div>
        )}
      </div>
    </ConfigProvider>
  );
};

export default BasicLayout;
