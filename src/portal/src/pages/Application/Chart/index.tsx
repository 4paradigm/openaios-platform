/*
 * @Author: liyuying
 * @Date: 2021-05-21 16:34:38
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 17:41:13
 * @Description: 应用商店
 */
import React, { useEffect } from 'react';
import { history, useDispatch, useSelector } from 'umi';
import { Breadcrumb, Icon, Avatar, OverflowToolTip } from 'cess-ui';
import { APPLICATION_INSTANCE_CREATE, APPLICATION_CHART_CREATE } from '@/router/url';
import { IApplicationChartState, ApplicationChartAction } from '../models/application-chart';
import { ChartMetadata } from '@/openApi/api';
import {
  CREATE_MY_APP_KEY,
  CHART_CATEGORY_LIST,
  CHART_CATEGORY_PUBLIC,
} from '@/constant/application';
import './index.less';
import Loading from '@/components/Loading';

const breadcrumb = (
  <Breadcrumb>
    <Breadcrumb.Item>应用市场</Breadcrumb.Item>
  </Breadcrumb>
);
const ApplicationChart = () => {
  const dispatch = useDispatch();
  const { dataSource, publicChartCount, isLoading }: IApplicationChartState = useSelector(
    (state: any) => state.applicationChart,
  );
  /**
   * 创建实例
   * @param item
   */
  const handleCreateInstance = (item: ChartMetadata) => {
    history.push(
      `${APPLICATION_INSTANCE_CREATE}/${item.name}?version=${item.version}&category=${item.category}`,
    );
  };
  /**
   * 创建我的应用
   */
  const handleCreateMyApp = () => {
    history.push(APPLICATION_CHART_CREATE);
  };
  const ChartItem = (categoryName: string, charts: ChartMetadata[]) => {
    return (
      <div className="application-chart-group" key={categoryName}>
        {categoryName ? (
          <div className="application-chart-group-sub-title">{categoryName}</div>
        ) : (
          <></>
        )}
        {charts.map((item) => {
          return item.name === CREATE_MY_APP_KEY ? (
            <div
              className="application-chart-item my-app-create"
              key="my-create"
              onClick={() => {
                handleCreateMyApp();
              }}
            >
              <Icon type="plus-circle" />
              <label>创建我的应用</label>
            </div>
          ) : (
            <div
              className="application-chart-item"
              key={`${item.category}-${item.name}-${item.version}`}
              onClick={() => {
                handleCreateInstance(item);
              }}
            >
              <Avatar
                className="application-chart-item-icon"
                icon={<Icon type="application" />}
                src={item.icon_link}
              />
              <div className="application-chart-item-name">
                <OverflowToolTip
                  title={item.showName || item.name}
                  line={2}
                  lineHeight={25}
                  width="100%"
                ></OverflowToolTip>
              </div>
              <div className="application-chart-item-version">
                {item.version}
                <div className="version-divider"></div>
              </div>
              <div className="application-chart-item-description">
                <OverflowToolTip title={item.description} width="100%" line={2}></OverflowToolTip>
              </div>
            </div>
          );
        })}
      </div>
    );
  };
  useEffect(() => {
    dispatch({
      type: ApplicationChartAction.GET_LIST,
    });
  }, []);
  return (
    <div className="application-chart comm-create-page">
      {breadcrumb}
      {isLoading ? (
        <Loading></Loading>
      ) : (
        <div className="application-chart-container">
          {CHART_CATEGORY_LIST.map((chartCategory) => {
            // 无内置应用
            if (publicChartCount === 0 && chartCategory.category === CHART_CATEGORY_PUBLIC) {
              return <div key={chartCategory.category}></div>;
            }
            return (
              <div key={chartCategory.category}>
                <>
                  <div className="application-chart-group-title">{chartCategory.categoryName}</div>
                  {chartCategory.subCategorys
                    ? chartCategory.subCategorys.map((subCategory) => {
                        if (
                          dataSource[subCategory.category] &&
                          dataSource[subCategory.category].length > 0
                        ) {
                          return ChartItem(
                            subCategory.categoryName,
                            (dataSource[subCategory.category] as any) || [],
                          );
                        }
                        return <></>;
                      })
                    : ChartItem('', (dataSource[chartCategory.category] as any) || [])}
                </>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};
export default ApplicationChart;
