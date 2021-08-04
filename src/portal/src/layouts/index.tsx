/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-07-07 15:12:25
 * @Description: file content
 */
import React, { useEffect, useState } from 'react';
import GlobalHeader from '@/components/GlobalHeader';
import './index.less';
import '../assets/styles/layout.less';
import { IAction } from '@/interfaces';
import zhCN from 'antd/lib/locale/zh_CN';
import { ConfigProvider, notification } from 'cess-ui';
import Loading from '@/components/Loading';
import keycloakClient from '@/keycloak';
import homeSvg from '@/assets/images/home.svg';
import { ICommonState, CommonActions } from 'umi';
import { useDispatch, useSelector } from 'react-redux';
import { Link, BrowserRouter as Router } from 'react-router-dom';

// 添加cnzz数据统计分析 --- start
if (!(window as any).addCnzzOnce) {
  // console.log('window.addCnzzOnce: ', (window as any).addCnzzOnce ); // 只执行一次
  var a = setTimeout(addCnzzfx, 500);
  function addCnzzfx() {
    (function () {
      var el = document.createElement('script');
      el.type = 'text/javascript';
      el.charset = 'utf-8';
      el.async = true;
      var ref: HTMLScriptElement = document.getElementsByTagName('script')[0];
      (ref.parentNode as Node & ParentNode).insertBefore(el, ref);
      el.src = 'https://w.cnzz.com/c.php?id=1280088195&async=1';
      const startTime = Date.now();
      el.onload = function () {
        const loadedTime = Date.now();
        const deltaTime = loadedTime - startTime;
        console.log('script loaded time:', deltaTime, 'ms');
        console.log('script dom:', this);
      };
    })();
  }

  (window as any).addCnzzOnce = true; //
}
// 添加cnzz数据统计分析 --- end

interface Props {
  dispatch: (a: IAction) => IAction;
}

const BasicLayout: React.FC<Props> = (props) => {
  const dispatch = useDispatch();
  const { loading, isMobile }: ICommonState = useSelector((state: any) => state.common);
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
      const scaleRate = window.innerWidth / (1200 + 36);
      document.body.style.transform = `scale(${scaleRate})`;
      (document.body.style as any)['transform-origin'] = 'left top';
      const rootDiv = document.getElementById('root');
      if (rootDiv) {
        rootDiv.style.height = `calc(100% * ${1 / scaleRate})`;
        rootDiv.style.width = `calc(100% * ${1 / scaleRate})`;
      }
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
      <div className={isMobile ? 'normal mobile' : 'normal'}>
        {loading ? (
          <Loading />
        ) : (
          <div className="layoutContainer">
            <div className="layoutSide">
              <Link to="/home">
                <img src={homeSvg} className="layoutHome" alt="home"></img>
              </Link>
            </div>
            <div className="layoutMain">
              <GlobalHeader />
              <div className="layoutNoHeader">{props.children}</div>
            </div>
          </div>
        )}
      </div>
    </ConfigProvider>
  );
};

export default BasicLayout;
