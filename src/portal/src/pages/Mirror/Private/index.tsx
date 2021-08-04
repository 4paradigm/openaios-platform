import React, { useMemo, useEffect } from 'react';
import { CopyOutlined } from '@ant-design/icons';
import { CessBaseTable, Button, Icon, Modal, OverflowToolTip } from 'cess-ui';
import { useSelector, useDispatch } from 'react-redux';
import { ColumnProps } from 'antd/lib/table';
import { PrivateImageAction, IMirrorPrivateState } from '../models/private-image';
import moment from 'moment';
import ImportModal from './ImportModal';
import TaskModal from './TaskModal';
import './index.less';
import { LoadingOutlined } from '@ant-design/icons';
import { Task } from '@/constant/task';
import CopyImageModal from './CopyImageModal';
import { ImageInfo } from '@/openApi/api';

const breadcrumbList = [
  {
    name: '私有镜像仓库',
  },
];
const MirrorPrivate = () => {
  const dispatch = useDispatch();
  const {
    dataSource,
    total,
    currentPage,
    taskList,
    copyImageLoading,
  }: IMirrorPrivateState = useSelector((state: any) => state.privateImage);
  useEffect(() => {
    dispatch({ type: PrivateImageAction.INITT_DATA });
  }, []);

  const importImage = () => {
    dispatch({
      type: PrivateImageAction.OPEN_IMAGE_MODAL,
    });
  };

  const deleteImage = (record: ImageInfo) => {
    Modal.confirm({
      title: '提示',
      content: `确定要删除这个镜像吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        dispatch({
          type: PrivateImageAction.DELETE_IMAGE,
          payload: {
            repo: record.repo,
            digest: record.digest,
          },
        });
      },
    });
  };
  /**
   * 拷贝镜像
   * @param record
   */
  const copyImage = (record: ImageInfo, loading: boolean) => {
    if (loading) {
      Modal.info({
        title: '镜像拷贝中...',
        content: `请等待上一镜像拷贝完成后再进行操作`,
        okText: '确认',
      });
    } else {
      dispatch({
        type: PrivateImageAction.UPDATE_STATE,
        payload: {
          copyImageModalVisible: true,
          copyImageSourceImage: record,
        },
      });
    }
  };

  const showTaskModal = () => {
    dispatch({
      type: PrivateImageAction.UPDATE_STATE,
      payload: {
        taskModalVisible: true,
      },
    });
  };

  const pageChange = (currentPage: number, pageSize: number) => {
    dispatch({
      type: PrivateImageAction.GET_LIST,
      payload: currentPage,
    });
  };

  const actions = useMemo(() => {
    return [
      <Button type="primary" className="task" key="task" onClick={showTaskModal}>
        任务列表
        {taskList.filter((data: any) => data.status === Task.InProgress).length > 0 && (
          <LoadingOutlined />
        )}
      </Button>,
      <Button type="primary" className="import" key="import" onClick={importImage}>
        导入
      </Button>,
    ];
  }, [taskList]);

  const columns: ColumnProps<ImageInfo>[] = useMemo(() => {
    return [
      {
        title: '仓库名',
        dataIndex: 'repo',
        key: 'repo',
        width: '300px',
      },
      {
        title: 'Tags',
        dataIndex: 'tags',
        key: 'tags',
        width: '300px',
        render: (value: any) => {
          return value.join(',');
        },
      },
      {
        title: '大小',
        dataIndex: 'size',
        key: 'size',
        width: '120px',
        render: (value, record) => {
          return <OverflowToolTip title={value} width={80} lineHeight={16}></OverflowToolTip>;
        },
      },
      {
        title: '导入时间',
        dataIndex: 'importing_time',
        key: 'importing_time',
        width: '220px',
        render: (value: any) => {
          return moment(new Date(value)).format('yyyy-MM-DD HH:mm:ss');
        },
      },
      {
        title: '操作',
        dataIndex: 'operate',
        key: 'operate',
        width: '180px',
        render: (value, record, index: number) => {
          return (
            <>
              <Button
                type="link"
                icon={<Icon type="delete" />}
                onClick={(e) => {
                  e.stopPropagation();
                  deleteImage(record);
                }}
              >
                删除
              </Button>
              {record.tags && record.tags.length > 0 ? (
                <Button
                  type="link"
                  icon={<CopyOutlined />}
                  onClick={(e) => {
                    e.stopPropagation();
                    copyImage(record, copyImageLoading);
                  }}
                >
                  拷贝
                </Button>
              ) : (
                ''
              )}
            </>
          );
        },
      },
    ];
  }, [copyImageLoading]);

  return (
    <div className="private-image">
      <CessBaseTable
        table={{
          rowKey: (record) => {
            return record.repo + record.tags + record.size;
          },
          dataSource,
          total,
          columns,
          currentPage,
          onChange: (currentPage, pageSize) => pageChange(currentPage, pageSize),
        }}
        actions={actions}
        breadcrumbList={breadcrumbList}
      />
      <ImportModal />
      <TaskModal />
      <CopyImageModal />
    </div>
  );
};

export default MirrorPrivate;
