/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-10 14:47:00
 * @Description: file content
 */
import { Model } from 'dva';
import { history } from 'umi';
import { CompetitionApiApiInstance, EnvironmentsApiInstance } from '@/openApi';
import { IAction } from '@/interfaces';
import { IMessage } from '@/interfaces/bussiness';
import { message } from 'cess-ui';
import { DEV_ENVIRONMENT } from '@/router/url';
import { delay } from '@/utils';

export const MessageDetailAction = {
  INIT_DATA: 'messageDetail/initData',
  SIGN_UP: 'messageDetail/signUp',
  INIT_ENV: 'messageDetail/initEnv',
  UPDATE_STATE: 'messageDetail/updateState',
};

export interface IMessageDetailState {
  msgInfo: IMessage;
  modalLoading: boolean;
  /* 是否已报名 */
  hasApplied: boolean;
  invitePersonNum: number;
}

const defaultState: IMessageDetailState = {
  msgInfo: { id: 0, name: '', descriptionMd: '', title: '', avl: false },
  modalLoading: false,
  hasApplied: false,
  invitePersonNum: 0,
};

const messageDetail: Model = {
  namespace: 'messageDetail',
  state: defaultState,
  effects: {
    *initData({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: defaultState,
      });
      // 获取当前数据
      yield put({
        type: 'getMessageInfo',
        payload,
      });
      yield put({
        type: 'checkStatus',
        payload,
      });
      yield put({
        type: 'getInvitePersonNum',
        payload,
      });
    },
    *getMessageInfo({ payload }, { call, put }) {
      const { data }: any = yield call(
        CompetitionApiApiInstance.competitionGet.bind(CompetitionApiApiInstance),
      );
      if (data) {
        let message = {};
        data.forEach((element: IMessage) => {
          // 当前查看的对象
          if (element.id && payload === element.id) {
            element.title = `【比赛】${element.name}`;
            try {
              element.descriptionMd = require(`../data/games/${element.id}/README.md`).default;
            } catch (error) {}
            try {
              // 初始化环境的配置
              const initEnvJson = require(`../data/games/${element.id}/INIT.json`);
              if (initEnvJson) {
                element.initEnvJson = initEnvJson;
              }
            } catch (error) {}
            try {
              // 调查问卷
              const queryFormn = require(`../data/games/${element.id}/FORM.json`);
              if (queryFormn && queryFormn.form) {
                element.formConfig = queryFormn.form;
              }
            } catch (error) {}
            try {
              // 邀请规则
              element.ruleMd = require(`../data/games/${element.id}/RULE.md`).default;
            } catch (error) {}

            message = element;
          }
        });
        yield put({
          type: 'updateState',
          payload: {
            msgInfo: message,
          },
        });
      }
    },
    // 验证状态（是否已报名、是否已环境初始化）
    *checkStatus({ payload }, { call, put }) {
      const { data } = yield call(
        CompetitionApiApiInstance.competitionCompetitionIDGet.bind(
          CompetitionApiApiInstance,
          payload,
        ),
      );
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            hasApplied: true,
          },
        });
      } else {
        yield put({
          type: 'updateState',
          payload: {
            hasApplied: false,
          },
        });
      }
    },
    // 获取已邀请人数
    *getInvitePersonNum({ payload }, { call, put }) {
      const { data } = yield call(
        CompetitionApiApiInstance.competitionCompetitionIDInvitationGet.bind(
          CompetitionApiApiInstance,
          payload,
        ),
      );
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            invitePersonNum: data || 0,
          },
        });
      } else {
        yield put({
          type: 'updateState',
          payload: {
            invitePersonNum: 0,
          },
        });
      }
    },
    // 报名
    *signUp({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          modalLoading: true,
        },
      });
      const { data } = yield call(
        CompetitionApiApiInstance.competitionCompetitionIDPost.bind(
          CompetitionApiApiInstance,
          payload.competitionID,
          payload.inviter,
          payload.formData || {},
        ),
      );
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            hasApplied: true,
            modalLoading: false,
          },
        });
      } else {
        yield put({
          type: 'updateState',
          payload: {
            modalLoading: false,
          },
        });
      }
    },
    // 初始化环境
    *initEnv({ payload }, { call, put }) {
      const initEnvJson = payload.initEnvJson;
      yield put({
        type: 'updateState',
        payload: {
          modalLoading: true,
        },
      });
      if (initEnvJson) {
        const { data } = yield call(
          EnvironmentsApiInstance.createEnvironment.bind(
            EnvironmentsApiInstance,
            initEnvJson.name,
            initEnvJson.config,
          ),
        );
        if (data) {
          message.success('初始化环境成功');
        }
        yield put({
          type: 'updateState',
          payload: {
            modalLoading: false,
          },
        });
        yield call(delay, 500);
        history.push(DEV_ENVIRONMENT);
      } else {
        message.error('无法初始化环境');
        yield put({
          type: 'updateState',
          payload: {
            modalLoading: false,
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

export default messageDetail;
