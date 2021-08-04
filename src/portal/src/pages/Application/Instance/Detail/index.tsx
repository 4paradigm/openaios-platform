/*
 * @Author: liyuying
 * @Date: 2021-05-21 16:05:14
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-22 15:54:58
 * @Description: 查看应用实例
 */
import React, { useEffect } from 'react';
import {
  useDispatch,
  useSelector,
  ApplicationInstanceDetailAction,
  IApplicationInstanceDetailState,
} from 'umi';
import {
  Breadcrumb,
  Button,
  Divider,
  CessBaseInfoBar,
  CessCard,
  CessCreateTitleBar,
  Modal,
  Tabs,
} from 'cess-ui';
import { APPLICATION_INSTANCE } from '@/router/url';
import { APPLICATION_INSTANCE_DETAIL_TABS } from '@/constant/application';
import ApplicationInstanceDetailPods from './components/Pods';
import ApplicationInstanceDetailService from './components/Service';
import './index.less';
import moment from 'moment';
import ApplicationInstanceDetailNotes from './components/Notes';

const { CessTitle, CessRightBtn } = CessCreateTitleBar;
const breadcrumb = (
  <Breadcrumb>
    <Breadcrumb.Item href={APPLICATION_INSTANCE}>应用管理</Breadcrumb.Item>
    <Breadcrumb.Item>应用实例详情</Breadcrumb.Item>
  </Breadcrumb>
);
const ApplicationInstanceDetail = ({ match }: any) => {
  const dispatch = useDispatch();
  const { appInstance, loading }: IApplicationInstanceDetailState = useSelector(
    (state: any) => state.applicationInstanceDetail,
  );
  /**
   * 删除实例
   * @param record
   */
  const handelInstanceDelete = () => {
    Modal.confirm({
      title: '删除应用实例',
      content: `确定要删除【${appInstance.instance_name}】这个应用实例吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        dispatch({
          type: ApplicationInstanceDetailAction.DELETE,
          payload: appInstance.instance_name,
        });
      },
    });
  };
  useEffect(() => {
    dispatch({
      type: ApplicationInstanceDetailAction.GET_DATA,
      payload: {
        instance_name: match.params.name,
      },
    });
  }, []);
  return (
    <div className="application-detail comm-detail-page">
      {breadcrumb}
      <CessCard title="">
        <CessCreateTitleBar>
          <CessTitle>{appInstance.instance_name}</CessTitle>
          <CessBaseInfoBar>
            <span>实例状态：</span>
            {appInstance.status}
            <Divider type="vertical"></Divider>
            <span>应用名称：</span>
            {appInstance.chart_name}
            <Divider type="vertical"></Divider>
            <span>应用版本：</span>
            {appInstance.chart_version}
            <Divider type="vertical"></Divider>
            <span>创建时间：</span>
            {moment(new Date(appInstance.create_tm || '')).format('yyyy-MM-DD HH:mm:ss')}
          </CessBaseInfoBar>
          <CessRightBtn>
            <Button onClick={handelInstanceDelete} loading={loading}>
              删除实例
            </Button>
          </CessRightBtn>
        </CessCreateTitleBar>
        <Tabs defaultActiveKey={APPLICATION_INSTANCE_DETAIL_TABS.NOTES}>
          <Tabs.TabPane
            tab={APPLICATION_INSTANCE_DETAIL_TABS.NOTES}
            key={APPLICATION_INSTANCE_DETAIL_TABS.NOTES}
          >
            <ApplicationInstanceDetailNotes
              instance_name={match.params.name || ''}
            ></ApplicationInstanceDetailNotes>
          </Tabs.TabPane>
          <Tabs.TabPane
            tab={APPLICATION_INSTANCE_DETAIL_TABS.PODS}
            key={APPLICATION_INSTANCE_DETAIL_TABS.PODS}
          >
            <ApplicationInstanceDetailPods
              instance_name={match.params.name || ''}
            ></ApplicationInstanceDetailPods>
          </Tabs.TabPane>
          <Tabs.TabPane
            tab={APPLICATION_INSTANCE_DETAIL_TABS.SERVICE}
            key={APPLICATION_INSTANCE_DETAIL_TABS.SERVICE}
          >
            <ApplicationInstanceDetailService
              instance_name={match.params.name || ''}
            ></ApplicationInstanceDetailService>
          </Tabs.TabPane>
        </Tabs>
      </CessCard>
    </div>
  );
};
export default ApplicationInstanceDetail;
