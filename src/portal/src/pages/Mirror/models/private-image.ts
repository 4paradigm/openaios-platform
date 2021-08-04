import { Model } from 'dva';
import { IAction } from '@/interfaces';
import { ITask } from '@/interfaces/bussiness';
import { ImagesApiInstance } from '@/openApi';
import { message } from 'cess-ui';

export const PrivateImageAction = {
  GET_LIST: 'privateImage/getList',
  DELETE_IMAGE: 'privateImage/deleteImage',
  IMPORT_IMAGE: 'privateImage/importImage',
  OPEN_IMAGE_MODAL: 'privateImage/openImageModal',
  GET_TASK_LIST: 'privateImage/getTaskList',
  DELETE_TASK: 'privateImage/deleteTask',
  CLOSE_MODAL: 'privateImage/closeModal',
  STOP_TASK: 'privateImage/stopTask',
  GET_TOTAL: 'privateImage/getTotal',
  UPDATE_STATE: 'privateImage/updateState',
  INITT_DATA: 'privateImage/initData',
  COPY_IMAGE: 'privateImage/copyImage',
};

export interface IMirrorPrivateState {
  dataSource: any;
  total: number;
  currentPage: number;
  modalVisible: boolean;
  modalLoading: boolean;
  copyImageModalVisible: boolean;
  copyImageLoading: boolean;
  copyImageSourceImage: any;
  registryList: { url: string; id: number }[];
  taskModalVisible: boolean;
  taskList: ITask[];
}

const defaultState: IMirrorPrivateState = {
  dataSource: [],
  total: 0,
  currentPage: 1,
  modalVisible: false,
  modalLoading: false,
  copyImageModalVisible: false,
  copyImageLoading: false,
  copyImageSourceImage: null,
  registryList: [],
  taskModalVisible: false,
  taskList: [],
};

const privateImage: Model = {
  namespace: 'privateImage',
  state: defaultState,
  effects: {
    // 初始化所有数据
    *initData({ payload }, { call, put, select }) {
      let { currentPage } = yield select((state: any) => state.privateImage);
      yield put({
        type: 'getList',
        payload: currentPage,
      });
      yield put({ type: 'getTotal' });
      yield put({ type: 'getTaskList' });
    },
    *getList({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.imagesGet.bind(ImagesApiInstance, payload, 10));
      let currentPage = payload;
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            dataSource: data,
            currentPage: currentPage,
          },
        });
      }
    },
    *getTotal({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.imagesInfoGet.bind(ImagesApiInstance));
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            total: data.image_count,
          },
        });
      }
    },
    *getTaskList({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.listImportingImages.bind(ImagesApiInstance));
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            taskList: data,
          },
        });
      }
    },
    *deleteImage({ payload }, { call, put, select }) {
      const { repo, digest } = payload;
      const { data } = yield call(
        ImagesApiInstance.imagesDelete.bind(ImagesApiInstance, repo, digest),
      );
      if (data) {
        let { currentPage, total } = yield select((state: any) => state.privateImage);
        if (total % 10 === 1) {
          // 如果删除的是当页的第一个也是唯一一个环境时，删除之后要把页码减1
          let nextPage = currentPage - 1;
          currentPage = nextPage ? nextPage : 1;
        }
        yield put({
          type: 'updateState',
          payload: {
            currentPage: currentPage,
          },
        });
        yield put({
          type: 'getList',
          payload: currentPage,
        });
        yield put({ type: 'getTotal' });
      }
    },
    *openImageModal({ payload }, { call, put }) {
      const { data } = yield call(ImagesApiInstance.imagesRegistryGet.bind(ImagesApiInstance));
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            registryList: data,
            modalVisible: true,
          },
        });
      }
    },
    *importImage({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          modalLoading: true,
        },
      });
      const { data } = yield call(
        ImagesApiInstance.imagesImportingPost.bind(ImagesApiInstance),
        payload.registryId,
        payload.repo,
        payload.tag,
      );
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            modalLoading: false,
            modalVisible: false,
          },
        });
        yield put({ type: 'initData' });
      } else {
        yield put({
          type: 'updateState',
          payload: {
            modalLoading: false,
          },
        });
      }
    },
    *deleteTask({ payload }, { call, put }) {
      const { data } = yield call(
        ImagesApiInstance.imagesImportingDelete.bind(ImagesApiInstance),
        payload,
      );
      if (data) {
        yield put({
          type: 'getTaskList',
        });
      }
    },
    *stopTask({ payload }, { call, put }) {
      const { data } = yield call(
        ImagesApiInstance.imagesImportingPut.bind(ImagesApiInstance),
        payload,
      );
      if (data) {
        yield put({
          type: 'getTaskList',
        });
      }
    },
    *copyImage({ payload }, { call, put }) {
      yield put({
        type: 'updateState',
        payload: {
          copyImageLoading: true,
        },
      });
      const { data } = yield call(
        ImagesApiInstance.imagesPut.bind(ImagesApiInstance),
        payload.srcRepo,
        payload.destRepo,
        payload.tag,
      );
      if (data) {
        yield put({
          type: 'updateState',
          payload: {
            copyImageLoading: false,
            copyImageModalVisible: false,
          },
        });
        message.success('拷贝镜像成功');
        yield put({ type: 'initData' });
      } else {
        yield put({
          type: 'updateState',
          payload: {
            copyImageLoading: false,
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
    closeModal(state, { payload }: IAction) {
      return {
        ...state,
        modalVisible: false,
        modalLoading: false,
        taskModalVisible: false,
      };
    },
  },
};

export default privateImage;
