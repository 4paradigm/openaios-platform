/*
 * @Author: liyuying
 * @Date: 2021-05-27 11:14:07
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-01 17:51:52
 * @Description: 分组的question
 */
import React from 'react';
import { useSelector } from 'umi';
import { Divider } from 'cess-ui';
import { IQuestionsGroup, IQuestionItem } from '@/interfaces';
import { IApplicationInstanceCreateState } from '@/pages/Application/models/application-instance-create';
import ConfigItem from '../ConfigItem';

import './index.less';

interface IProps {
  questionsGroup: IQuestionsGroup;
}

const ConfigGroup = ({ questionsGroup }: IProps) => {
  const { applicationInstance }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  const checkShow = (item: IQuestionItem) => {
    let isShow = true;
    if (item.showConditions && item.showConditions.length > 0) {
      if (!applicationInstance.answers) {
        return false;
      }
      item.showConditions.forEach((element) => {
        if (
          !applicationInstance.answers ||
          applicationInstance.answers[element.key] + '' !== element.value
        ) {
          isShow = false;
        }
      });
    } else {
      isShow = true;
    }
    return isShow;
  };
  return (
    <>
      {questionsGroup.groupName ? (
        <Divider className="config-group-lable">{questionsGroup.groupName}</Divider>
      ) : (
        <></>
      )}
      {questionsGroup.questionList.map((item) => {
        if (checkShow(item)) {
          return <ConfigItem questionItem={item} key={item.variable}></ConfigItem>;
        }
      })}
    </>
  );
};
export default ConfigGroup;
