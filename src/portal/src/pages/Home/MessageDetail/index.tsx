/*
 * @Author: liyuying
 * @Date: 2021-04-25 11:52:23
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2021-07-12 13:26:44
 * @Description: 公告信息查看弹窗
 */
import React, { useState, useRef, useEffect } from 'react';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import ReactMarkdown from '@uiw/react-markdown-preview';
import { Button, Icon, Form, Input, Select, Radio, Breadcrumb, message } from 'cess-ui';
import { IMessageDetailState, MessageDetailAction } from 'umi';
import { useSelector, useDispatch } from 'react-redux';
import keycloakClient from '@/keycloak';
import { INDEX } from '@/router/url';
import Loading from '@/components/Loading';

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

  /* 是否选择SSH这种环境交互方式 */
  const [selectSSH, setSelectSSH] = useState(0);

  /* 设置ssh的id_rsa.pub */
  const [idRsaPub, setIdRsaPub] = useState('');

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

  const {
    msgInfo,
    ssh,
    modalLoading,
    hasApplied,
    invitePersonNum,
    hasData,
  }: IMessageDetailState = useSelector((state: any) => state.messageDetail);

  // console.log('messageDetail-page-msgInfo:', msgInfo);

  const selectSSHHandler = (event: any) => {
    const selectSSHVal = event.target.value;
    setSelectSSH(selectSSHVal);
    // console.log( 'selectSSHVal:', selectSSHVal );
  };

  const idRsaPubHandler = (event: any) => {
    const idRsaPubVal = event.target.value;
    setIdRsaPub(idRsaPubVal);
    // console.log( 'idRsaPubVal:', idRsaPubVal );
  };

  /**
   * 初始化环境
   */
  const handelInitEnv = () => {
    // console.log('handelInitEnv--msgInfo:', msgInfo);
    // console.log('selectSSHVal:', selectSSH, selectSSH===1 );
    // console.log('idRsaPubVal:', idRsaPub);
    if (selectSSH === 1) {
      // 选择SSH的环境交互方式：true
      if (idRsaPub === '') {
        message.warning('已经选择SSH的环境交互方式，请先输入ssh的id_rsa.pub，再创建环境！');
        return;
      }
    }
    dispatch({
      type: MessageDetailAction.INIT_ENV,
      payload: {
        msgInfo,
        ssh: {
          enable: selectSSH === 1,
          'id_rsa.pub': idRsaPub,
        },
      },
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
      {hasData ? (
        <>
          <div className="view-msg-container">
            <div className="view-msg-content" id="view-msg-content">
              <h2>{msgInfo.name}</h2>
              {msgInfo.participant ? (
                <div className="participant">
                  已有<label className="person-num">{msgInfo.participant}</label>人参与当前赛事
                </div>
              ) : (
                ''
              )}
              <div className="invite-bar">
                <div className="left">
                  <div className="invite-title">邀请参赛拿大奖</div>
                  <label>邀请链接</label>
                  <Input className="invite_url" readOnly value={inviteUrl}></Input>
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
                    目前您在此比赛已经邀请成功
                    <label className="person-num">{invitePersonNum}</label>
                    人次
                  </div>
                </div>
                <div className="right">
                  <ReactMarkdown source={msgInfo.ruleMd} linkTarget={'_blank'} />
                </div>
              </div>
              <ReactMarkdown
                className="description"
                source={msgInfo.descriptionMd}
                linkTarget={'_blank'}
              />
              {msgInfo && msgInfo.formConfig && !hasApplied && !isRegistionEnd ? (
                <div className="game-form">
                  <div className="game-form-title">调查问卷</div>
                  <Form ref={formRef} requiredMark={false} layout="vertical" preserve={false}>
                    {msgInfo.formConfig.map((item) => {
                      return (
                        <Form.Item
                          label={item.lable}
                          rules={item.rules}
                          name={item.key}
                          key={item.key}
                        >
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
                    <div>
                      是否选择SSH的环境交互方式：
                      <Radio.Group onChange={selectSSHHandler}>
                        {/* {item.options?.map((item) => {
                          return ( */}
                        <Radio className="select-ssh-1" key={'select-ssh-1'} value={1}>
                          是
                        </Radio>
                        <Radio key={'select-ssh-0'} value={0}>
                          否
                        </Radio>
                        {/* );
                        })} */}
                      </Radio.Group>
                    </div>
                    <div>
                      {selectSSH ? (
                        <Input.TextArea
                          className="id-rsa-pub-input"
                          rows={3}
                          minLength={1}
                          maxLength={9999}
                          placeholder={`请输入ssh的id_rsa.pub`}
                          onChange={idRsaPubHandler}
                        />
                      ) : null}
                    </div>
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
        </>
      ) : (
        <Loading />
      )}
    </div>
  );
};
export default MessageDetail;
