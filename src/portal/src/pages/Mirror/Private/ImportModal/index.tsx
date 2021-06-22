/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-04-25 17:57:26
 * @Description: file content
 */
import React, { useRef } from 'react';
import { Modal, Form, Input, Select } from 'cess-ui';
import { useDispatch, useSelector } from 'umi';
import { PrivateImageAction } from '../../models/private-image';
const { Option } = Select;

function ImportModal() {
  const { modalVisible, modalLoading, registryList } = useSelector(
    (state: any) => state.privateImage,
  );
  const dispatch = useDispatch();
  const formRef = useRef(null);

  const handleOk = (data: any) => {
    dispatch({
      type: PrivateImageAction.IMPORT_IMAGE,
      payload: data,
    });
  };

  const handleCancle = () => {
    dispatch({ type: PrivateImageAction.CLOSE_MODAL });
  };

  return (
    <Modal
      title="导入镜像"
      visible={modalVisible}
      closable={false}
      confirmLoading={modalLoading}
      onOk={() => {
        (formRef.current as any)
          .validateFields()
          .then((values: any) => {
            handleOk(values);
          })
          .catch(() => {});
      }}
      okText="导入"
      onCancel={() => handleCancle()}
      destroyOnClose={true}
      centered
      width={480}
    >
      <Form ref={formRef} requiredMark={false} layout="vertical" preserve={false}>
        <Form.Item
          label="registry"
          rules={[{ required: true, message: '请选择registry' }]}
          name="registryId"
        >
          <Select style={{ width: '433px' }} placeholder="请选择registry">
            {registryList.map((registry: any) => (
              <Option key={registry.id} value={registry.id}>
                {registry.url}
              </Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="repo"
          rules={[
            { required: true, message: '请输入repo' },
            () => ({
              validator(_, value) {
                if (!value || new RegExp('/').test(value)) {
                  return Promise.resolve();
                } else {
                  return Promise.reject(
                    '必须包含"/"字符,如果是官方镜像用户名需要加上library，ex: library/ubuntu',
                  );
                }
              },
            }),
          ]}
          name="repo"
        >
          <Input
            maxLength={100}
            placeholder='请输入repo，必须包含"/"字符，官方镜像用户名需要加上library'
            autoComplete="off"
          />
        </Form.Item>
        <Form.Item label="tag" rules={[{ required: true, message: '请输入tag' }]} name="tag">
          <Input maxLength={100} placeholder="请输入tag" autoComplete="off" />
        </Form.Item>
      </Form>
    </Modal>
  );
}

export default ImportModal;
