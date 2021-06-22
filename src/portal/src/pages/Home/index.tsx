/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-21 17:07:37
 * @Description: file content
 */
import React, { useMemo, useEffect, useRef } from 'react';
import { CessCard, Breadcrumb, Empty, Tooltip, Carousel, OverflowToolTip, Icon } from 'cess-ui';
import { Link, history } from 'umi';
import { useDispatch, useSelector } from 'react-redux';
import { HomeAction, IHomeState } from './models/home';
// import BalanceIcon from '@/assets/images/balance.svg';
// import CostIcon from '@/assets/images/cost.svg';
// import TimeIcon from '@/assets/images/time.svg';
import './index.less';
import OverFlowToolTip from 'cess-ui/lib/components/overflow-tool-tip';
import ResourceChart from './ResourceChart';
import { IMessage, IBanner } from '@/interfaces/bussiness';
import { BANNER_LSIT } from './data/banner/banner';
import { INDEX_MESSAGE, APPLICATION, DEV_ENVIRONMENT, APPLICATION_INSTANCE } from '@/router/url';

const Home = () => {
  const dispatch = useDispatch();
  const { taskInfo, msgInfo, userInfo }: IHomeState = useSelector((state: any) => state.home);
  const breadCrumb = useMemo(() => {
    return (
      <Breadcrumb>
        <Breadcrumb.Item>首页</Breadcrumb.Item>
      </Breadcrumb>
    );
  }, []);

  const handleViewMsgDetail = (msgInfo: IMessage) => {
    history.push(`${INDEX_MESSAGE}/${msgInfo.id}`);
  };
  /**
   * 点击banner
   * @param bannerInfo
   */
  const handleBannerClick = (bannerInfo: IBanner) => {
    console.log(bannerInfo);
  };
  const initHomeDat = async () => {
    dispatch({ type: HomeAction.GET_MESSAGE_INFO });
    await dispatch({ type: HomeAction.GET_INFO_TASK });
    await dispatch({ type: HomeAction.GET_USER_INFO });
  };
  useEffect(() => {
    initHomeDat();
  }, []);

  return (
    <div className="home">
      {breadCrumb}
      {/* <Carousel autoplay autoplaySpeed={3000} className="home-banner">
        {BANNER_LSIT.map((item) => {
          return (
            <img
              src={item.image}
              alt={item.title}
              key={item.id}
              onClick={() => {
                handleBannerClick(item);
              }}
            ></img>
          );
        })}
      </Carousel> */}
      <div className="home-container">
        <div className="home-left">
          <CessCard>
            <h3 className="card-title">活动与赛事</h3>
            <div className="home-msg-container">
              {msgInfo && msgInfo.length > 0
                ? msgInfo.map((msg) => {
                    return (
                      <p
                        className="home-msg-item"
                        key={msg.id}
                        onClick={() => {
                          handleViewMsgDetail(msg);
                        }}
                      >
                        <OverFlowToolTip title={msg.title} width={1000} line={1}></OverFlowToolTip>
                      </p>
                    );
                  })
                : '暂无公告'}
            </div>
          </CessCard>
          <CessCard>
            <h3 className="card-title">我的资源信息</h3>
            <div className="home-resource-bar">
              <div className="description">
                <p className="balance-item">
                  <span className="balance-item-label">余额：</span>
                  <Icon type="balance" className="balance-icon" />
                  <label className="balance">
                    <OverflowToolTip
                      title={(userInfo && userInfo.balance) || '0.00'}
                      width={220}
                      lineHeight={32}
                    ></OverflowToolTip>
                  </label>
                </p>
                <p className="balance-item">
                  <span className="balance-item-label">每分钟消耗：</span>
                  <Icon type="cost" className="balance-icon" />
                  <label className="balance">
                    <OverflowToolTip
                      title={(taskInfo && taskInfo.perCost) || '0.00'}
                      width={220}
                      lineHeight={32}
                    ></OverflowToolTip>
                  </label>
                </p>
                <p className="balance-item">
                  <span className="balance-item-label">预计可用：</span>
                  <Icon type="time" className="balance-icon" />
                  <label className="balance">
                    <OverflowToolTip
                      title={(userInfo && userInfo.costTime) || '0 min'}
                      width={220}
                      lineHeight={32}
                    ></OverflowToolTip>
                  </label>
                </p>
              </div>
              <div className="chart">
                {taskInfo && taskInfo.task_list && taskInfo.task_list.length > 0 ? (
                  <ResourceChart></ResourceChart>
                ) : (
                  <Empty msg="暂无任务"></Empty>
                )}
              </div>
            </div>
          </CessCard>
          <CessCard>
            <h3 className="card-title">我的实例信息</h3>
            <p>
              共 {(taskInfo && taskInfo.env_num) || 0} 个开发环境实例
              <Link to={DEV_ENVIRONMENT}> &gt;&gt;查看更多</Link>
            </p>
            <p>
              共 {(taskInfo && taskInfo.app_num) || 0} 个应用实例
              <Link to={APPLICATION_INSTANCE}> &gt;&gt;查看更多</Link>
            </p>
          </CessCard>
        </div>
      </div>
    </div>
  );
};

export default Home;
