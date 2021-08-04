import React, { useMemo } from 'react';
import { useDispatch, useSelector, EnvironmentListAction, IEnvironmentListState } from 'umi';
import { ColumnProps } from 'antd/lib/table';
import { ApplicationInstanceEvent } from '@/openApi/api';
import { Modal, CessBaseTable, Empty } from 'cess-ui';
import './index.less';

/*
 * @Author: liyuying
 * @Date: 2021-06-01 20:00:04
 * @LastEditors: liyuying
 * @LastEditTime: 2021-07-05 20:51:41
 * @Description: 开发环境Events
 */
const DevEnvironmentEvents = () => {
  const dispatch = useDispatch();
  const { eventVisible, eventList }: IEnvironmentListState = useSelector(
    (state: any) => state.environmentList,
  );
  const handleCancel = () => {
    dispatch({
      type: EnvironmentListAction.UPDATE_STATUS,
      payload: { eventVisible: false },
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
      visible={eventVisible}
      closable={true}
      footer={null}
      destroyOnClose={false}
      onCancel={handleCancel}
      centered
      className="dev-enviromrnt-events"
      width={880}
    >
      <CessBaseTable
        table={{
          maxHeight: 300,
          rowKey: (record: any): string => record.age + record.message,
          dataSource: eventList,
          columns,
          unitMsg: '个Event',
          renderEmpty: <Empty msg="当前开发环境无Event" />,
        }}
      />
    </Modal>
  );
};
export default DevEnvironmentEvents;
