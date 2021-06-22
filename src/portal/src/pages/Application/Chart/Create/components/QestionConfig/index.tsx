/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:25:40
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-07 15:31:31
 * @Description: file content
 */
import React from 'react';
import JsYmal from 'js-yaml';
import {
  useDispatch,
  IApplicationChartCreateState,
  useSelector,
  ApplicationChartCreateAction,
} from 'umi';
import CodeMirror from '@uiw/react-codemirror';
import 'codemirror/keymap/sublime';
import 'codemirror/theme/darcula.css';
import './index.less';
import { Upload, notification, Button } from 'cess-ui';
import AppChartPerviewFormModal from './PerviewFormModal';

/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:25:23
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-04 19:36:11
 * @Description: 创建我的应用--question.yaml配置
 */
const ApplicationChartQuestionConfig = () => {
  const dispatch = useDispatch();
  const { questionYaml }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );

  /**
   * 编辑修改
   * @param value
   */
  const handelEditorChange = (editor: any) => {
    const value = editor.getValue();
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { questionYaml: value },
    });
  };
  /**
   * 上传yaml文件
   */
  const handleQuestionYamlUpload = (files: any) => {
    if (typeof FileReader === 'undefined') {
      notification.error({ message: '您的浏览器不支持FileReader接口' });
      return false;
    }
    let reader = new FileReader();
    reader.readAsText(files, 'utf-8');
    reader.onload = () => {
      console.log(reader.result);
      dispatch({
        type: ApplicationChartCreateAction.UPDATE_STATE,
        payload: { questionYaml: reader.result || '' },
      });
    };
    return true;
  };
  /**
   * 测试answer.yaml 文件是否符合规则
   */
  const handleTestQuestionYaml = () => {
    if (questionYaml) {
      try {
        const flatObj: any = JsYmal.load(questionYaml);
        console.log(flatObj);
        // // 合并对象的key,成一个层级
        // const answers = transferObjectFromLevlesToFlat(flatObj);
        // console.log(answers);
        notification.success({ message: 'YAML文件配置无误' });
      } catch (error) {
        notification.error({ message: 'YAML文件格式错误' });
      }
    } else {
      notification.error({ message: 'YAML文件为空' });
    }
  };
  /**
   * 预览表单
   */
  const handlePreviewHome = () => {
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { modalVisible: true },
    });
  };
  return (
    <div className="application-chart-question-config">
      <div className="question-config-button-bar">
        <Button type="primary" onClick={handleTestQuestionYaml}>
          测试YAML配置
        </Button>
        <Upload
          action=""
          accept=".yaml,.yml"
          beforeUpload={(files) => {
            return handleQuestionYamlUpload(files);
          }}
          showUploadList={false}
        >
          <Button type="primary">从文件读取</Button>
        </Upload>
        <Button type="primary" onClick={handlePreviewHome}>
          表单预览
        </Button>
      </div>

      <CodeMirror
        value={questionYaml}
        onChange={(value) => {
          handelEditorChange(value);
        }}
        options={{
          theme: 'darcula',
          keyMap: 'sublime',
          mode: 'YAML',
        }}
      />
      <AppChartPerviewFormModal></AppChartPerviewFormModal>
    </div>
  );
};
export default ApplicationChartQuestionConfig;
