/*
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-24 17:51:41
 * @Description: file content
 */
export enum ENVIRONMENT_STATUS {
  Pending = 'Pending',
  Running = 'Running',
  Succeeded = 'Succeeded',
  Failed = 'Failed',
  Unknown = 'Unknown',
  Killed = 'Killed',
}
export const ENVIRONMENT_STATUS_TAG: any = {
  Pending: 'default',
  Running: 'processing',
  Succeeded: 'success',
  Failed: 'error',
  Killed: 'error',
  Unknown: 'warning',
};

export const FILE_ICON: any = {
  VIDEO: 'shipin',
  DOC: 'wendang',
  IMAGE: 'tuxiang',
  AUDIO: 'yinpin',
  EXCEL: 'Excel',
  FLASH: 'Flash',
  PDF: 'PDF',
  PPT: 'PPT',
  WORD: 'Word',
  OTHERS: 'cankaowenjian',
};
