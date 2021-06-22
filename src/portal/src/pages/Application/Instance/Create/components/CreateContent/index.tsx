/*
 * @Author: liyuying
 * @Date: 2021-06-07 16:11:16
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-20 13:43:33
 * @Description: file content
 */
import React, { useRef, useImperativeHandle } from 'react';
import { Collapse, notification } from 'cess-ui';
import PreviewYaml from '../PreviewYaml';
import VersionSelect from '../VersionSelect';
import AppReadMe from '../AppReadMe';
import ReadMe from '../ReadMe';
import AppInstanceConfigForm from '../ConfigForm';

import './index.less';
import { APPLICATION_INSTANCE_CREATE_COLLAPSE } from '@/constant/application';
import { IApplicationInstanceCreateState, useDispatch, useSelector } from 'umi';

const { Panel } = Collapse;
/*
 * @Author: liyuying
 * @Date: 2021-06-07 16:11:16
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-07 16:13:43
 * @Description: 创建实例的内容部分（去除顶部导航和底部按钮）
 */
const AppInstanceCreateContent = ({}, ref: any) => {
  const { applicationTemplate }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  const formRef = useRef();
  /**
   * 暴露方法给父组件
   */
  useImperativeHandle(ref, () => ({
    createApplication: () => {
      // 已绑定
      if (formRef.current) {
        (formRef.current as any).createApplication();
      } else {
        notification.error({ message: '请先进行应用配置' });
      }
    },
  }));
  /**
   * 切换展示的折叠面板
   * @param value
   */
  const handleCollapseChange = (value: any) => {};
  return (
    <div className="application-create-contanner">
      {/* 版本选择 */}
      <div className="application-create-version-bar">
        <VersionSelect></VersionSelect>
      </div>
      {/* 顶部展示md */}
      <AppReadMe></AppReadMe>
      <Collapse
        defaultActiveKey={[APPLICATION_INSTANCE_CREATE_COLLAPSE.CONFIOG_OPTIONS]}
        onChange={handleCollapseChange}
      >
        {/* 详情md */}
        <Panel
          header={APPLICATION_INSTANCE_CREATE_COLLAPSE.DETAIL_DISCRIBTION}
          key={APPLICATION_INSTANCE_CREATE_COLLAPSE.DETAIL_DISCRIBTION}
        >
          <ReadMe></ReadMe>
        </Panel>
        {/* 配置信息 */}
        <Panel
          header={APPLICATION_INSTANCE_CREATE_COLLAPSE.CONFIOG_OPTIONS}
          key={APPLICATION_INSTANCE_CREATE_COLLAPSE.CONFIOG_OPTIONS}
        >
          <AppInstanceConfigForm ref={formRef}></AppInstanceConfigForm>
        </Panel>
        {/* 预览 */}
        <Panel
          header={APPLICATION_INSTANCE_CREATE_COLLAPSE.PREVIEW}
          key={APPLICATION_INSTANCE_CREATE_COLLAPSE.PREVIEW}
        >
          <PreviewYaml></PreviewYaml>
        </Panel>
      </Collapse>
    </div>
  );
};
export default React.forwardRef(AppInstanceCreateContent);
