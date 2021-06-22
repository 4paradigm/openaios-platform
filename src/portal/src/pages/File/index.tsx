import React, { useEffect, useMemo } from 'react';
import {
  Button,
  Icon,
  CessBaseTable,
  Modal,
  message,
  CessDropdown,
  Menu,
  OverflowToolTip,
} from 'cess-ui';
import { FolderFilled, HomeOutlined } from '@ant-design/icons';
import { useDispatch, useSelector } from 'react-redux';
import { ColumnProps } from 'antd/lib/table';
import { IFileListState, FileListActions } from './models/file-list';
import CreateFolderModal from './CreateFolderModal';
import UploadFileModal from './UploadFileModal';
import moment from 'moment';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { arrSplit, getFileIcon } from '@/utils';

import './index.less';

const breadcrumbList = [
  {
    name: '文件管理',
  },
];

function File() {
  const dispatch = useDispatch();
  useEffect(() => {
    dispatch({
      type: FileListActions.GET_LIST,
    });
  }, [dispatch]);

  const { dataSource, filePath, downIndex }: IFileListState = useSelector(
    (state: any) => state.fileList,
  );

  const createFolder = () => {
    dispatch({
      type: FileListActions.UPDATE_STATUS,
      payload: {
        folderModalVisible: true,
      },
    });
  };

  const uploadFile = () => {
    dispatch({
      type: FileListActions.UPDATE_STATUS,
      payload: {
        fileModalVisible: true,
      },
    });
  };

  const deleteFile = (name: string) => {
    Modal.confirm({
      title: '提示',
      content: `确定要删除【${name}】这个文件吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        dispatch({
          type: FileListActions.DELETE_FILE,
          payload: {
            path: `${filePath}${name}`,
          },
        });
      },
    });
  };

  const downLoad = (name: string, index: number) => {
    dispatch({
      type: FileListActions.DOWNLOAD_FILE,
      payload: {
        filePath,
        name,
        index,
      },
    });
  };

  const handleCopy = (path: string) => {
    message.success(`${path}，已复制到剪贴板`);
  };

  const toDir = (path: string) => {
    dispatch({
      type: FileListActions.CHANGE_PATH,
      payload: {
        path,
      },
    });
  };

  const changePath = (e: any) => {
    dispatch({
      type: FileListActions.CHANGE_PATH,
      payload: {
        path: e.key,
      },
    });
  };

  const actions = useMemo(() => {
    return [
      <Button type="default" key="upload" onClick={uploadFile}>
        上传文件
      </Button>,
      <Button type="primary" key="create" icon={<Icon type="plus" />} onClick={createFolder}>
        新建文件夹
      </Button>,
    ];
  }, []);

  const filePathShow = useMemo(() => {
    let pathArr: string[] = filePath.split('/').filter((s) => s);
    const [ellPath, showPath] = arrSplit(pathArr);
    return (
      <div className="path-breadcrumb">
        <span>
          当前目录：
          <span className="pathName" onClick={() => toDir('/')}>
            <HomeOutlined />
          </span>
        </span>
        {ellPath.length > 0 && (
          <CessDropdown
            placement="bottomRight"
            overlay={
              <Menu onClick={(e) => changePath(e)}>
                {ellPath.map((path: string, index: number) => {
                  const key = ellPath.slice(0, index + 1).join('/');
                  return <Menu.Item key={`/${key}/`}>{path}</Menu.Item>;
                })}
              </Menu>
            }
          >
            <span className="path-divider">/</span>
            <span className="pathName">...</span>
          </CessDropdown>
        )}
        {showPath.map((name: any, index: number) => {
          const path = showPath.slice(0, index + 1).join('/');
          const prePath = ellPath.join('/');
          return (
            <span key={path}>
              <span className="path-divider">/</span>
              <span
                className="pathName"
                onClick={() => toDir(`${prePath ? '/' + prePath : ''}/${path}/`)}
              >
                {name}
              </span>
            </span>
          );
        })}
      </div>
    );
  }, [filePath]);

  const columns: ColumnProps<any>[] = useMemo(() => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        key: 'name',
        render: (value, record) => {
          return (
            <div className="column-name">
              {record.is_dir && (
                <div className="fileName" onClick={() => toDir(`${filePath}${record.name}/`)}>
                  <FolderFilled />
                  <OverflowToolTip title={value} width={180} lineHeight={16}></OverflowToolTip>
                </div>
              )}
              {!record.is_dir && (
                <>
                  <Icon type={getFileIcon(value)} />
                  <OverflowToolTip title={value} width={180} lineHeight={16}></OverflowToolTip>
                </>
              )}
            </div>
          );
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
        title: '创建时间',
        dataIndex: 'modification_time',
        key: 'modification_time',
        width: '220px',
        render: (value, record) => {
          return moment(new Date(value)).format('yyyy-MM-DD HH:mm:ss');
        },
      },
      {
        title: '操作',
        dataIndex: 'operate',
        key: 'operate',
        width: '250px',
        render: (value, record, index: number) => {
          return (
            <>
              <Button
                disabled={record.is_dir || downIndex === index}
                type="link"
                onClick={(e) => {
                  e.stopPropagation();
                  downLoad(record.name, index);
                }}
                icon={<Icon type="download" />}
              >
                下载
              </Button>
              <Button
                type="link"
                icon={<Icon type="delete" />}
                onClick={(e) => {
                  e.stopPropagation();
                  deleteFile(record.name);
                }}
              >
                删除
              </Button>
              <CopyToClipboard
                text={filePath + record.name}
                onCopy={() => handleCopy(filePath + record.name)}
              >
                <Button
                  onClick={(e) => {
                    e.stopPropagation();
                  }}
                  type="link"
                  icon={<Icon type="save" />}
                >
                  复制
                </Button>
              </CopyToClipboard>
            </>
          );
        },
      },
    ];
  }, [filePath, downIndex]);

  return (
    <div className="file">
      <CessBaseTable
        table={{
          rowKey: 'name',
          dataSource,
          columns,
        }}
        titleNav={{
          others: filePathShow,
          actions,
        }}
        breadcrumbList={breadcrumbList}
      />
      <CreateFolderModal />
      <UploadFileModal />
    </div>
  );
}

export default File;
