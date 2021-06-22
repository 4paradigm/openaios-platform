/*
 * @Author: liyuying
 * @Date: 2021-05-25 16:06:31
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-17 19:33:18
 * @Description: file content
 */
import React, { useState, useEffect } from 'react';
import { decode } from 'js-base64';
import { useSelector } from 'umi';
import ReactMarkdown from '@uiw/react-markdown-preview';
import { IApplicationInstanceCreateState } from '@/pages/Application/models/application-instance-create';

const ReadMe = () => {
  const [readMeDoc, setReadMeDoc] = useState('');
  const { applicationTemplate }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  useEffect(() => {
    if (applicationTemplate.files) {
      for (const name in applicationTemplate.files) {
        if (name === 'README.md') {
          setReadMeDoc(decode(applicationTemplate.files[name] || ''));
        }
      }
    }
  }, [applicationTemplate.files]);
  return <ReactMarkdown source={readMeDoc} linkTarget={'_blank'} />;
};
export default ReadMe;
