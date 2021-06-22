import { Model } from 'dva';
import { message } from 'cess-ui';
import { StorageApiInstance } from '@/openApi';

export const FileListActions = {
  UPDATE_STATUS: 'fileList/updateStatus',
  CLOSE_MODAL: 'fileList/closeModal',
  CREATE_FOLDER: 'fileList/createFolder',
  GET_LIST: 'fileList/getList',
  UPLOAD_FILE: 'fileList/uploadFile',
  DELETE_FILE: 'fileList/deleteFile',
  DOWNLOAD_FILE: 'fileList/downloadFile',
  CHANGE_PATH: 'fileList/changePath',
};

export interface IFileListState {
  dataSource: any;
  filePath: string;
  folderModalVisible: boolean;
  folderLoading: boolean;
  fileModalVisible: boolean;
  fileLoading: boolean;
  file: any;
  fileList: any;
  downIndex: number;
  controller: any; // umi-request取消请求的控制器
}

const defaultState: IFileListState = {
  dataSource: [],
  filePath: '/',
  folderModalVisible: false,
  folderLoading: false,
  fileModalVisible: false,
  fileLoading: false,
  file: null,
  fileList: [],
  downIndex: -1,
  controller: null, // umi-request取消请求的控制器
};

const fileList: Model = {
  namespace: 'fileList',
  state: defaultState,
  effects: {
    *uploadFile({ payload }, { call, select, put }) {
      const { filePath } = yield select((state: { fileList: IFileListState }) => state.fileList);
      const controller = new AbortController(); // create a controller
      const { signal } = controller; // 使用AbortController.signal属性获取对其关联的AbortSignal对象的引用
      yield put({
        type: 'updateStatus',
        payload: {
          fileLoading: true,
          controller: controller,
        },
      });
      try {
        yield call(
          StorageApiInstance.storageUploadPost.bind(StorageApiInstance, filePath, payload.file, {
            timeout: 5000000,
            signal,
          }),
        );
        yield put({
          type: 'closeModal',
        });
        yield put({
          type: 'getList',
        });
      } catch {
        yield put({
          type: 'updateStatus',
          payload: {
            fileLoading: false,
          },
        });
      }
    },
    *getList({ payload }, { call, select, put }) {
      const { filePath } = yield select((state: { fileList: IFileListState }) => state.fileList);
      const { data } = yield call(
        StorageApiInstance.getDirectory.bind(StorageApiInstance, filePath),
      );

      // 没数据返回null
      yield put({
        type: 'updateStatus',
        payload: {
          dataSource: data || [],
        },
      });
    },
    *deleteFile({ payload }, { call, select, put }) {
      yield call(StorageApiInstance.deleteDirectoryOrFile.bind(StorageApiInstance, payload.path));
      message.success('删除成功');
      yield put({
        type: 'getList',
      });
    },
    *downloadFile({ payload }, { call, select, put }) {
      const { filePath, name, index } = payload;
      yield put({
        type: 'updateStatus',
        payload: {
          downIndex: index,
        },
      });
      const { data } = yield call(
        StorageApiInstance.storageDownloadGet.bind(StorageApiInstance, filePath + name),
      );
      if (data) {
        const blob = new Blob([data], { type: 'application/octet-stream' });
        const downloadElement: any = document.createElement('a');
        const href = window.URL.createObjectURL(blob); //创建下载的链接
        downloadElement.href = href;
        downloadElement.download = name; //下载后文件名
        document.body.appendChild(downloadElement);
        downloadElement.click(); //点击下载
        document.body.removeChild(downloadElement); //下载完成移除元素
        window.URL.revokeObjectURL(href); //释放掉blob对象
        yield put({
          type: 'updateStatus',
          payload: {
            downIndex: -1,
          },
        });
      } else {
        yield put({
          type: 'updateStatus',
          payload: {
            downIndex: -1,
          },
        });
      }
    },
    *createFolder({ payload }, { call, put, select }) {
      yield put({
        type: 'updateStatus',
        payload: {
          folderLoading: true,
        },
      });
      const { filePath } = yield select((state: { fileList: IFileListState }) => state.fileList);
      const { data } = yield call(
        StorageApiInstance.createDirectory.bind(StorageApiInstance, `${filePath}${payload.path}`),
      );
      if (data) {
        yield put({
          type: 'closeModal',
        });
        yield put({
          type: 'getList',
        });
      } else {
        yield put({
          type: 'updateStatus',
          payload: {
            folderLoading: false,
          },
        });
      }
    },
    *changePath({ payload }, { call, put, select }) {
      yield put({
        type: 'updateStatus',
        payload: {
          filePath: payload.path,
        },
      });
      yield put({
        type: 'getList',
      });
    },
  },
  reducers: {
    updateStatus(state, { payload }: any) {
      return {
        ...state,
        ...payload,
      };
    },
    closeModal(state, { payload }: any) {
      return {
        ...state,
        folderModalVisible: false,
        folderLoading: false,
        fileModalVisible: false,
        fileLoading: false,
        file: null,
        fileList: [],
      };
    },
  },
};

export default fileList;
