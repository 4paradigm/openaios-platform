import React, { useEffect } from 'react';
import { encode } from 'js-base64';
import AppInstanceCreateContent from '../../../../Instance/Create/components/CreateContent';
import './index.less';
import {
  useDispatch,
  IApplicationChartCreateState,
  useSelector,
  ApplicationInstanceCreateAction,
} from 'umi';

/*
 * @Author: liyuying
 * @Date: 2021-06-03 17:25:23
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-07 16:36:57
 * @Description: 创建我的应用--预览与提交
 */
const ApplicationChartAppPreview = () => {
  const dispatch = useDispatch();
  const { questionYaml, chartData, appReadMeDoc }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );
  useEffect(() => {
    (chartData.files || {})[`questions.yml`] = encode(questionYaml);
    (chartData.files || {})[`app-readme.md`] = encode(appReadMeDoc);
    const template = chartData;
    dispatch({
      type: ApplicationInstanceCreateAction.INIT_DATA,
      payload: template,
    });
  }, [chartData, questionYaml, appReadMeDoc]);
  return (
    <div>
      <AppInstanceCreateContent></AppInstanceCreateContent>
    </div>
  );
};
export default ApplicationChartAppPreview;
