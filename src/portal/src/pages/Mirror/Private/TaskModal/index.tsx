import React, { useMemo } from 'react';
import { Modal, CessBaseTable, Button } from 'cess-ui';
import { ColumnProps } from 'antd/lib/table';
import { useDispatch, useSelector } from 'umi';
import { PrivateImageAction } from '../../models/private-image';
import moment from 'moment';
import { Task } from '@/constant/task';
import './index.less';

function TaskModal() {
  const { taskModalVisible, taskList } = useSelector((state: any) => state.privateImage);
  const dispatch = useDispatch();

  const deleteTask = (importing_id: string) => {
    Modal.confirm({
      title: '提示',
      content: `确定要删除这个任务吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        dispatch({
          type: PrivateImageAction.DELETE_TASK,
          payload: importing_id,
        });
      },
    });
  };

  const stopTask = (importing_id: string) => {
    dispatch({
      type: PrivateImageAction.STOP_TASK,
      payload: importing_id,
    });
  };

  const handleCancle = () => {
    dispatch({ type: PrivateImageAction.CLOSE_MODAL });
  };

  const columns: ColumnProps<any>[] = useMemo(() => {
    return [
      {
        title: '仓库名',
        dataIndex: 'repo',
        key: 'repo',
      },
      {
        title: 'Tag',
        dataIndex: 'tag',
        key: 'tag',
      },
      {
        title: '状态',
        dataIndex: 'status',
        key: 'status',
      },
      {
        title: '创建时间',
        dataIndex: 'start_time',
        key: 'start_time',
        render: (value: any) => {
          return moment(new Date(value)).format('yyyy-MM-DD HH:mm:ss');
        },
      },
      {
        title: '结束时间',
        dataIndex: 'end_time',
        key: 'end_time',
        render: (value: any, record: any) => {
          return record.status === Task.Pending ||
            record.status === Task.InProgress ||
            /* 开始时间在结束时间之后 */
            new Date(record.start_time).getTime() >= new Date(record.end_time).getTime()
            ? '--'
            : moment(new Date(value)).format('yyyy-MM-DD HH:mm:ss');
        },
      },
      {
        title: 'Registry',
        dataIndex: 'registry',
        key: 'registry',
        render: (value: any) => {
          return value.url;
        },
      },
      {
        title: '操作',
        dataIndex: 'operate',
        key: 'operate',
        width: '130px',
        render: (value, record, index: number) => {
          return (
            <>
              <Button
                type="link"
                disabled={!(record.status === Task.InProgress || record.status === Task.Pending)}
                onClick={(e) => {
                  e.stopPropagation();
                  stopTask(record.importing_id);
                }}
              >
                停止
              </Button>
              <Button
                type="link"
                disabled={record.status === Task.InProgress}
                onClick={(e) => {
                  e.stopPropagation();
                  deleteTask(record.importing_id);
                }}
              >
                删除
              </Button>
            </>
          );
        },
      },
    ];
  }, []);

  return (
    <Modal
      title="任务列表"
      visible={taskModalVisible}
      closable={true}
      footer={null}
      className="task-modal"
      destroyOnClose={true}
      centered
      width={950}
      onCancel={handleCancle}
    >
      <CessBaseTable
        table={{
          rowKey: 'importing_id',
          dataSource: taskList,
          columns,
        }}
      />
    </Modal>
  );
}

export default TaskModal;
