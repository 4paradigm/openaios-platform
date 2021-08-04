import React from 'react';
import { Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import { mount } from 'enzyme';
import keycloakClient from '@/keycloak';
import GlobalHeader from '../index';

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const menuState = mockStore({
  common: {
    isMobile: false,
  },
});

describe('Test for <GlobalHeader />', () => {
  const globalHeader = mount(
    <Provider store={menuState}>
      <GlobalHeader />
    </Provider>,
  );

  test('basic render with keycloak', () => {
    expect(globalHeader.find('header span').first().text()).toContain('社区版');
  });
});
