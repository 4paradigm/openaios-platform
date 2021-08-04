import React from 'react';
import { Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import { mount } from 'enzyme';
import BasicLayout from '../WithAsideNav';

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const basicWebState = mockStore({
  common: {
    isMobile: false,
  },
});

const basicMobileState = mockStore({
  common: {
    isMobile: true,
  },
});

describe('BasicLayout with aside nav', () => {
  it('Test for basic render on mobile device', () => {
    const basicLayoutmobile = mount(
      <Provider store={basicMobileState}>
        <BasicLayout />
      </Provider>,
    );
    expect(basicLayoutmobile.find('.footer').exists()).toBe(false);
  });

  it('Test for basic render on web', () => {
    const basicLayout = mount(
      <Provider store={basicWebState}>
        <BasicLayout children={'test render'} />
      </Provider>,
    );
    expect(basicLayout.find('.content').text()).toContain('test render');
    expect(basicLayout.find('.footer').text()).toContain('联系我们');
    expect(basicLayout.find('.footer').exists()).toBe(true);
  });
});
