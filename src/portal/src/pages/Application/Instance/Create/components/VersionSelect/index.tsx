import React, { useEffect } from 'react';
import {
  IApplicationInstanceCreateState,
  ApplicationInstanceCreateAction,
} from '@/pages/Application/models/application-instance-create';
import { useSelector, useDispatch } from 'umi';
import { Select, Form } from 'cess-ui';

/*
 * @Author: liyuying
 * @Date: 2021-05-28 20:20:34
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-17 19:57:06
 * @Description: file content
 */
const VersionSelect = () => {
  const dispatch = useDispatch();
  const { applicationTemplate, applicationInstance }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  const [form] = Form.useForm();

  const layout = {
    labelCol: { flex: '78px' },
    wrapperCol: { flex: 'auto' },
  };
  /**
   * 修改版本
   * @param version
   */
  const changeVersion = (version: any) => {
    dispatch({
      type: ApplicationInstanceCreateAction.GET_TEMPLATE,
      payload: {
        category: applicationTemplate.metadata?.category,
        name: applicationTemplate.metadata?.name,
        version,
      },
    });
  };
  useEffect(() => {
    form.setFieldsValue({
      version: applicationInstance.version,
    });
  }, [applicationInstance.version]);
  return (
    <Form
      labelAlign="left"
      {...layout}
      scrollToFirstError={true}
      name="version"
      form={form}
      colon={false}
    >
      <Form.Item label="应用版本" name="version" className="app-config-item">
        {applicationTemplate.version_list && applicationTemplate.version_list.length > 0 ? (
          <Select
            value={applicationInstance.version}
            onChange={(value) => {
              changeVersion(value);
            }}
            style={{ width: '400px' }}
          >
            {(applicationTemplate.version_list || []).map((item) => {
              return (
                <Select.Option key={item} value={item}>
                  {item}
                </Select.Option>
              );
            })}
          </Select>
        ) : (
          applicationInstance.version
        )}
      </Form.Item>
    </Form>
  );
};
export default VersionSelect;
