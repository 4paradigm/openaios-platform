/*
 * @Author: liyuying
 * @Date: 2021-05-25 16:06:31
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-17 19:32:52
 * @Description: file content
 */
import React, { useState, useEffect } from 'react';
import { decode } from 'js-base64';
import { useSelector } from 'umi';
import ReactMarkdown from '@uiw/react-markdown-preview';
import { Avatar, Icon } from 'cess-ui';
import { IApplicationInstanceCreateState } from '@/pages/Application/models/application-instance-create';

import './index.less';

const AppReadMe = () => {
  const [readMeDoc, setReadMeDoc] = useState('');
  const { applicationTemplate }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  useEffect(() => {
    setReadMeDoc('');
    if (applicationTemplate.files) {
      for (const name in applicationTemplate.files) {
        if (name === 'app-readme.md') {
          setReadMeDoc(decode(applicationTemplate.files[name] || ''));
        }
      }
    }
  }, [applicationTemplate.files]);
  return (
    <div className="app-read-me-page">
      <div className="icon">
        <Avatar icon={<Icon type="application" />} src={applicationTemplate.metadata?.icon_link} />
      </div>
      <div className="doc">
        {readMeDoc ? (
          <ReactMarkdown source={readMeDoc} linkTarget={'_blank'} skipHtml={false} />
        ) : (
          applicationTemplate.metadata?.description
        )}
      </div>
    </div>
  );
};
export default AppReadMe;
