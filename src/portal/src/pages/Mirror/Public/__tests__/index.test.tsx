import React from 'react';
import { Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import { mount } from 'enzyme';
import MirrorPublic from '../index';

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const publicState = mockStore({
  publicImage: {
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
        tags: ['0.1.1'],
      },
    ],
    total: 3,
    currentPage: 1,
  },
});

describe('Test for public mirror', () => {
  const publicMirror = mount(
    <Provider store={publicState}>
      <MirrorPublic />
    </Provider>,
  );

  test('Basic render', () => {
    // console.log(publicMirror.find('th').hostNodes().length);
    expect(publicMirror.find('th').hostNodes()).toHaveLength(4);
  });
});
