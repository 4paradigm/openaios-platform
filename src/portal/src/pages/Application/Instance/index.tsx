import React, { useEffect, useMemo } from 'react';
import { ColumnProps } from 'antd/lib/table';
import { history, useDispatch, useSelector } from 'umi';
import { Button, Icon, CessBaseTable, Empty, Modal } from 'cess-ui';
import { APPLICATION_CHART, APPLICATION_INSTANCE_DETAIL } from '@/router/url';
import moment from 'moment';
import {
  IApplicationInstanceListState,
  ApplicationInstanceListAction,
} from '../models/application-instance-list';
import { ApplicationInstanceMetadata } from '@/openApi/api';

/*
 * @Author: liyuying
 * @Date: 2021-05-20 17:22:22
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-01 17:42:52
 * @Description: file content
 */
const breadcrumbList = [
  {
    name: '应用实例',
  },
];
const ApplicationPage = () => {
  const dispatch = useDispatch();
  const { dataSource, total, currentPage }: IApplicationInstanceListState = useSelector(
    (state: any) => state.applicationInstanceList,
  );
  const actions = useMemo(() => {
    return [
      <Button
        type="primary"
        key="create"
        onClick={() => {
          history.push(APPLICATION_CHART);
        }}
        icon={<Icon type="plus" />}
      >
        创建应用实例
      </Button>,
    ];
  }, []);

  const columns: ColumnProps<ApplicationInstanceMetadata>[] = useMemo(() => {
    return [
      {
        title: '实例名称',
        dataIndex: 'instance_name',
        key: 'instance_name',
        width: '301px',
        render: (value, record) => {
          return record.instance_name;
        },
      },
      {
        title: '实例状态',
        dataIndex: 'status',
        key: 'status',
        width: '120px',
        render: (value, record) => {
          return record.status;
        },
      },
      {
        title: '应用名称',
        dataIndex: 'chart_name',
        key: 'chart_name',
        width: '240px',
        render: (value, record) => {
          return record.chart_name;
        },
      },
      {
        title: '应用版本',
        dataIndex: 'chart_version',
        key: 'chart_version',
        width: '120px',
      },
      {
        title: '创建时间',
        dataIndex: 'create_tm',
        key: 'create_tm',
        width: '200px',
        render: (value, record) => {
          return moment(new Date(record.create_tm || '')).format('yyyy-MM-DD HH:mm:ss');
        },
      },
      {
        title: '操作',
        dataIndex: 'operate',
        key: 'operate',
        width: '90px',
        render: (value, record) => {
          return (
            <Button
              type="link"
              icon={<Icon type="delete" />}
              onClick={(e) => {
                e.stopPropagation();
                handelInstanceDelete(record);
              }}
            >
              删除
            </Button>
          );
        },
      },
    ];
  }, []);
  const pageChange = (currentPage: number) => {
    dispatch({
      type: ApplicationInstanceListAction.GET_LIST,
      payload: currentPage,
    });
  };
  /**
   * 删除实例
   * @param record
   */
  const handelInstanceDelete = (record: ApplicationInstanceMetadata) => {
    Modal.confirm({
      title: '删除应用实例',
      content: `确定要删除【${record.instance_name}】这个应用实例吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: () => {
        dispatch({
          type: ApplicationInstanceListAction.DELETE,
          payload: record.instance_name,
        });
      },
    });
  };

  useEffect(() => {
    dispatch({
      type: ApplicationInstanceListAction.GET_LIST,
      payload: currentPage,
    });
  }, [currentPage]);

  return (
    <div className="application-list">
      <CessBaseTable
        table={{
          rowKey: (record): string => record.instance_name,
          dataSource,
          total,
          columns,
          currentPage,
          onChange: (currentPage) => pageChange(currentPage),
          onRow: (record) => {
            history.push(`${APPLICATION_INSTANCE_DETAIL}${record.instance_name || ''}`);
          },
          unitMsg: '个应用实例',
          renderEmpty: (
            <Empty
              msg="您还没有应用实例"
              action={
                <Button
                  type="link"
                  onClick={() => {
                    history.push(APPLICATION_CHART);
                  }}
                >
                  创建应用实例
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

export default ApplicationPage;
