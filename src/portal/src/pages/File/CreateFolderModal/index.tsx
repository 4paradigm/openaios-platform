/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-01 18:31:47
 * @Description: file content
 */
import React, { useRef } from 'react';
import { Modal, Form, Input } from 'cess-ui';
import { useSelector, useDispatch } from 'react-redux';
import { IFileListState, FileListActions } from '../models/file-list';

const CreateFolderModal = () => {
  const formRef = useRef(null);
  const dispatch = useDispatch();
  const { folderModalVisible, folderLoading }: IFileListState = useSelector(
    (state: any) => state.fileList,
  );

  const handleOk = (data: any) => {
    dispatch({
      type: FileListActions.CREATE_FOLDER,
      payload: {
        path: data.path,
      },
    });
  };

  const handleCancle = () => {
    dispatch({
      type: FileListActions.CLOSE_MODAL,
    });
  };

  return (
    <Modal
      title="新建文件夹"
      visible={folderModalVisible}
      closable={false}
      confirmLoading={folderLoading}
      onOk={() => {
        (formRef.current as any)
          .validateFields()
          .then((values: any) => {
            handleOk(values);
          })
          .catch(() => {});
      }}
      okText="保存更改"
      onCancel={() => handleCancle()}
      destroyOnClose={true}
      centered
      className="create-folder-modal"
      width={480}
    >
      <Form ref={formRef} requiredMark={false} layout="vertical" preserve={false}>
        <Form.Item
          label="文件夹名称"
          rules={[{ required: true, message: '请输入文件夹名称' }]}
          name="path"
        >
          <Input maxLength={32} placeholder="请输入文件夹名称" autoComplete="off" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CreateFolderModal;
