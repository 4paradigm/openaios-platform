/*
 * @Author: liyuying
 * @Date: 2021-04-25 11:52:23
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-17 19:34:08
 * @Description: 公告信息查看弹窗
 */
import React, { useState, useRef, useEffect } from 'react';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import ReactMarkdown from '@uiw/react-markdown-preview';
import { Button, Icon, Form, Input, Select, Radio, Breadcrumb, message } from 'cess-ui';
import { useSelector, useDispatch, IMessageDetailState, MessageDetailAction } from 'umi';
import keycloakClient from '@/keycloak';
import { INDEX } from '@/router/url';

import './index.less';

let timer: any = { id: '' };
const breadCrumb = (
  <Breadcrumb>
    <Breadcrumb.Item href={INDEX}>首页</Breadcrumb.Item>
    <Breadcrumb.Item>活动与赛事</Breadcrumb.Item>
  </Breadcrumb>
);
const MessageDetail = ({ match, location }: any) => {
  const dispatch = useDispatch();
  const formRef = useRef(null);
  const [countDownDay, setCountDownDay] = useState(0);
  const [countDownTime, setCountDownTime] = useState('00:00:00');
  /* 允许初始化环境 */
  const [isCanStartEnv, setIsCanStartEnv] = useState(false);
  /* 报名已截止 */
  const [isRegistionEnd, setIsRegistionEnd] = useState(false);
  /* 报名未开始 */
  const [isRegistioNotStart, setRegistioNotStart] = useState(false);
  /* 邀请的url */
  const inviteUrl = `${window.location.origin}${
    window.location.pathname
  }?user=${keycloakClient.getUserId()}`;
  const { msgInfo, modalLoading, hasApplied, invitePersonNum }: IMessageDetailState = useSelector(
    (state: any) => state.messageDetail,
  );
  /**
   * 初始化环境
   */
  const handelInitEnv = () => {
    dispatch({
      type: MessageDetailAction.INIT_ENV,
      payload: msgInfo,
    });
  };
  /**
   * 报名
   */
  const handelApply = () => {
    if (formRef.current as any) {
      (formRef.current as any)
        .validateFields()
        .then((values: any) => {
          console.log(values);
          dispatch({
            type: MessageDetailAction.SIGN_UP,
            payload: {
              competitionID: msgInfo.id,
              formData: values,
              inviter: location.query.user || '',
            },
          });
        })
        .catch(() => {
          const messageDetailDiv = document.getElementById('view-msg-content');
          if (messageDetailDiv) {
            messageDetailDiv.scrollTop = messageDetailDiv.scrollHeight || 0;
          }
        });
    } else {
      dispatch({
        type: MessageDetailAction.SIGN_UP,
        payload: {
          competitionID: msgInfo.id,
        },
      });
    }
  };

  // 倒计时
  const countDown = (startTime: any, endTime: any) => {
    const usedTime = Math.floor(
      (new Date(endTime).getTime() - new Date(startTime).getTime()) / 1000,
    );
    if (usedTime <= 0) {
      setCountDownDay(0);
      setCountDownTime('00:00:00');
      setIsCanStartEnv(true);
      return;
    }
    let day: number | string = Math.floor(usedTime / (60 * 60 * 24));
    let hour: number | string = Math.floor(usedTime / (60 * 60)) - day * 24;
    let minute: number | string = Math.floor(usedTime / 60) - day * 24 * 60 - hour * 60;
    let second: number | string =
      Math.floor(usedTime) - day * 24 * 60 * 60 - hour * 60 * 60 - minute * 60;
    if (hour <= 9) hour = '0' + hour;
    if (minute <= 9) minute = '0' + minute;
    if (second <= 9) second = '0' + second;
    setCountDownDay(day);
    setCountDownTime(`${hour}:${minute}:${second}`);
    const timeOut = setTimeout(() => {
      countDown(startTime + 1000, endTime);
    }, 1000);
    timer.id = timeOut;
  };
  const handleCopy = (path: string) => {
    message.success(`${path}，已复制到剪贴板`);
  };
  useEffect(() => {
    /**
     * 校验当前人员在该比赛中的状态
     */
    dispatch({
      type: MessageDetailAction.INIT_DATA,
      payload: match.params.id,
    });
  }, []);
  useEffect(() => {
    // 已获取详情
    if (msgInfo.id) {
      // 判断报名截止
      if (new Date(msgInfo.deadline || '').getTime() < new Date().getTime()) {
        setIsRegistionEnd(true);
      } else if (new Date(msgInfo.beginning || '').getTime() > new Date().getTime()) {
        // 判断报名未开始
        setRegistioNotStart(true);
        countDown(new Date().getTime(), msgInfo.beginning);
      } else {
        setRegistioNotStart(false);
        setIsRegistionEnd(false);
      }
      // 判断可以开始初始化环境
      if (msgInfo.avl) {
        setIsCanStartEnv(true);
      } else {
        setIsCanStartEnv(false);
      }
      // if (msgInfo.initBeginning) {
      //   setIsCanStartEnv(false);
      //   countDown(new Date().getTime(), msgInfo.initBeginning);
      // } else {
      //   setIsCanStartEnv(true);
      // }
    }
    return () => {
      if (timer && timer.id) {
        clearTimeout(timer.id);
        timer = null;
      }
    };
  }, [msgInfo]);
  return (
    <div className="message-detail-page">
      {breadCrumb}
      <div className="view-msg-container">
        <div className="view-msg-content" id="view-msg-content">
          <h2>{msgInfo.name}</h2>
          <div className="invite-bar">
            <div className="left">
              <div className="invite-title">邀请参赛拿大奖</div>
              <label>邀请链接</label>
              <Input readOnly value={inviteUrl}></Input>
              <CopyToClipboard text={inviteUrl} onCopy={() => handleCopy(inviteUrl)}>
                <Button
                  onClick={(e) => {
                    e.stopPropagation();
                  }}
                  type="primary"
                >
                  复制并发送给他人
                </Button>
              </CopyToClipboard>
              <div className="invite-person">
                目前您在此比赛已经邀请成功<label className="person-num">{invitePersonNum}</label>
                人次
              </div>
            </div>
            <div className="right">
              <ReactMarkdown source={msgInfo.ruleMd} linkTarget={'_blank'} />
            </div>
          </div>
          <ReactMarkdown source={msgInfo.descriptionMd} linkTarget={'_blank'} />
          {msgInfo && msgInfo.formConfig && !hasApplied && !isRegistionEnd ? (
            <div className="game-form">
              <div className="game-form-title">调查问卷</div>
              <Form ref={formRef} requiredMark={false} layout="vertical" preserve={false}>
                {msgInfo.formConfig.map((item) => {
                  return (
                    <Form.Item label={item.lable} rules={item.rules} name={item.key} key={item.key}>
                      {item.type === 'INPUT' ? (
                        <Input
                          maxLength={item.maxLength || 30}
                          placeholder={item.placeholder}
                          autoComplete="off"
                        />
                      ) : item.type === 'SELECT' ? (
                        <Select
                          showSearch
                          style={{ width: '842px' }}
                          placeholder={item.placeholder}
                        >
                          {item.options?.map((item) => {
                            return (
                              <Select.Option key={item} value={item}>
                                {item}
                              </Select.Option>
                            );
                          })}
                        </Select>
                      ) : item.type === 'RADIO' ? (
                        <Radio.Group>
                          {item.options?.map((item) => {
                            return (
                              <Radio key={item} value={item}>
                                {item}
                              </Radio>
                            );
                          })}
                        </Radio.Group>
                      ) : (
                        ''
                      )}
                    </Form.Item>
                  );
                })}
              </Form>
            </div>
          ) : (
            ''
          )}
        </div>
      </div>
      <div className="view-msg-operation-bar">
        {hasApplied ? (
          <>
            <label className="success-info">
              您已报名成功
              <Icon type="success" />
            </label>
            {msgInfo && !isCanStartEnv ? (
              // <div>
              //   距离赛事启动还有<span className="count-down-num">{countDownDay}</span>天
              //   <span className="count-down-num">{countDownTime}</span>
              //   ，请您耐心等待！
              // </div>
              <div>赛事还未启动 ，请您耐心等待！</div>
            ) : (
              <>
                初始化过程会为您申请计算资源，请确保您的参赛时间，谨慎选择！
                <Button type="primary" loading={modalLoading} onClick={handelInitEnv}>
                  一键初始化环境
                </Button>
              </>
            )}
          </>
        ) : msgInfo && !isRegistionEnd && !isRegistioNotStart ? (
          <Button type="primary" loading={modalLoading} onClick={handelApply}>
            我要报名
          </Button>
        ) : isRegistionEnd ? (
          '报名已截止'
        ) : (
          <div>
            距离赛事报名开始还有<span className="count-down-num">{countDownDay}</span>天
            <span className="count-down-num">{countDownTime}</span>
            ，请您耐心等待！
          </div>
        )}
      </div>
    </div>
  );
};
export default MessageDetail;
