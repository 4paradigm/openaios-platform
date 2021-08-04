/*
 * @Author: liyuying
 * @Date: 2021-06-20 12:56:30
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-22 16:01:51
 * @Description: file content
 */
import React from 'react';
import { useEffect } from 'react';
import { useDispatch, useSelector, AppServiceAction } from 'umi';
import { AppNotesAction, IAppNotesState } from './models/app-notes';
import './index.less';
import { Input, Empty } from 'cess-ui';

interface Iprop {
  instance_name: string;
}
/*
 * @Author: liyuying
 * @Date: 2021-06-01 17:46:21
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-02 14:16:52
 * @Description:实例详情--Notes
 */
const ApplicationInstanceDetailNotes = ({ instance_name }: Iprop) => {
  const dispatch = useDispatch();
  const { notes }: IAppNotesState = useSelector((state: any) => state.appNotes);
  const initNotes = async () => {
    await dispatch({
      type: AppServiceAction.GET_DATA,
      payload: instance_name,
    });
    await dispatch({
      type: AppNotesAction.GET_DATA,
      payload: instance_name,
    });
  };
  /**
   * 获取notes信息
   */
  useEffect(() => {
    if (instance_name) {
      initNotes();
    }
  }, [instance_name, dispatch]);
  return (
    <div className="application-instance-detail-notes">
      {notes ? (
        <Input.TextArea value={notes || ''} readOnly autoSize={true}></Input.TextArea>
      ) : (
        <div className="empty-container">
          <Empty msg="当前实例无 Notes"></Empty>
        </div>
      )}
    </div>
  );
};
export default ApplicationInstanceDetailNotes;
