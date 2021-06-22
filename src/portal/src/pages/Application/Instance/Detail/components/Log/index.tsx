import React, { useEffect } from 'react';
import { Modal, InputNumber, Button } from 'cess-ui';
import { useSelector, useDispatch } from 'umi';
import { LazyLog, ScrollFollow } from 'react-lazylog';
import {
  LOG_DEAFAULT_TAIL_LINE,
  LOG_MIN_TAIL_LINE,
  LOG_MAX_TAIL_LINE,
} from '@/constant/application';
import { IAppContainerLogState, AppContainerLogAction } from './models/app-container-log';
import './index.less';

/*
 * @Author: liyuying
 * @Date: 2021-06-01 18:28:07
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-18 11:10:35
 * @Description: 展示容器的日志
 */

const AppContainerLog = () => {
  const dispatch = useDispatch();
  const {
    modalVisible,
    containerName,
    logContent,
    xmlHttp,
    tailLines,
  }: IAppContainerLogState = useSelector((state: any) => state.appContainerLog);
  /**
   * 停止获取log
   */
  const abortHttp = () => {
    if (xmlHttp) {
      xmlHttp.abort();
      dispatch({
        type: AppContainerLogAction.UPDATE_STATUS,
        payload: { xmlHttp: null },
      });
    }
  };
  /**
   * 设置Tail的大小
   */
  const setTailLines = (value: number) => {
    dispatch({
      type: AppContainerLogAction.UPDATE_STATUS,
      payload: { tailLines: value || LOG_DEAFAULT_TAIL_LINE },
    });
  };
  /**
   * 关闭弹窗
   */
  const handleCancel = () => {
    dispatch({
      type: AppContainerLogAction.UPDATE_STATUS,
      payload: { modalVisible: false, tailLines: LOG_DEAFAULT_TAIL_LINE, logContent: '' },
    });
    abortHttp();
  };
  const getLogData = () => {
    dispatch({
      type: AppContainerLogAction.UPDATE_STATUS,
      payload: { logContent: '' },
    });
    dispatch({
      type: AppContainerLogAction.GET_DATA,
      payload: {
        callBack: (logs: string) => {
          dispatch({
            type: AppContainerLogAction.UPDATE_STATUS,
            payload: { logContent: logs },
          });
        },
      },
    });
  };
  useEffect(() => {
    if (containerName && modalVisible) {
      getLogData();
    }
    return () => {
      abortHttp();
    };
  }, [containerName, modalVisible]);
  return (
    <Modal
      title="容器日志"
      visible={modalVisible}
      closable={true}
      footer={null}
      destroyOnClose={false}
      onCancel={handleCancel}
      centered
      maskClosable={false}
      className="app-container-log"
      width={880}
    >
      <div>
        <div className="log-search-bar">
          <label>
            Tail：
            <InputNumber
              value={tailLines}
              min={LOG_MIN_TAIL_LINE}
              max={LOG_MAX_TAIL_LINE}
              onChange={setTailLines}
            ></InputNumber>
          </label>
          <Button type="primary" onClick={getLogData}>
            重新获取日志
          </Button>
          <Button onClick={abortHttp} disabled={!xmlHttp}>
            停止获取日志
          </Button>
        </div>
      </div>
      <div className="log-container">
        <ScrollFollow
          startFollowing
          render={({ onScroll, follow }: any) => (
            <LazyLog
              extraLines={1}
              enableSearch
              selectableLines
              text={logContent || ` `}
              caseInsensitive
              stream
              onScroll={onScroll}
              follow={follow}
            />
          )}
        />
      </div>
    </Modal>
  );
};
export default AppContainerLog;
