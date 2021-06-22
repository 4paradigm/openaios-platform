import React, { useMemo } from 'react';
import { useDispatch, useSelector } from 'umi';
import { ColumnProps } from 'antd/lib/table';
import { ApplicationInstanceEvent } from '@/openApi/api';
import {
  AppContainerGroupEventsAction,
  IAppContainerGroupEventsState,
} from './models/app-container-group-events';
import { Modal, CessBaseTable, Empty } from 'cess-ui';
import './index.less';

/*
 * @Author: liyuying
 * @Date: 2021-06-01 20:00:04
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-03 11:35:13
 * @Description: 容器组Events
 */
const AppContainerGroupEvents = () => {
  const dispatch = useDispatch();
  const { modalVisible, eventList }: IAppContainerGroupEventsState = useSelector(
    (state: any) => state.appContainerGroupEvents,
  );
  const handleCancel = () => {
    dispatch({
      type: AppContainerGroupEventsAction.UPDATE_STATUS,
      payload: { modalVisible: false },
    });
  };
  const columns: ColumnProps<ApplicationInstanceEvent>[] = [
    {
      title: 'Type',
      dataIndex: 'type',
      key: 'type',
      width: '120px',
    },
    {
      title: 'Reason',
      dataIndex: 'reason',
      key: 'reason',
      width: '120px',
    },
    {
      title: 'Age',
      dataIndex: 'age',
      key: 'age',
      width: '140px',
    },
    {
      title: 'From',
      dataIndex: 'from',
      key: 'from',
      width: '120px',
    },
    {
      title: 'Message',
      dataIndex: 'message',
      key: 'message',
    },
  ];

  return (
    <Modal
      title="Events"
      visible={modalVisible}
      closable={true}
      footer={null}
      destroyOnClose={false}
      onCancel={handleCancel}
      centered
      className="app-container-group-events"
      width={880}
    >
      <CessBaseTable
        table={{
          maxHeight: 300,
          rowKey: (record: any): string => record.age,
          dataSource: eventList,
          columns,
          unitMsg: '个Event',
          renderEmpty: <Empty msg="当前容器组无Event" />,
        }}
      />
    </Modal>
  );
};
export default AppContainerGroupEvents;
