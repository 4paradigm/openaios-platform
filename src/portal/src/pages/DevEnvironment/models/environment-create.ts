import { Model } from 'dva';
import { IEnvironmentData, IImage } from '@/interfaces/bussiness';
import { ImagesApiInstance, EnvironmentsApiInstance } from '@/openApi';
import { IAction } from '@/interfaces';
import juperterIcon from '@/assets/images/jupter.png';
import sshIcon from '@/assets/images/ssh.png';
import { history } from 'umi';

export const EnvironmentCreateAction = {
  UPDATE_STATE: 'environmentCreate/updateState',
  UPDATE_DATA: 'environmentCreate/updateData',
  GET_DATA: 'environmentCreate/getData',
  CREATE_ENVIRONMENT: 'environmentCreate/createEnvironment',
};

// 默认环境数据
const ENVIRONMENT_DATA: IEnvironmentData = {
  image: {
    repo: '',
    source: undefined,
    tags: [],
  },
  mounts: [],
  compute_unit: '',
  ssh: {
    enable: false,
    'id_rsa.pub': '',
  },
  jupyter: {
    enable: false,
    token: '',
  },
};

// 错误信息map
export const ErrorMap: any = {
  image: '',
  mounts: '',
  compute_unit: '',
  ssh: '',
  jupyter: '',
  // ssh和jupyter均未选择，
  interact: '',
  name: '',
};

export interface IEnvironmentCreateState {
  loading: boolean;
  name: string;
  publicMirror: IImage[];
  privateMirror: IImage[];
  environmentData: IEnvironmentData;
  errorMap: any;
  interactList: {
    key: string;
    name: string;
    desc: string;
    pswKey: string;
    icon: any;
    source: string;
  }[];
}

const defaultState: IEnvironmentCreateState = {
  loading: false,
  name: '',
  publicMirror: [],
  privateMirror: [],
  environmentData: JSON.parse(JSON.stringify(ENVIRONMENT_DATA)),
  errorMap: ErrorMap,
  interactList: [
    {
      name: 'JupyterLab',
      key: 'jupyter',
      desc: '交互编辑器，推荐使用',
      pswKey: 'token',
      icon: juperterIcon,
      source: 'public',
    },
    {
      name: 'SSH服务',
      key: 'ssh',
      desc: '允许远程机器SSH登录',
      pswKey: 'id_rsa.pub',
      icon: sshIcon,
      source: 'all',
    },
  ],
};

const environmentCreate: Model = {
  namespace: 'environmentCreate',
  state: defaultState,
  effects: {
    *getData({ payload }, { call, put }) {
      yield put({ type: 'initData' });
      yield put({ type: 'getComputeUnit' });
      yield put({ type: 'getPublicMirror' });
      yield put({ type: 'getPrivateMirror' });
    },
    *getPublicMirror({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.publicImagesGet.bind(ImagesApiInstance, 'env'));
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            publicMirror: data,
          },
        });
      }
    },
    *getPrivateMirror({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.imagesGet.bind(ImagesApiInstance));
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            privateMirror: data,
          },
        });
      }
    },
    *createEnvironment({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          loading: true,
        },
      });
      const { data } = yield call(
        EnvironmentsApiInstance.createEnvironment.bind(
          EnvironmentsApiInstance,
          payload.name,
          payload.data,
        ),
      );
      if (data) {
        history.push('/devEnvironment');
        yield put({
          type: 'updateState',
          payload: {
            loading: false,
            name: '',
            environmentData: JSON.parse(JSON.stringify(ENVIRONMENT_DATA)),
          },
        });
      } else {
        yield put({
          type: 'updateState',
          payload: {
            loading: false,
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
    updateData(state, { payload }: IAction) {
      return {
        ...state,
        environmentData: {
          ...state.environmentData,
          ...payload,
        },
      };
    },
    initData(state, { payload }: IAction) {
      return {
        ...state,
        environmentData: JSON.parse(JSON.stringify(ENVIRONMENT_DATA)),
        errorMap: ErrorMap,
      };
    },
  },
};

export default environmentCreate;
