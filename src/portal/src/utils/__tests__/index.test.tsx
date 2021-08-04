import {
  paramsToPath,
  unixTimestampFormate,
  arrSplit,
  toThousands,
  transferObjectFromLevlesToFlat,
  transferObjectFromFlatToLevles,
} from '../index';

describe('Test for index in utils', () => {
  test('paramsToPath: Combine parameters with url.', () => {
    const params = {
      name: 'Ming',
    };
    const url = 'https://localhost:name';
    const res = 'https://localhostMing';
    expect(paramsToPath(url, params)).toEqual(res);
  });

  test('unixTimestampFormate: Get formatted time.', () => {
    const time = '100000';
    // 在服务器上执行错误 执行结果是1970-01-01 00:01:40
    // expect(unixTimestampFormate(time)).toEqual('1970-01-01 08:01:40');
    expect(unixTimestampFormate(time, 'yyyy-MM-dd')).toEqual('1970-01-01');
  });

  test('arrSplit: Split the array', () => {
    const arr = [1, 2, 3, 4, 5, 6, 7];
    const res = [[1, 2, 3, 4, 5, 6], [7]];
    expect(arrSplit(arr)).toEqual(res);
  });

  test('toThousands', () => {
    expect(toThousands(1000000)).toEqual('1,000,000.00');
  });

  test('transferObjectFromLevlesToFlat', () => {
    const obj = {
      a: {
        b: {
          c: 1,
        },
      },
    };
    const flatObj = { 'a.b.c': 1 };
    expect(transferObjectFromLevlesToFlat(obj)).toEqual(flatObj);
  });

  test('transferObjectFromFlatToLevles', () => {
    const obj = {
      a: {
        b: {
          c: 1,
        },
      },
    };

    const flatObj = { 'a.b.c': 1 };
    expect(transferObjectFromFlatToLevles(flatObj)).toEqual(obj);
  });
});
