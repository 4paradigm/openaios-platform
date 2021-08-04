/*
 * @Author: liyuying
 * @Date: 2021-06-22 18:06:24
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-24 16:27:44
 * @Description: file content
 */
import React, { useRef } from 'react';
import { Modal, Form, Input } from 'cess-ui';
import { useSelector, useDispatch } from 'react-redux';
import { IMirrorPrivateState, PrivateImageAction } from 'umi';
import './index.less';

const CopyImageModal = () => {
  const formRef = useRef(null);
  const dispatch = useDispatch();
  const {
    copyImageModalVisible,
    copyImageLoading,
    copyImageSourceImage,
  }: IMirrorPrivateState = useSelector((state: any) => state.privateImage);

  const handleOk = (data: any) => {
    dispatch({
      type: PrivateImageAction.COPY_IMAGE,
      payload: {
        srcRepo: copyImageSourceImage.repo,
        tag: copyImageSourceImage.tags && [copyImageSourceImage.tags[0]],
        destRepo: data.name,
      },
    });
  };

  const handleCancle = () => {
    if (copyImageLoading) {
      Modal.confirm({
        title: '镜像拷贝中...',
        content: `确认关闭？拷贝成功后自动刷新列表`,
        okText: '确认',
        cancelText: '取消',
        onOk: () => {
          dispatch({
            type: PrivateImageAction.UPDATE_STATE,
            payload: {
              copyImageModalVisible: false,
            },
          });
        },
      });
    } else {
      dispatch({
        type: PrivateImageAction.UPDATE_STATE,
        payload: {
          copyImageModalVisible: false,
        },
      });
    }
  };

  return (
    <Modal
      title="拷贝镜像"
      visible={copyImageModalVisible}
      closable={false}
      confirmLoading={copyImageLoading}
      onOk={() => {
        (formRef.current as any)
          .validateFields()
          .then((values: any) => {
            handleOk(values);
          })
          .catch(() => {});
      }}
      okText="拷贝"
      onCancel={() => handleCancle()}
      destroyOnClose={true}
      centered
      className="copy-image-modal"
      width={480}
    >
      <div className="copy-image-item">
        <label>源仓库名称</label>
        <span> {(copyImageSourceImage && copyImageSourceImage.repo) || ''}</span>
      </div>
      <div className="copy-image-item">
        <label>Tags</label>
        <span> {(copyImageSourceImage && copyImageSourceImage.tags.join(',')) || ''}</span>
      </div>
      <Form ref={formRef} requiredMark={false} layout="vertical" preserve={false}>
        <Form.Item
          label="目标仓库名称"
          className="target"
          rules={[
            { required: true, message: '请输入仓库名称' },
            () => ({
              validator(_, value) {
                if (!value) {
                  return Promise.resolve();
                }
                if (!value.trim()) {
                  return Promise.reject(new Error('仓库名称不可以只包含空格'));
                } else {
                  const fitNameList = new RegExp(/[_./a-z0-9]*/).exec(value) || [];
                  if (fitNameList.length > 0 && fitNameList[0] === value) {
                    return Promise.resolve();
                  } else {
                    return Promise.reject(
                      new Error(
                        "仓库名被分解为路径组件。仓库名必须至少有一个小写字母、字母数字字符，可选句点、破折号或下划线分隔。如果仓库名有两个或多个路径组件，则它们必须用正斜杠('/')分隔。",
                      ),
                    );
                  }
                }
              },
            }),
          ]}
          name="name"
        >
          <Input maxLength={64} placeholder="请输入目标仓库名称" autoComplete="off" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default CopyImageModal;
