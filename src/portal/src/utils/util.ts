/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-08 14:52:35
 * @Description: file content
 */
// 将body处理成query
export function dealBodyToQuery(body: any) {
  let paramStr = '?';
  for (const p in body) {
    paramStr += p + '=' + body[p] + '&';
  }
  paramStr = paramStr.substring(0, paramStr.length - 1);
  return paramStr;
}
