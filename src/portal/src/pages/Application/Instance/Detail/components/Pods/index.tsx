import React, { useEffect } from 'react';
import { useDispatch, useSelector, EnvironmentListAction } from 'umi';
import AppContainerLog from '../Log';
import { Button, Collapse, OverflowToolTip, Empty, Tag } from 'cess-ui';
import { AppContainerLogAction } from '../Log/models/app-container-log';
import AppContainerGroupEvents from '../Events';
import { AppContainerGroupEventsAction } from '../Events/models/app-container-group-events';
import { AppPodsAction, IAppPodsState } from './models/app-pods';
import {
  ApplicationInstanceContainerPorts,
  ApplicationInstanceContainer,
  ApplicationInstanceEvent,
  ApplicationInstancePod,
} from '@/openApi/api';
import DockerIcon from '@/assets/images/application/docker.svg';
import './index.less';
import moment from 'moment';
import { POD_STATE_TAG } from '@/constant/application';
interface Iprop {
  instance_name: string;
}
/*
 * @Author: liyuying
 * @Date: 2021-06-01 17:46:21
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-02 14:16:52
 * @Description:实例详情--Pods
 */
const ApplicationInstanceDetailPods = ({ instance_name }: Iprop) => {
  const dispatch = useDispatch();
  const { podsList }: IAppPodsState = useSelector((state: any) => state.appPods);
  /**
   * 展示日志弹窗
   */
  const handleShowLog = (pod: ApplicationInstancePod, container: ApplicationInstanceContainer) => {
    dispatch({
      type: AppContainerLogAction.UPDATE_STATUS,
      payload: {
        modalVisible: true,
        podName: pod.name,
        containerName: container.name,
        logContent: ``,
      },
    });
  };
  /**
   * 展示Event
   */
  const handleShowEvents = (eventList: ApplicationInstanceEvent[]) => {
    dispatch({
      type: AppContainerGroupEventsAction.UPDATE_STATUS,
      payload: {
        modalVisible: true,
        eventList,
      },
    });
  };
  /**
   * 打开webterminal
   */
  const openWebTerminal = (
    pod: ApplicationInstancePod,
    container: ApplicationInstanceContainer,
  ) => {
    dispatch({
      type: AppPodsAction.OPEN_WEB_TERMINAL,
      payload: {
        podName: pod.name,
        containerName: container.name,
      },
    });
  };
  const getShowPords = (pordList: ApplicationInstanceContainerPorts[]) => {
    const showPords: string[] = [];
    pordList.forEach((item) => {
      showPords.push(`${item.container_port}/${item.protocol}`);
    });
    return showPords.join('，');
  };
  /**
   * 获取pods信息
   */
  useEffect(() => {
    if (instance_name) {
      dispatch({
        type: AppPodsAction.GET_DATA,
        payload: instance_name,
      });
    }
  }, [instance_name, dispatch]);
  return (
    <div className="application-instance-detail-pods">
      {podsList && podsList.length > 0 ? (
        <>
          <Collapse defaultActiveKey={[]}>
            {podsList.map((item, index) => {
              return (
                <Collapse.Panel
                  header={
                    <div className="pods-group-header">
                      <div className="pods-group-title">
                        {item.name}{' '}
                        {item.state ? (
                          <Tag color={POD_STATE_TAG[item.state] as any}>{item.state}</Tag>
                        ) : (
                          ''
                        )}
                      </div>
                      <div className="pods-group-ctime memo">
                        创建于{moment(new Date(item.create_tm || '')).format('yyyy-MM-DD HH:mm:ss')}
                      </div>
                      <Button
                        onClick={(e) => {
                          e.stopPropagation();
                          handleShowEvents(item.events || []);
                        }}
                      >
                        Events
                      </Button>
                    </div>
                  }
                  key={item.name || index}
                >
                  <div className="pods-container">容器</div>
                  {item.containers?.map((containerItem) => {
                    return (
                      <div
                        className="pods-container-item"
                        key={`${item.name}-${containerItem.name}`}
                      >
                        <div className="container-item-icon">
                          <img src={DockerIcon} alt="docker"></img>
                        </div>
                        <div className="container-item-name column">
                          <div>
                            <OverflowToolTip
                              title={containerItem.name}
                              width="100%"
                            ></OverflowToolTip>
                          </div>
                          <div className="memo">
                            <OverflowToolTip
                              title={`镜像：${containerItem.image}`}
                              width="100%"
                            ></OverflowToolTip>
                          </div>
                        </div>
                        <div className="container-item-state column">
                          <div>{containerItem.state}</div>
                          <div className="memo">状态</div>
                        </div>

                        <div className="container-item-pords column">
                          {containerItem.ports && containerItem.ports.length > 0 ? (
                            <>
                              <div>
                                <OverflowToolTip
                                  title={getShowPords(containerItem.ports || [])}
                                  width="100%"
                                ></OverflowToolTip>
                              </div>
                              <div className="memo">端口</div>
                            </>
                          ) : (
                            ''
                          )}
                        </div>
                        <div className="container-item-terminal column">
                          <Button
                            type="link"
                            onClick={(e) => {
                              e.stopPropagation();
                              openWebTerminal(item, containerItem);
                            }}
                          >
                            打开
                          </Button>
                          <div className="memo">WebTerminal</div>
                        </div>
                        <Button
                          className="container-item-log"
                          onClick={(e) => {
                            e.stopPropagation();
                            handleShowLog(item, containerItem);
                          }}
                        >
                          日志
                        </Button>
                      </div>
                    );
                  })}
                </Collapse.Panel>
              );
            })}
          </Collapse>

          <AppContainerLog></AppContainerLog>
          <AppContainerGroupEvents></AppContainerGroupEvents>
        </>
      ) : (
        <div className="empty-container">
          <Empty msg="当前实例暂无 Pods"></Empty>
        </div>
      )}
    </div>
  );
};
export default ApplicationInstanceDetailPods;
