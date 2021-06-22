/*
 * @Author: liyuying
 * @Date: 2021-05-31 18:22:02
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-10 12:21:22
 * @Description: 以Yaml形式编辑内容
 */
import React, { useState, useEffect } from 'react';
import CodeMirror from '@uiw/react-codemirror';
import 'codemirror/keymap/sublime';
import 'codemirror/theme/darcula.css';
import {
  IApplicationInstanceCreateState,
  ApplicationInstanceCreateAction,
} from '@/pages/Application/models/application-instance-create';
import { useSelector, useDispatch } from 'umi';
import './index.less';

const EditAnswersYaml = () => {
  const dispatch = useDispatch();
  const { applicationInstance, applicationTemplate }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );

  /**
   * 编辑修改
   * @param value
   */
  const handelEditorChange = (editor: any) => {
    const value = editor.getValue();
    dispatch({
      type: ApplicationInstanceCreateAction.UPDATE_STATE,
      payload: {
        applicationInstance: { ...applicationInstance, answersYaml: value },
      },
    });
  };
  return (
    <div className="application-edit-answer-yaml">
      <CodeMirror
        value={applicationInstance.answersYaml}
        onChange={(value) => {
          handelEditorChange(value);
        }}
        options={{
          theme: 'darcula',
          keyMap: 'sublime',
          mode: 'YAML',
        }}
      />
    </div>
  );
};
export default EditAnswersYaml;
