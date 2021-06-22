/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-31 17:28:27
 * @Description: file content
 */
import { Model } from 'dva';
import { UserApiInstance, CompetitionApiApiInstance } from '@/openApi';
import { IAction } from '@/interfaces';
import { IMessage } from '@/interfaces/bussiness';
import { delay, toThousands } from '@/utils';

export const HomeAction = {
  GET_INFO_TASK: 'home/getInfoTask',
  GET_USER_INFO: 'home/getUserInfo',
  GET_MESSAGE_INFO: 'home/getMessageInfo',
};

export interface IHomeState {
  userInfo: any;
  taskInfo: any;
  msgInfo: IMessage[];
}

const defaultState: IHomeState = {
  userInfo: null,
  taskInfo: null,
  msgInfo: [],
};

const home: Model = {
  namespace: 'home',
  state: defaultState,
  effects: {
    *getInfoTask({ payload }, { call, put, select }) {
      const { data }: any = yield call(UserApiInstance.userTasksGet.bind(UserApiInstance));
      if (data) {
        // 每分钟消耗
        data.perCost = 0;
        data.task_list.forEach((item: any) => {
          data.perCost = data.perCost + item.price * item.number;
        });
        data.perCost = toThousands(data.perCost, 3);
        yield put({
          type: 'updateState',
          payload: {
            taskInfo: data,
          },
        });
      }
    },
    *getUserInfo({ payload }, { call, put, select }) {
      const { data }: any = yield call(UserApiInstance.getUser.bind(UserApiInstance));
      const { taskInfo } = yield select((state: { home: IHomeState }) => state.home);
      if (data) {
        data.costTime = data.balance ? '+ ∞' : '0 min';
        if (taskInfo && taskInfo.perCost && data.balance) {
          if (Number(taskInfo.perCost)) {
            const minCount = Math.floor(data.balance / Number(taskInfo.perCost));
            const hour = Math.floor(minCount / 60);
            const min = Math.floor(minCount % 60);
            data.costTime = `${hour}h ${min}min`;
          } else {
            data.costTime = '+ ∞';
          }
        }
        data.balance = toThousands(data.balance, 3);
        yield put({
          type: 'updateState',
          payload: {
            userInfo: data,
          },
        });
        // 一分钟刷新一次余额
        if (window.location.pathname === '/home') {
          yield call(delay, 60 * 1000);
          yield put({
            type: 'getUserInfo',
          });
        }
      }
    },
    *getMessageInfo({ payload }, { call, put }) {
      const { data }: any = yield call(
        CompetitionApiApiInstance.competitionGet.bind(CompetitionApiApiInstance),
      );
      if (data) {
        const messageList: any[] = [];
        data.forEach((element: IMessage) => {
          if (element.id) {
            element.title = `【比赛】${element.name}`;
            try {
              element.descriptionMd = require(`../data/games/${element.id}/README.md`).default;
              // 初始化环境的配置
              const initEnvJson = require(`../data/games/${element.id}/INIT.json`);
              if (initEnvJson) {
                element.initEnvJson = initEnvJson;
              }
              // 调查问卷
              const queryFormn = require(`../data/games/${element.id}/FORM.json`);
              if (queryFormn && queryFormn.form) {
                element.formConfig = queryFormn.form;
              }
            } catch (error) {}
            messageList.push(element);
          }
        });
        yield put({
          type: 'updateState',
          payload: {
            msgInfo: messageList,
          },
        });
      }
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

export default home;
