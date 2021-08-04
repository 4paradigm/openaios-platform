import React from 'react';
import { Provider } from 'react-redux';
import { BrowserRouter as Router } from 'react-router-dom';
import configureMockStore from 'redux-mock-store';
import { mount } from 'enzyme';
import Home from '../index';
import MessageDetail from '../MessageDetail/index';

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const homeState = mockStore({
  common: {
    isMobile: false,
  },
  home: {
    taskInfo: {
      app_num: 2,
      env_num: 1,
      perCost: '0.000',
      task_list: [],
    },
    msgInfo: [
      {
        name: '【AI应用与异构内存编程挑战赛】编程挑战赛道',
        participant: 307,
        title: '【比赛】【AI应用与异构内存编程挑战赛】编程挑战赛道',
        id: 'intel-ai-coding-20210601',
        deadline: '2021-08-31T16:00:00Z',
        beginning: '2021-05-31T13:00:00Z',
        avl: true,
        descriptionMd:
          '# [ai-coding] AI 应用与异构内存编程挑战赛 - 编程挑战赛道\n\n[☞ 点此阅读赛题](https://memark.io/ai_2021/devel_track.html)',
        computingResource: [
          {
            description: '4 cores, 8Gi mem, 64Gi PMEM',
            id: 'PMEM-competition',
            price: 0.004,
          },
        ],
      },
    ],
    userInfo: {
      balance: '740.888',
      costTime: '+ ∞',
      name: 'liuxueping',
    },
  },
  messageDetail: {
    msgInfo: {
      avl: true,
      beginning: '2021-05-31T13:00:00Z',
      computingResource: [
        {
          description: '4 cores, 8Gi mem, 64Gi PMEM',
          id: 'PMEM-competition',
          price: 0.004,
        },
      ],
      deadline: '2021-08-31T16:00:00Z',
      descriptionMd:
        '# [ai-coding] AI 应用与异构内存编程挑战赛 - 编程挑战赛道\n\n[☞ 点此阅读赛题](https://memark.io/ai_2021/devel_track.html)',
      formConfig: [
        {
          key: 'organization',
          lable: '学校或者工作单位',
          placeholder: '请输入学校或者工作单位',
          rules: [
            {
              message: '请输入学校或者工作单位',
              required: true,
            },
          ],
        },
        {
          key: 'phone',
          lable: '联系电话（为方便赛事相关事宜联系，建议填写）',
          placeholder: '请输入联系电话',
          rules: [],
          type: 'INPUT',
        },
        {
          key: 'PMem',
          lable: '对持久内存（PMem）的了解程度',
          placeholder: '请选择对持久内存（PMem）的了解程度',
          options: [
            '从未听说',
            '了解大致概念',
            '有过使用或者编程持久内存的初级经验',
            '非常了解，有持久内存相关的实际项目经验',
          ],
          rules: [{ required: true, message: '请选择对持久内存（PMem）的了解程度' }],
          type: 'RADIO',
        },
      ],
      id: 'intel-ai-coding-20210601',
      initEnvJson: {
        config: {},
        name: 'ai-coding',
      },
      name: '【AI应用与异构内存编程挑战赛】编程挑战赛道',
      participant: 307,
      ruleMd: '活动规则，此处是活动规则，此处是活动规则',
      title: '【比赛】【AI应用与异构内存编程挑战赛】编程挑战赛道',
    },
    ssh: {
      enable: false,
      'id_rsa.pub': '',
    },
    modalLoading: false,
    hasApplied: true,
    invitePersonNum: 0,
    hasData: true,
  },
});

const match = {
  isExact: true,
  path: '/home/message/:id',
  url: '/home/message/intel-ai-coding-20210601',
  params: {
    id: 'intel-ai-coding-20210601',
  },
};

const location = {
  search: '',
  state: undefined,
  query: {},
  pathname: '/home/message/intel-ai-coding-20210601',
  hash:
    '#state=d99e4a02-2880-45c6-87e4-7b97f17d03dc&session_state=1e21ff95-52ec-435c-8b9d-e940ace66f52&code=5678fb16-61d5-48bb-a6ea-9cb717929cc4.1e21ff95-52ec-435c-8b9d-e940ace66f52.97c4d8b5-5e50-4769-a87e-d5ff4edb7804',
};

describe('Test for <Home>', () => {
  test('basic render for <Home />', () => {
    const home = mount(
      <Provider store={homeState}>
        <Router>
          <Home />
        </Router>
      </Provider>,
    );
    // 活动与赛事模块
    expect(home.find('.card-title').first().text()).toContain('活动与赛事');
    expect(home.find('.home-msg-item').text()).toContain('AI应用与异构内存编程挑战赛');
    // 我的资源信息
    expect(home.find('.balance-item-label').first().text()).toContain('余额：');
    expect(home.find('.balance').first().text()).toContain('740.888');
    expect(home.find('.balance-item-label').at(1).text()).toContain('每分钟消耗：');
    expect(home.find('.balance').at(1).text()).toContain('0.000');
    expect(home.find('.balance-item-label').at(2).text()).toContain('预计可用：');
    expect(home.find('.balance').at(2).text()).toContain('+');
    // 我的实例信息
    expect(home.find('.env').text()).toContain('1');
    expect(home.find('.app').text()).toContain('2');
  });

  test('basic render for <MessageDetail />', () => {
    const msgDetail = mount(
      <Provider store={homeState}>
        <MessageDetail match={match} locatio={location} />
      </Provider>,
    );

    expect(msgDetail.find('#view-msg-content > h2').text()).toContain(
      '【AI应用与异构内存编程挑战赛】编程挑战赛道',
    );
    expect(msgDetail.find('.participant').text()).toContain('307');

    expect(msgDetail.find('.invite-title').text()).toContain('邀请参赛拿大奖');
    expect(msgDetail.find('.left > label').text()).toContain('邀请链接');

    expect(msgDetail.find('.invite-person .person-num').text()).toContain('0');
    expect(msgDetail.find('.right').text()).toContain('活动规则');
    expect(msgDetail.find('.description').hostNodes().text()).toContain('应用与异构内存编程挑战赛');

    // msgDetail.find('.select-ssh-1').hostNodes().simulate('click');

    // console.log(msgDetail.find('.id-rsa-pub-input').debug());
    // expect(msgDetail.find('.id-rsa-pub-input').exists()).toContain(true);
  });
});
