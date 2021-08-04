/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-19 16:10:00
 * @Description: file content
 */
import React, { useMemo, useEffect } from 'react';
import { CessBaseTable, OverflowToolTip } from 'cess-ui';
import { useSelector, useDispatch } from 'react-redux';
import { ColumnProps } from 'antd/lib/table';
import { PublicImageAction, IMirrorPublicState } from '../models/public-image';
import moment from 'moment';

const breadcrumbList = [
  {
    name: '公有镜像仓库',
  },
];
const MirrorPublic = () => {
  const dispatch = useDispatch();
  const { dataSource, total, currentPage }: IMirrorPublicState = useSelector(
    (state: any) => state.publicImage,
  );

  useEffect(() => {
    dispatch({
      type: PublicImageAction.GET_LIST,
      payload: currentPage,
    });
    dispatch({
      type: PublicImageAction.GET_TOTAL,
    });
  }, []);

  const columns: ColumnProps<any>[] = useMemo(() => {
    return [
      {
        title: '仓库名',
        dataIndex: 'repo',
        key: 'repo',
        width: '390px',
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
    ];
  }, []);
  const pageChange = (currentPage: number, pageSize: number) => {
    dispatch({
      type: PublicImageAction.GET_LIST,
      payload: currentPage,
    });
  };
  return (
    <div className="public-image">
      <CessBaseTable
        table={{
          rowKey: (record) => {
            return record.repo + record.tags;
          },
          dataSource,
          total,
          columns,
          currentPage,
          onChange: (currentPage, pageSize) => pageChange(currentPage, pageSize),
        }}
        breadcrumbList={breadcrumbList}
      />
    </div>
  );
};

export default MirrorPublic;
