import React, { useMemo, useEffect } from 'react';
import { CessBaseTable, Button, Icon, Steps, Empty, Tooltip, Modal, Tag } from 'cess-ui';
import { useSelector, useDispatch } from 'react-redux';
import { ColumnProps } from 'antd/lib/table';
import { history } from 'umi';
import { IEnvironmentListState, EnvironmentListAction } from './models/environment-list';
import { ENVIRONMENT_CREATE } from '@/router/url';
import './index.less';
import moment from 'moment';
import juperterIcon from '@/assets/images/jupter.png';
import sshIcon from '@/assets/images/ssh.png';
import webTerminalIcon from '@/assets/images/web-terminal.png';
import { ENVIRONMENT_STATUS, ENVIRONMENT_STATUS_TAG } from '@/constant/environment';

const { Step } = Steps;

const breadcrumbList = [
  {
    name: '开发环境',
  },
];
const DevEnvironment = () => {
  const dispatch = useDispatch();
  const { dataSource, total, currentPage }: IEnvironmentListState = useSelector(
    (state: any) => state.environmentList,
  );

  useEffect(() => {
    dispatch({
      type: EnvironmentListAction.GET_LIST,
      payload: currentPage,
    });
  }, []);
  const pageChange = (currentPage: number, pageSize: number) => {
    dispatch({
      type: EnvironmentListAction.GET_LIST,
      payload: currentPage,
    });
    dispatch({
      type: EnvironmentListAction.UPDATE_STATUS,
      payload: currentPage,
    });
  };
  const delteEnvironment = (name: string) => {
    Modal.confirm({
      title: '提示',
      content: `确定要删除【${name}】这个环境吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        dispatch({
          type: EnvironmentListAction.DELETE_ENVIRONMENT,
          payload: name,
        });
      },
    });
  };

  const openWebTerminal = (pod: string) => {
    dispatch({
      type: EnvironmentListAction.OPEN_TERMINAL,
      payload: pod,
    });
  };

  const actions = useMemo(() => {
    return [
      <Tooltip
        key="tooltip"
        placement="left"
        title={
          <Steps current={2} hiddenDivider={true}>
            <Step key="step1" title="创建一个开发环境"></Step>
            <Step key="step2" title="启用开发环境"></Step>
          </Steps>
        }
      >
        <span>
          <Icon type="info-outline" />
        </span>
      </Tooltip>,
      <Button
        type="primary"
        key="create"
        onClick={() => {
          history.push(ENVIRONMENT_CREATE);
        }}
        icon={<Icon type="plus" />}
      >
        创建开发环境
      </Button>,
    ];
  }, []);

  const columns: ColumnProps<any>[] = useMemo(() => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        key: 'name',
        render: (value, record: any) => {
          return record.staticInfo.name;
        },
      },
      {
        title: '状态',
        dataIndex: 'state',
        key: 'state',
        render: (value, record: any) => {
          return (
            <>
              <Tag color={ENVIRONMENT_STATUS_TAG[value] as any}>{value}</Tag>
              {value === ENVIRONMENT_STATUS.Unknown || value === ENVIRONMENT_STATUS.Killed ? (
                <Tooltip title={record.staticInfo.description}>
                  <Icon type="info-outline" />
                  &nbsp;
                </Tooltip>
              ) : (
                ''
              )}
            </>
          );
        },
      },
      {
        title: '交互方式',
        dataIndex: 'interact',
        key: 'interact',
        width: '216px',
        render: (value, record: any) => {
          return (
            <div className="interact-td">
              {record.staticInfo.environmentConfig.jupyter.enable && (
                <span className="item">
                  <Tooltip
                    key="tooltip"
                    placement="top"
                    title={record.staticInfo.environmentConfig.jupyter.token}
                  >
                    <img src={juperterIcon} alt="juperty" />
                  </Tooltip>
                  {record.state === ENVIRONMENT_STATUS.Running && (
                    <a href={record.staticInfo.notebook_url} target="_blank">
                      打开
                    </a>
                  )}
                </span>
              )}
              {record.staticInfo.environmentConfig.ssh.enable && (
                <span className="item">
                  <Tooltip
                    key="tooltip"
                    placement="top"
                    title={Object.keys(record.sshInfo).map((key) => {
                      return (
                        <p key={key}>
                          {key}：{record.sshInfo[key]}
                        </p>
                      );
                    })}
                  >
                    <img src={sshIcon} alt="ssh" />
                  </Tooltip>
                </span>
              )}
              <span className="item">
                <img src={webTerminalIcon} alt="terminal" />
                {record.state === ENVIRONMENT_STATUS.Running && (
                  <Button
                    type="link"
                    onClick={() => {
                      openWebTerminal(record.pod_name);
                    }}
                  >
                    打开
                  </Button>
                )}
              </span>
            </div>
          );
        },
      },
      {
        title: '算力规格',
        dataIndex: 'compute_unit',
        key: 'compute_unit',
        render: (value, record: any) => {
          return record.staticInfo.environmentConfig.compute_unit;
        },
      },
      {
        title: '创建时间',
        dataIndex: 'create_tm',
        key: 'create_tm',
        render: (value, record: any) => {
          return moment(new Date(record.staticInfo.create_tm)).format('yyyy-MM-DD HH:mm:ss');
        },
      },
      {
        title: '操作',
        dataIndex: 'operate',
        key: 'operate',
        width: '80px',
        render: (value, record, index: number) => {
          return (
            <Button
              type="link"
              icon={<Icon type="delete" />}
              onClick={(e) => {
                e.stopPropagation();
                delteEnvironment(record.staticInfo.name);
              }}
            >
              删除
            </Button>
          );
        },
      },
    ];
  }, []);

  return (
    <div className="dev-environment">
      <CessBaseTable
        table={{
          rowKey: (record: any): string => record.staticInfo.name,
          dataSource,
          total,
          columns,
          currentPage,
          onChange: (currentPage, pageSize) => pageChange(currentPage, pageSize),
          renderEmpty: (
            <Empty
              msg="您还没有云上开发环境"
              action={
                <Button
                  type="link"
                  onClick={() => {
                    history.push(ENVIRONMENT_CREATE);
                  }}
                >
                  创建环境
                </Button>
              }
            />
          ),
        }}
        actions={actions}
        breadcrumbList={breadcrumbList}
      />
    </div>
  );
};

export default DevEnvironment;
