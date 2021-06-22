/*
 * @Author: liyuying
 * @Date: 2021-05-26 15:45:48
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-10 15:38:10
 * @Description: 根据配置，渲染展示的组件
 */
import React from 'react';
import { useSelector } from 'umi';
import { Form, Input, InputNumber, Radio, Select } from 'cess-ui';
import { IQuestionItem } from '@/interfaces';
import { APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE } from '@/constant/application';
import ComputeUnitRadio from '@/components/ComputeUnitRadio';
import { IApplicationInstanceCreateState } from '@/pages/Application/models/application-instance-create';
import ConfigGroup from '../ConfigGroup';
import './index.less';

interface IProps {
  questionItem: IQuestionItem;
}

const ConfigItem = ({ questionItem }: IProps) => {
  const { applicationInstance }: IApplicationInstanceCreateState = useSelector(
    (state: any) => state.applicationInstanceCreate,
  );
  return (
    <>
      <Form.Item
        label={questionItem.label}
        extra={questionItem.description}
        name={questionItem.variable}
        rules={questionItem.rules}
        initialValue={questionItem.default || ''}
        className="app-config-item"
      >
        {questionItem.type === APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.STRING ? (
          <Input
            minLength={questionItem.min_length || 0}
            maxLength={questionItem.max_length || 256}
            placeholder={`请输入 ${questionItem.label}`}
          />
        ) : questionItem.type === APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.INT ? (
          <InputNumber
            value={questionItem.default}
            min={!questionItem.min && questionItem.min !== 0 ? -999999 : questionItem.min}
            max={!questionItem.max && questionItem.max !== 0 ? 999999 : questionItem.max}
          ></InputNumber>
        ) : questionItem.type === APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.MULTILINE ? (
          <Input.TextArea
            rows={4}
            minLength={questionItem.min_length || 0}
            maxLength={questionItem.max_length || 9999}
            placeholder={`请输入 ${questionItem.label}`}
          />
        ) : questionItem.type === APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.BOOLEAN ? (
          <Radio.Group value={Boolean(questionItem.default)}>
            <Radio value={true}>True</Radio>
            <Radio value={false}>False</Radio>
          </Radio.Group>
        ) : questionItem.type === APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.ENUM ? (
          <Select
            showSearch
            style={{ width: '516px' }}
            placeholder={`请选择 ${questionItem.label}`}
            value={questionItem.default}
          >
            {(questionItem.options || []).map((item) => {
              return (
                <Select.Option key={item} value={item}>
                  {item}
                </Select.Option>
              );
            })}
          </Select>
        ) : questionItem.type === APPLICATION_TEMPLATE_QUESTION_ITEM_TYPE.COMPUTE_UNIT ? (
          <Radio.Group>
            <ComputeUnitRadio></ComputeUnitRadio>
          </Radio.Group>
        ) : (
          <></>
        )}
      </Form.Item>
      {/* 处理子配置 */}
      {questionItem.subquestions &&
      applicationInstance.answers &&
      questionItem.show_subquestion_if === applicationInstance.answers[questionItem.variable] ? (
        <>
          <ConfigGroup
            questionsGroup={{ groupName: '', questionList: questionItem.subquestions }}
          ></ConfigGroup>
        </>
      ) : (
        <></>
      )}
    </>
  );
};
export default ConfigItem;
