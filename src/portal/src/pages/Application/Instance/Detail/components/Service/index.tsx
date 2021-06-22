/*
 * @Author: liyuying
 * @Date: 2021-06-01 18:14:46
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 15:56:04
 * @Description: 实例详情--Service
 */
import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'umi';
import { ColumnProps } from 'antd/lib/table';
import { IAppServiceState, AppServiceAction } from './models/app-service';
import { ApplicationInstanceService } from '@/openApi/api';
import { CessBaseTable, Empty } from 'cess-ui';
interface Iprop {
  instance_name: string;
}
const ApplicationInstanceDetailService = ({ instance_name }: Iprop) => {
  const dispatch = useDispatch();
  const { serviceList }: IAppServiceState = useSelector((state: any) => state.appService);
  const columns: ColumnProps<ApplicationInstanceService>[] = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Type',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: 'Cluster_Ip',
      dataIndex: 'cluster_ip',
      key: 'cluster_ip',
    },
    {
      title: 'External_Ip',
      dataIndex: 'from',
      key: 'from',
      render: (value, record) => {
        return record.external_ips && record.external_ips.length > 0
          ? record.external_ips.join('，')
          : '<none>';
      },
    },
    {
      title: 'Ports',
      dataIndex: 'ports',
      key: 'ports',
      render: (value, record) => {
        const showPorts: string[] = [];
        (record.ports || []).forEach((element) => {
          let showPort = '';
          if (element.name) {
            showPort += `【${element.name}】`;
          }
          if (element.port) {
            showPort += element.port;
          }
          if (element.node_port) {
            showPort += `:${element.node_port}`;
          }
          if (element.protocol) {
            showPort += `/${element.protocol}`;
          }
          showPorts.push(showPort);
        });
        return showPorts.join('，');
      },
    },
  ];
  /**
   * 获取service信息
   */
  useEffect(() => {
    if (instance_name) {
      dispatch({
        type: AppServiceAction.GET_DATA,
        payload: instance_name,
      });
    }
  }, [instance_name, dispatch]);
  return (
    <CessBaseTable
      table={{
        maxHeight: 300,
        rowKey: (record: any): string => record.name,
        dataSource: serviceList,
        columns,
        renderEmpty: <Empty msg="当前实例暂无 Service" />,
      }}
    />
  );
};
export default ApplicationInstanceDetailService;
