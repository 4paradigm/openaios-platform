import { Model } from 'dva';
import { Action } from 'umi';

export interface IAction<T = any> extends Action<string> {
  type: string;
  payload?: T;
}

export interface IModel<T = any> extends Model {
  state: T;
}
