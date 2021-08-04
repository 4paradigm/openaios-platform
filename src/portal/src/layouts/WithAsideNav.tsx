/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-22 15:51:16
 * @Description: file content
 */
import React from 'react';
import AsideNav from '@/components/AsideNav';
import './index.less';
import { ICommonState } from 'umi';
import { useSelector } from 'react-redux';

const BasicLayout: React.FC = (props) => {
  const { isMobile }: ICommonState = useSelector((state: any) => state.common);
  return (
    <div className={isMobile ? `layoutMobile layoutWithNav` : 'layoutWithNav'}>
      <AsideNav />
      <div className="content">
        {props.children}
        {isMobile ? (
          ''
        ) : (
          <div className="footer">
            联系我们：opensource@4paradigm.com &nbsp;| &nbsp; 版权所有 © 2014—2021 第四范式
            保留所有权利 京 ICP 备 16062885 号-1 京公网安备11010802024074 号
          </div>
        )}
      </div>
    </div>
  );
};

export default BasicLayout;
