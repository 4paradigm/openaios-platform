import React, { useEffect } from 'react';
import JsYmal from 'js-yaml';
import { QuestionYaml } from './questionYaml';

/*
 * @Author: liyuying
 * @Date: 2021-05-20 17:22:22
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-27 14:20:09
 * @Description: file content
 */
const ApplicationPage = () => {
  const yamlToJson = () => {
    try {
      const doc: any = JsYmal.load(QuestionYaml);
      doc.openapi = 'zheshiyigezidingyide';
      console.log(doc);
      console.log(JsYmal.dump(doc));
    } catch (e) {
      console.log(e);
    }
  };
  useEffect(() => {
    yamlToJson();
  });
  return <div>应用</div>;
};

export default ApplicationPage;
