/*
 * @Author: liyuying
 * @Date: 2021-05-21 16:04:54
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 14:36:18
 * @Description: 应用实例创建
 */
import React, { useEffect, useRef } from 'react';
import { history, useSelector, useDispatch } from 'umi';
import { Breadcrumb, Collapse, CessCreateBtnBar, Button } from 'cess-ui';
import { APPLICATION_CHART } from '@/router/url';
import {
  IApplicationInstanceCreateState,
  ApplicationInstanceCreateAction,
} from '../../models/application-instance-create';

import './index.less';
import AppInstanceCreateContent from './components/CreateContent';
import Loading from '@/components/Loading';

const { LeftBtnBar, RightBtnBar } = CessCreateBtnBar;
const breadcrumb = (
  <Breadcrumb>
    <Breadcrumb.Item href={APPLICATION_CHART}>应用市场</Breadcrumb.Item>
    <Breadcrumb.Item>创建应用实例</Breadcrumb.Item>
  </Breadcrumb>
);

const ApplicationInstanceCreate = ({ match, history }: any) => {
  const dispatch = useDispatch();
  const { loading, pageLoading }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  const formRef = useRef();

  const back = () => {
    history.push(APPLICATION_CHART);
  };
  /**
   * 创建实例
   */
  const handelCreateApplication = () => {
    (formRef.current as any).createApplication();
  };
  useEffect(() => {
    dispatch({
      type: ApplicationInstanceCreateAction.UPDATE_STATE,
      payload: {
        pageLoading: true,
      },
    });
    dispatch({
      type: ApplicationInstanceCreateAction.GET_TEMPLATE,
      payload: {
        name: match.params.name,
        category: history.location.query.category || '',
        version: history.location.query.version || '',
      },
    });
  }, []);
  return (
    <div className="application-create comm-create-page">
      {breadcrumb}
      {pageLoading ? (
        <Loading />
      ) : (
        <>
          <AppInstanceCreateContent ref={formRef}></AppInstanceCreateContent>
          <CessCreateBtnBar>
            <LeftBtnBar>
              <Button onClick={back}>退出</Button>
            </LeftBtnBar>
            <RightBtnBar>
              <Button type="primary" loading={loading} onClick={handelCreateApplication}>
                启动应用
              </Button>
            </RightBtnBar>
          </CessCreateBtnBar>
        </>
      )}
    </div>
  );
};
export default ApplicationInstanceCreate;
