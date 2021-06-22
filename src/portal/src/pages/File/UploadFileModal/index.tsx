/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 11:43:29
 * @Description: file content
 */
import React, { useEffect } from 'react';
import { Modal, Upload, message } from 'cess-ui';
import { InboxOutlined } from '@ant-design/icons';
import { useSelector, useDispatch } from 'react-redux';
import { IFileListState, FileListActions } from '../models/file-list';
const { Dragger } = Upload;

const CreateFolderModal = () => {
  const dispatch = useDispatch();
  const { fileModalVisible, fileLoading, file, fileList, controller }: IFileListState = useSelector(
    (state: any) => state.fileList,
  );
  const handleOk = () => {
    dispatch({
      type: FileListActions.UPLOAD_FILE,
      payload: {
        file,
      },
    });
  };

  const handleCancle = () => {
    dispatch({
      type: FileListActions.CLOSE_MODAL,
    });
    if (fileLoading) {
      // 如果用户在文件上传过程中点击了取消，那么就取消刚才上传的文件的请求(通过umi-request的取消请求机制)
      controller.abort();
    }
  };
  useEffect(() => {
    const beforeunload = (event: any) => {
      if (fileModalVisible && fileLoading) {
        // 只有IE可以自定义文案
        event.returnValue = '文件会终止上传';
      }
    };
    window.addEventListener('beforeunload', beforeunload);
    return () => {
      window.removeEventListener('beforeunload', beforeunload);
    };
  });

  const handleRemoveFile = () => {
    dispatch({
      type: FileListActions.UPDATE_STATUS,
      payload: {
        file: null,
        fileList: [],
      },
    });
    return false;
  };

  const handleChangeFile = ({ file, fileList }: any) => {
    const lastFile = [...fileList].pop();
    const arr: any = new Array(1);
    arr[0] = lastFile;
    dispatch({
      type: FileListActions.UPDATE_STATUS,
      payload: {
        file,
        fileList: arr,
      },
    });
  };

  return (
    <Modal
      title="上传文件"
      visible={fileModalVisible}
      closable={false}
      confirmLoading={fileLoading}
      onOk={() => {
        handleOk();
      }}
      okButtonProps={{ disabled: !file }}
      okText="上传"
      onCancel={() => handleCancle()}
      destroyOnClose={true}
      centered
      className="create-folder-modal"
      width={480}
    >
      <Dragger
        fileList={fileList}
        onChange={(file) => handleChangeFile(file)}
        onRemove={handleRemoveFile}
        beforeUpload={() => {
          return false;
        }}
      >
        <p className="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p className="ant-upload-text">Click or drag file to this area to upload</p>
        {/* <p className="ant-upload-hint">
        Support for a single or bulk upload. Strictly prohibit from uploading company data or other
        band files
      </p> */}
      </Dragger>
    </Modal>
  );
};

export default CreateFolderModal;
