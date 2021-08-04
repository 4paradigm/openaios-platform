import React from 'react';
import { Provider } from 'react-redux';
import configureMockStore from 'redux-mock-store';
import { mount, render } from 'enzyme';
import ComputeUnitRadio from '../index';

const conpute_unit_list = [
  {
    description: '2 gpus, 20 cores, 128Gi mem',
    id: '2GPU',
    name: '2GPU',
    price: 0.008,
  },
  {
    description: '4 cores, 8Gi mem, 64Gi PMEM',
    id: 'PMEM',
    name: 'PMEM',
    price: 0.004,
  },
  {
    description: '4 cores, 8Gi mem, 64Gi PMEM',
    id: 'PMEM-competition',
    name: 'PMEM-competition',
    price: 0.004,
  },
  {
    description: '1 core, 1Gi mem',
    id: 'single-core',
    name: 'single-core',
    price: 0.001,
  },
  {
    description: '3 cores, 12Gi mem',
    id: 'triple-core',
    name: 'triple-core',
    price: 0.003,
  },
];

const middlewares: any = [];
const mockStore = configureMockStore(middlewares);
const unitState = mockStore({
  computeUnitRadio: {
    computeUnitList: conpute_unit_list,
  },
});

describe('Test for <ComputeUnitRadio /> component', () => {
  const comUnit = mount(
    <Provider store={unitState}>
      <ComputeUnitRadio />
    </Provider>,
  );
  test('basic render', () => {
    expect(comUnit.find("[value='2GPU']").hostNodes()).toHaveLength(1);
    expect(comUnit.find("[value='PMEM']").hostNodes()).toHaveLength(1);
    expect(comUnit.find("[value='single-core']").hostNodes()).toHaveLength(1);
    expect(comUnit.find("[value='triple-core']").hostNodes()).toHaveLength(1);
  });
});
