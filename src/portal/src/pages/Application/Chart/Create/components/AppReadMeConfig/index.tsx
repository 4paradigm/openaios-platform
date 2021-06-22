/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:45:11
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-17 20:05:25
 * @Description: file content
 */
import React from 'react';
import { CloudUploadOutlined } from '@ant-design/icons';
import { message, Upload, notification, Alert } from 'cess-ui';
import {
  useDispatch,
  IApplicationChartCreateState,
  useSelector,
  ApplicationChartCreateAction,
} from 'umi';
import MarkdownEditor from '@uiw/react-markdown-editor';
import 'codemirror/keymap/sublime';
import 'codemirror/theme/darcula.css';
import './index.less';

/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:25:23
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-03 19:47:51
 * @Description: 创建我的应用--简介编辑
 */
const ApplicationChartAppReadMeConfig = () => {
  const dispatch = useDispatch();
  const { appReadMeDoc }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );
  const handleMdChange = (editor: any, data: any, value: string) => {
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { appReadMeDoc: value },
    });
  };
  /**
   * 上传readme文件
   */
  const handleMdlUpload = (files: any) => {
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
        payload: { appReadMeDoc: reader.result || '' },
      });
    };
    return true;
  };
  return (
    <div className="application-chart-readMe-config-page">
      <Alert type="info" message="请编辑或上传 Markdown 文件"></Alert>
      <Upload
        action=""
        accept=".md"
        beforeUpload={(files) => {
          return handleMdlUpload(files);
        }}
        showUploadList={false}
      >
        <CloudUploadOutlined />
      </Upload>
      <MarkdownEditor
        className="application-chart-readMe-config"
        value={appReadMeDoc}
        onChange={handleMdChange}
        height={240}
        options={{
          autofocus: true,
          showCursorWhenSelecting: true,
          theme: 'darcula',
          mode: 'Markdown',
        }}
      />
    </div>
  );
};
export default ApplicationChartAppReadMeConfig;
