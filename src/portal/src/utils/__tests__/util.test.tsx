import { dealBodyToQuery } from '../util';

describe('Test for util', () => {
  test('dealBodyToQuery', () => {
    const bodyPara = {
      name: 'Ming',
      age: '15',
    };
    const query = '?name=Ming&age=15';
    expect(dealBodyToQuery(bodyPara)).toEqual(query);
  });
});
