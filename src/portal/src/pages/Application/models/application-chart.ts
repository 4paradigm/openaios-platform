/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-21 13:20:23
 * @Description: file content
 */
import { Model } from 'dva';
import { AppstoreApiInstance } from '@/openApi';
import { IAction, IChartList } from '@/interfaces';
// import { templateList } from '../mock/template';
import { ChartCategory, Chart, ChartMetadata } from '@/openApi/api';
import { CREATE_MY_APP_KEY } from '@/constant/application';

export const ApplicationChartAction = {
  UPDATE_STATUS: 'applicationChart/updateStatus',
  GET_LIST: 'applicationChart/getList',
};

export interface IApplicationChartState {
  dataSource: IChartList;
  publicChartCount: number;
  isLoading: boolean;
}

const defaultState: IApplicationChartState = {
  dataSource: {},
  publicChartCount: 0,
  isLoading: true,
};

const applicationChart: Model = {
  namespace: 'applicationChart',
  state: defaultState,
  effects: {
    *getList({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          isLoading: true,
        },
      });
      const { data } = yield call(
        AppstoreApiInstance.getAppstoreChartList.bind(AppstoreApiInstance),
      );
      const chartMap: IChartList = {};
      let publicChartCount = 0;
      if (data) {
        const templateList = data.items;
        templateList.forEach((item: ChartMetadata) => {
          const originDesc = item.description;
          if (originDesc?.startsWith('[[')) {
            const showNameIndex = originDesc?.indexOf(']]:');
            item.showName = originDesc.substring(2, showNameIndex);
            item.description = originDesc.substring(showNameIndex + 3, originDesc.length);
          }
          if (item.category !== ChartCategory.Private) {
            publicChartCount = publicChartCount + 1;
          }
          // 为官方应用排序（强制）
          if (item.category === ChartCategory.PublicOfficial) {
            if (item.name === 'openmldb') {
              (item as any).position = 1;
            } else if (item.name === 'openembedding-dev') {
              (item as any).position = 2;
            } else if (item.name === 'pafka') {
              (item as any).position = 3;
            } else {
              (item as any).position = 999;
            }
          }
          if (chartMap[item.category || ChartCategory.Private]) {
            (chartMap[item.category || ChartCategory.Private] || []).push(item);
          } else {
            chartMap[item.category || ChartCategory.Private] = [item];
          }
        });
      }
      if (!chartMap[ChartCategory.Private]) {
        chartMap[ChartCategory.Private] = [];
      }
      (chartMap[ChartCategory.Private] || []).push({ name: CREATE_MY_APP_KEY });
      // 针对官方应用强制排序
      (chartMap[ChartCategory.PublicOfficial] || []).sort((item1: any, item2: any) => {
        return item1.position - item2.position;
      });
      yield put({
        type: 'updateState',
        payload: {
          dataSource: chartMap,
          publicChartCount,
          isLoading: false,
        },
      });
    },
  },
  reducers: {
    updateState(state, { payload }: IAction) {
      return {
        ...state,
        ...payload,
      };
    },
  },
};

export default applicationChart;
