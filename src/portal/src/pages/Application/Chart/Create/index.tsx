/*
 * @Author: liyuying
 * @Date: 2021-06-03 14:36:30
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-03 19:13:26
 * @Description: file content
 */
import React, { useEffect } from 'react';
import './index.less';
import { Breadcrumb, Steps, Card } from 'cess-ui';
import { APPLICATION_CHART } from '@/router/url';
import { APPLICATION_CHART_CREATE_STEPS } from '@/constant/application';
import {
  useDispatch,
  IApplicationChartCreateState,
  useSelector,
  ApplicationChartCreateAction,
} from 'umi';
import ApplicationChartUpload from './components/UploadChart';
import ApplicationChartAppReadMeConfig from './components/AppReadMeConfig';
import ApplicationChartQuestionConfig from './components/QestionConfig';
import ApplicationChartAppPreview from './components/AppPreview';
import ApplicationChartCreateButtonBar from './components/CreateButtonBar';
/*
 * @Author: liyuying
 * @Date: 2021-06-03 14:36:30
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-03 17:59:43
 * @Description: 应用模版--创建
 */
const breadcrumb = (
  <Breadcrumb>
    <Breadcrumb.Item href={APPLICATION_CHART}>应用市场</Breadcrumb.Item>
    <Breadcrumb.Item>创建我的应用</Breadcrumb.Item>
  </Breadcrumb>
);
const ApplicationChartCreate = () => {
  const dispatch = useDispatch();
  const { loading, currentStep }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );
  useEffect(() => {
    dispatch({
      type: ApplicationChartCreateAction.INIT_DATA,
    });
  }, []);
  return (
    <div className="application-chart-create comm-create-page">
      {breadcrumb}
      <Card>
        <Steps current={currentStep}>
          <Steps.Step title={APPLICATION_CHART_CREATE_STEPS.UPLOAD.title} />
          <Steps.Step title={APPLICATION_CHART_CREATE_STEPS.APP_README.title} />
          <Steps.Step title={APPLICATION_CHART_CREATE_STEPS.QUESTION_CONFIG.title} />
          <Steps.Step title={APPLICATION_CHART_CREATE_STEPS.PREVIEW.title} />
        </Steps>
        <div className="application-chart-create-content">
          {currentStep === APPLICATION_CHART_CREATE_STEPS.UPLOAD.step ? (
            <ApplicationChartUpload />
          ) : currentStep === APPLICATION_CHART_CREATE_STEPS.APP_README.step ? (
            <ApplicationChartAppReadMeConfig />
          ) : currentStep === APPLICATION_CHART_CREATE_STEPS.QUESTION_CONFIG.step ? (
            <ApplicationChartQuestionConfig />
          ) : currentStep === APPLICATION_CHART_CREATE_STEPS.PREVIEW.step ? (
            <ApplicationChartAppPreview />
          ) : (
            <></>
          )}
        </div>
      </Card>
      <ApplicationChartCreateButtonBar />
    </div>
  );
};
export default ApplicationChartCreate;
