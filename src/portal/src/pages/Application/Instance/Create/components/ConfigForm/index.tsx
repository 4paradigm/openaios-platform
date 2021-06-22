import React, { useImperativeHandle } from 'react';
import JsYmal from 'js-yaml';
import { Form, Input, Button, Upload, Modal, notification } from 'cess-ui';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import {
  useDispatch,
  IApplicationInstanceCreateState,
  useSelector,
  ApplicationInstanceCreateAction,
} from 'umi';
import EditAnswersYaml from '../EditAnswersYaml';
import ConfigGroup from './ConfigGroup';
import {
  APPLICATION_BASIC_CONFIG_VAR_MAP,
  EDIT_APPLICATION_INSTANCE_ANSWER_MODE,
  APPLICATION_BASIC_CONFIG_VAR,
} from '@/constant/application';
import { transferObjectFromFlatToLevles, transferObjectFromLevlesToFlat } from '@/utils';
import './index.less';

/*
 * @Author: liyuying
 * @Date: 2021-06-07 14:07:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-07 14:12:31
 * @Description: file content
 */
const layout = {
  labelCol: { flex: '240px' },
  wrapperCol: { flex: 'auto' },
};
const AppInstanceConfigForm = ({}, ref: any) => {
  const dispatch = useDispatch();
  const {
    applicationInstance,
    applicationTemplate,
    editMode,
  }: IApplicationInstanceCreateState = useSelector((state: any) => state.applicationInstanceCreate);
  const [form] = Form.useForm();
  /**
   * 配置表单的更改
   */
  const handleFieldsChange = (changedFields: any, allFields: any) => {
    try {
      // 修改的是实例名称
      if (
        changedFields &&
        changedFields[0].name &&
        changedFields[0].name[0] === APPLICATION_BASIC_CONFIG_VAR.NAME
      ) {
        dispatch({
          type: ApplicationInstanceCreateAction.UPDATE_STATE,
          payload: {
            applicationInstance: { ...applicationInstance, name: changedFields[0].value },
          },
        });
      } else {
        const answers: any = {};
        allFields.forEach((element: any) => {
          if (element.name && element.name[0]) {
            if (APPLICATION_BASIC_CONFIG_VAR_MAP[element.name[0]]) {
            } else {
              answers[element.name[0]] = element.value;
            }
          }
        });
        // 配置answer.yaml
        // 拆分对象的key为多个层级
        const levelsAnswers = transferObjectFromFlatToLevles(answers);
        const answersYaml = JsYmal.dump(levelsAnswers);
        applicationTemplate.files = {
          ...applicationTemplate.files,
          'answers.yaml': answersYaml,
        };
        dispatch({
          type: ApplicationInstanceCreateAction.UPDATE_STATE,
          payload: {
            applicationInstance: { ...applicationInstance, answers, answersYaml },
            applicationTemplate,
          },
        });
      }
    } catch (error) {
      console.warn(error);
    }
  };
  /**
   * 切换配置模式
   */
  const handleEditModeSwitch = () => {
    let yamlDoc = '';
    // 切换为Yaml配置
    if (editMode === EDIT_APPLICATION_INSTANCE_ANSWER_MODE.FORM) {
      try {
        // 配置answer.yaml
        // 拆分对象的key为多个层级
        const answers = transferObjectFromFlatToLevles(applicationInstance.answers);
        yamlDoc = JsYmal.dump(answers);
      } catch (error) {
        console.warn(error);
      }
      dispatch({
        type: ApplicationInstanceCreateAction.UPDATE_STATE,
        payload: {
          applicationInstance: { ...applicationInstance, answersYaml: yamlDoc },
          editMode: EDIT_APPLICATION_INSTANCE_ANSWER_MODE.YAML,
        },
      });
    } else {
      // 切换为form配置
      Modal.confirm({
        title: '切换为编辑表单',
        icon: <ExclamationCircleOutlined />,
        content: '由YAML切换为表单，会丢失在YAML中编辑的内容，确认切换？',
        okText: '确认',
        cancelText: '取消',
        onOk: () => {
          dispatch({
            type: ApplicationInstanceCreateAction.UPDATE_STATE,
            payload: {
              editMode: EDIT_APPLICATION_INSTANCE_ANSWER_MODE.FORM,
            },
          });
        },
      });
    }
  };
  /**
   * 上传answer.yaml文件
   */
  const handleAnswerYamlUpload = (files: any) => {
    if (typeof FileReader === 'undefined') {
      notification.error({ message: '您的浏览器不支持FileReader接口' });
      return false;
    }
    let reader = new FileReader();
    reader.readAsText(files, 'utf-8');
    reader.onload = () => {
      applicationTemplate.files = {
        ...applicationTemplate.files,
        'answers.yaml': reader.result + '' || '',
      };
      dispatch({
        type: ApplicationInstanceCreateAction.UPDATE_STATE,
        payload: {
          applicationInstance: { ...applicationInstance, answersYaml: reader.result },
          applicationTemplate,
        },
      });
    };
    return true;
  };
  /**
   * 测试answer.yaml 文件是否符合规则
   */
  const handleTestAnswerYaml = () => {
    if (applicationInstance.answersYaml) {
      try {
        const flatObj: any = JsYmal.load(applicationInstance.answersYaml);
        // 合并对象的key,成一个层级
        const answers = transferObjectFromLevlesToFlat(flatObj);
        console.log(answers);
        notification.success({ message: 'YAML文件配置无误' });
      } catch (error) {
        notification.error({ message: 'YAML文件格式错误' });
      }
    } else {
      notification.error({ message: 'YAML文件为空' });
    }
  };
  /**
   * 暴露方法给父组件
   */
  useImperativeHandle(ref, () => ({
    createApplication: async () => {
      // 已绑定
      if (form && form.__INTERNAL__ && form.__INTERNAL__.name) {
        try {
          const values = await form.validateFields();
          console.log('Success:', values);
          const answers = JsYmal.load(applicationInstance.answersYaml || '');
          dispatch({
            type: ApplicationInstanceCreateAction.CREATE_INSTANCE,
            payload: {
              appInstanceName: values.name,
              answers,
              url: applicationTemplate.metadata?.url,
            },
          });
        } catch (errorInfo) {
          const message = errorInfo.errorFields[0].errors[0];
          notification.error({ message });
          console.log('Failed:', errorInfo);
        }
      } else {
        notification.error({ message: '请先进行应用配置' });
      }
    },
  }));
  return (
    <Form
      className="app-create-form"
      labelAlign="left"
      {...layout}
      scrollToFirstError={true}
      name="basic"
      form={form}
      colon={false}
      onFieldsChange={handleFieldsChange}
    >
      <Form.Item
        label="实例名称"
        name={APPLICATION_BASIC_CONFIG_VAR.NAME}
        rules={[
          { required: true, message: '请输入应用实例名称' },
          () => ({
            validator(_, value) {
              if (!value) {
                return Promise.resolve();
              }
              if (!value.trim()) {
                return Promise.reject(new Error('应用实例名称不可以只包含空格'));
              } else {
                const fitNameList = new RegExp(/[a-z]([-a-z0-9]*[a-z0-9])?/).exec(value) || [];
                if (fitNameList.length > 0 && fitNameList[0] === value) {
                  return Promise.resolve();
                } else {
                  return Promise.reject(
                    new Error(
                      '应用实例名称以小写字母开头，只能包含小写英文字母、数字和中划线，且不能以中划线结尾',
                    ),
                  );
                }
              }
            },
          }),
        ]}
        className="app-config-item"
      >
        <Input placeholder="请输入应用实例名称" />
      </Form.Item>
      <div className="app-create-switch-bar">
        {applicationInstance &&
        applicationInstance.questionsGroupList &&
        applicationInstance.questionsGroupList.length > 0 ? (
          <Button
            type="primary"
            onClick={() => {
              handleEditModeSwitch();
            }}
          >
            {editMode === EDIT_APPLICATION_INSTANCE_ANSWER_MODE.FORM ? '编辑YAML' : '编辑表单'}
          </Button>
        ) : (
          <></>
        )}
        {editMode === EDIT_APPLICATION_INSTANCE_ANSWER_MODE.YAML ? (
          <>
            <Upload
              className="ant-btn-primary"
              action=""
              accept=".yaml"
              beforeUpload={(files) => {
                return handleAnswerYamlUpload(files);
              }}
              showUploadList={false}
            >
              <Button type="primary">从文件读取</Button>
            </Upload>

            <Button type="primary" onClick={handleTestAnswerYaml}>
              测试YAML配置
            </Button>
          </>
        ) : (
          <></>
        )}
      </div>
      {editMode === EDIT_APPLICATION_INSTANCE_ANSWER_MODE.FORM ? (
        applicationInstance &&
        applicationInstance.questionsGroupList &&
        applicationInstance.questionsGroupList.length > 0 ? (
          applicationInstance.questionsGroupList.map((item) => {
            return <ConfigGroup questionsGroup={item} key={item.groupName}></ConfigGroup>;
          })
        ) : (
          <></>
        )
      ) : (
        <EditAnswersYaml></EditAnswersYaml>
      )}
    </Form>
  );
};
export default React.forwardRef(AppInstanceConfigForm);
