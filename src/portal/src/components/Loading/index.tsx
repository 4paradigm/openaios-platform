import React from 'react';
import { Spin } from 'cess-ui';
import style from './index.less';

export default function () {
  return (
    <div className={style.container}>
      <Spin tip="正在加载中..." />
    </div>
  );
}
