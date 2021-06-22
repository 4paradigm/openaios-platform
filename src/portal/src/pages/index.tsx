import React from 'react';
import { useHistory } from 'react-router-dom';

// 路由跳转决策组件
function Redirect() {
  const history = useHistory();

  const refererURL = localStorage.getItem('REFERER_URL');

  if (refererURL && refererURL.startsWith(window.location.origin)) {
    localStorage.removeItem('REFERER_URL');
    history.push(refererURL.replace(window.location.origin, ''));
  } else {
    history.push('/home');
  }

  return <div></div>;
}

export default Redirect;
