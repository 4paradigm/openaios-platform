import { Modal } from 'cess-ui';
import React, { useEffect } from 'react';
import { encode } from 'js-base64';
import {
  useDispatch,
  IApplicationChartCreateState,
  useSelector,
  ApplicationChartCreateAction,
  ApplicationInstanceCreateAction,
} from 'umi';
import AppInstanceConfigForm from '../../../../../Instance/Create/components/ConfigForm';

/*
 * @Author: liyuying
 * @Date: 2021-06-07 13:59:12
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-07 14:05:04
 * @Description: 预览question.yaml配置的form
 */
const AppChartPerviewFormModal = () => {
  const dispatch = useDispatch();
  const { modalVisible, questionYaml, chartData }: IApplicationChartCreateState = useSelector(
    (state: any) => state.applicationChartCreate,
  );
  const handleCancel = () => {
    dispatch({
      type: ApplicationChartCreateAction.UPDATE_STATE,
      payload: { modalVisible: false },
    });
  };
  useEffect(() => {
    // 展示弹窗时，处理yaml
    if (modalVisible) {
      (chartData.files || {})[`questions.yml`] = encode(questionYaml);
      const template = chartData;
      dispatch({
        type: ApplicationInstanceCreateAction.INIT_DATA,
        payload: template,
      });
    }
  }, [questionYaml, modalVisible]);
  return (
    <Modal
      title="表单预览"
      visible={modalVisible}
      closable={true}
      footer={null}
      destroyOnClose={false}
      onCancel={handleCancel}
      centered
      className="app-container-group-events"
      width={880}
    >
      <div>
        <AppInstanceConfigForm></AppInstanceConfigForm>
      </div>
    </Modal>
  );
};
export default AppChartPerviewFormModal;
