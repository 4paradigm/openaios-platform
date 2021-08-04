import React from 'react';
import { Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import { mount } from 'enzyme';
import AsideNav from '../index';
import { useHistory } from 'react-router-dom';

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const menuState = mockStore({
  common: {
    isMobile: false,
  },
});

describe('Test for <AsideNav /> component', () => {
  const history = useHistory();
  const asideNav = mount(
    <Provider store={menuState}>
      <AsideNav />
    </Provider>,
  );
  const devEnv = asideNav.find('.devEnv').hostNodes();
  const appManage = asideNav.find('.appManage').hostNodes();
  const appChart = asideNav.find('.appChart').hostNodes();
  const appInstance = asideNav.find('.appInstance').hostNodes();
  const file = asideNav.find('.file').hostNodes();
  const mirror = asideNav.find('.mirror').hostNodes();
  const priMirror = asideNav.find('.priMirror').hostNodes();
  const pubMirror = asideNav.find('.pubMirror').hostNodes();

  test('basic render', () => {
    // menu基础显示正常
    expect(devEnv).toHaveLength(1);
    expect(appManage).toHaveLength(1);
    expect(appChart).toHaveLength(1);
    expect(appInstance).toHaveLength(1);
    expect(file).toHaveLength(1);
    expect(mirror).toHaveLength(1);
    expect(priMirror).toHaveLength(1);
    expect(pubMirror).toHaveLength(1);
  });

  test('test for sub menu', () => {
    expect(appManage.find('.appChart').hostNodes()).toHaveLength(1);
    expect(appManage.find('.appInstance').hostNodes()).toHaveLength(1);
    expect(mirror.find('.priMirror').hostNodes()).toHaveLength(1);
    expect(mirror.find('.pubMirror').hostNodes()).toHaveLength(1);
  });

  test('test for the function called handleClick', () => {
    devEnv.simulate('click');
    expect(history.location.pathname).toEqual('/devEnvironment');

    appChart.simulate('click');
    expect(history.location.pathname).toEqual('/application_chart');

    appInstance.simulate('click');
    expect(history.location.pathname).toEqual('/application_instance');

    file.simulate('click');
    expect(history.location.pathname).toEqual('/file');

    priMirror.simulate('click');
    expect(history.location.pathname).toEqual('/private_mirror');

    pubMirror.simulate('click');
    expect(history.location.pathname).toEqual('/public_mirror');
  });
});
