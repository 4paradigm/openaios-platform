import React from 'react';
import { Button, CessCreateBtnBar, message, notification } from 'cess-ui';
import { APPLICATION_CHART_CREATE_STEPS } from '@/constant/application';
import {
  history,
  useDispatch,
  ApplicationChartCreateAction,
  IApplicationChartCreateState,
  useSelector,
  FileListActions,
} from 'umi';
import JsZip from 'jszip';
import { decode } from 'js-base64';
import { APPLICATION_CHART } from '@/router/url';

/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:56:09
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-03 18:09:59
 * @Description: 操作按钮栏目
 */
const { LeftBtnBar, RightBtnBar } = CessCreateBtnBar;

const ApplicationChartCreateButtonBar = () => {
  const dispatch = useDispatch();
  const { loading, currentStep, chartData }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );
  /**
   * 提交我的应用
   */
  const handelCreateMyApp = () => {
    let zip = new JsZip();
    zip.folder(`${chartData.metadata?.name}`);
    for (const filePath in chartData.files) {
      if (filePath !== 'answers.yaml') {
        zip.file(
          `${chartData.metadata?.name}/${filePath}`,
          decode(chartData.files[filePath] || ''),
        );
      }
    }
    zip
      .generateAsync({
        type: 'nodebuffer',
        compression: 'DEFLATE',
        compressionOptions: {
          level: 9,
        },
      })
      .then((content) => {
        dispatch({
          type: ApplicationChartCreateAction.CREATE_APPLICATION,
          payload: {
            file: new File([content], `${chartData.metadata?.name}.zip`),
          },
        });
      });
  };
  /**
   * 下一步
   */
  const handelNextStep = () => {
    if (currentStep === APPLICATION_CHART_CREATE_STEPS.UPLOAD.step) {
      // 需要上传应用包
      if (!chartData.metadata?.name) {
        notification.error({ message: '请上传应用' });
        return;
      }
    }
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { currentStep: currentStep + 1 },
    });
  };
  /**
   * 上一步
   */
  const handelPreStep = () => {
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { currentStep: currentStep - 1 },
    });
  };
  /**
   * 退出
   */
  const handelGoChartShop = () => {
    history.push(APPLICATION_CHART);
  };
  return (
    <CessCreateBtnBar>
      <LeftBtnBar>
        {currentStep === APPLICATION_CHART_CREATE_STEPS.UPLOAD.step ? (
          <Button onClick={handelGoChartShop}>退出</Button>
        ) : (
          <Button onClick={handelPreStep}>上一步</Button>
        )}
      </LeftBtnBar>
      <RightBtnBar>
        {currentStep === APPLICATION_CHART_CREATE_STEPS.PREVIEW.step ? (
          <Button type="primary" onClick={handelCreateMyApp} loading={loading}>
            提交应用
          </Button>
        ) : (
          <Button type="primary" onClick={handelNextStep}>
            下一步
          </Button>
        )}
      </RightBtnBar>
    </CessCreateBtnBar>
  );
};
export default ApplicationChartCreateButtonBar;
