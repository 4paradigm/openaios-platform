import React from 'react';
import { Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import { mount } from 'enzyme';
import MirrorPrivate from '../index';
import CopyImageModal from '../CopyImageModal/index';
import ImportModal from '../ImportModal/index';
import TaskModal from '../TaskModal/index';

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const priState = mockStore({
  privateImage: {
    dataSource: [
      {
        digest: 'sha256:c1e807f25111a1a658567c12df4e0bf8f5d99c9d6e2e5a043cee68ac355277ca',
        importing_time: '2021-07-05T04:52:01.894Z',
        repo: 'env/game/pmem-competition',
        size: '3.5G',
        tags: ['0.0.8'],
      },
      {
        digest: 'sha256:c1e807f25111a1a658567c12df4e0bf8f5d99c9d6e2e5a043cee68ac355277ca',
        importing_time: '2021-07-02T06:38:59.385Z',
        repo: 'pmem-competition',
        size: '3.5G',
        tags: ['0.0.8'],
      },
      {
        digest: 'sha256:2d8d27a6e324e38ca99072609185e08ae7f17256923c82d719b2d414322265e2',
        importing_time: '2021-07-05T04:43:48.179Z',
        repo: 'env/game/pmem-dev',
        size: '604.6M',
        tags: [],
      },
    ],
    total: 3,
    currentPage: 1,
    taskModalVisible: true,
    taskList: [
      {
        end_time: '2021-07-28T07:27:51Z',
        importing_id: 154,
        registry: { id: 1, url: 'https://hub.docker.com' },
        repo: 'library/python',
        start_time: '2021-07-28T07:21:58.609Z',
        status: 'Succeed',
        tag: 'latest',
      },
      {
        end_time: '2021-07-16T06:35:43.848Z',
        importing_id: 140,
        registry: { id: 1, url: 'https://hub.docker.com' },
        repo: 'library/python',
        start_time: '2021-07-16T06:35:32.817Z',
        status: 'NotFound',
        tag: 'stable,latest',
      },
    ],
    registryList: [
      {
        id: 1,
        url: 'https://hub.docker.com',
      },
    ],
    copyImageModalVisible: true,
    copyImageLoading: false,
    copyImageSourceImage: {
      digest: 'sha256:acb4e43d0c66c168e72ceaba5913cde472e4a17017cec9346969c9725a1fea94',
      importing_time: '2021-07-28T07:27:51.641Z',
      repo: 'library/python',
      size: '10.9G',
      tags: ['3.6'],
    },
    modalVisible: true,
    modalLoading: false,
  },
});

describe('Test for private mirror', () => {
  test('Basic render of table', () => {
    const priMirror = mount(
      <Provider store={priState}>
        <MirrorPrivate />
      </Provider>,
    );
    const allItems = priMirror.find('tr');
    // table 显示正常
    const titles = allItems.first();
    expect(titles.find('th').first().text()).toContain('仓库名');
    expect(titles.find('th').at(1).text()).toContain('Tags');
    expect(titles.find('th').at(2).text()).toContain('大小');
    expect(titles.find('th').at(3).text()).toContain('导入时间');
    expect(titles.find('th').at(4).text()).toContain('操作');

    const firstData = allItems.at(2);
    expect(firstData.find('td').first().text()).toContain('env/game/pmem-competition');
    expect(firstData.find('td').at(1).text()).toContain('0.0.8');
    expect(firstData.find('td').at(2).text()).toContain('3.5G');
    expect(firstData.find('td').at(3).text()).toContain('2021-07-05');
    expect(firstData.find('td').at(4).find('button')).toHaveLength(2);

    // 如果tags.length == 0 不显示copy，拿最后一条数据校验此特征
    expect(allItems.last().find('button').last().text()).not.toContain('copy');

    expect(priMirror.find('.task').hostNodes().text()).toContain('任务列表');
    expect(priMirror.find('.import').hostNodes().text()).toContain('导 入');
  });
  test('Basic render of <CopyImageModal>', () => {
    const copyImg = mount(
      <Provider store={priState}>
        <CopyImageModal />
      </Provider>,
    );
    const imgName = copyImg.find('.copy-image-item').first();
    const tag = copyImg.find('.copy-image-item').last();
    expect(imgName.children().first().text()).toContain('源仓库名称');
    expect(imgName.children().last().text()).toContain('library/python');
    expect(tag.children().first().text()).toContain('Tags');
    expect(tag.children().last().text()).toContain('3.6');
    expect(copyImg.find('input').prop('placeholder')).toContain('请输入目标仓库名称');

    // console.log(copyImg.find('.ant-btn-primary').hostNodes().last().debug());
    // copyImg.find('.ant-btn-primary').simulate('click');
    // console.log(copyImg.find("[role='alert']").debug());

    // 注意，此处没有测输入值后的变化，只检测了未输入值
    // expect(copyImg.find('button').last().simulate('click')).toContain("'name' is required");
    // expect(copyImg.find('input').prop('value')).toContain('things');
    // console.log(copyImg.text());
    // expect(copyImg.find('.ant-form-item-explain').exists()).toBe(true);
    // console.log(copyImg.find('.target').hostNodes().debug());
  });
  test('Basic render of <importModal>', () => {
    const importModal = mount(
      <Provider store={priState}>
        <ImportModal />
      </Provider>,
    );
    const allText = importModal.text();
    // 校验渲染
    expect(allText).toContain('registry');
    expect(allText).toContain('repo');
    expect(allText).toContain('tag');

    // importModal.find('#registryId').hostNodes().simulate('click');
    // console.log(importModal.find('.registry_list').debug());
    // 这里拿不到select的option，并不能自动渲染
    // expect(importModal.find('.registry_list').hostNodes().text()).toContain('https://hub.docker.com');
  });
  test('Basic render of <TaskModal>', () => {
    const taskModal = mount(
      <Provider store={priState}>
        <TaskModal />
      </Provider>,
    );
    const allItems = taskModal.find('tr');
    // table 显示正常
    const titles = allItems.first();
    expect(titles.find('th').first().text()).toContain('仓库名');
    expect(titles.find('th').at(1).text()).toContain('Tag');
    expect(titles.find('th').at(2).text()).toContain('状态');
    expect(titles.find('th').at(3).text()).toContain('创建时间');
    expect(titles.find('th').at(4).text()).toContain('结束时间');
    expect(titles.find('th').at(5).text()).toContain('Registry');
    expect(titles.find('th').at(6).text()).toContain('操作');

    const firstData = allItems.at(2);
    expect(firstData.find('td').first().text()).toContain('library/python');
    expect(firstData.find('td').at(1).text()).toContain('latest');
    expect(firstData.find('td').at(2).text()).toContain('Succeed');
    expect(firstData.find('td').at(3).text()).toContain('2021-07-28');
    expect(firstData.find('td').at(4).text()).toContain('2021-07-28');
    expect(firstData.find('td').at(5).text()).toContain('https://hub.docker.com');
  });
});
