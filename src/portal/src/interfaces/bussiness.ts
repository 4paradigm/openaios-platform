/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-10 16:19:52
 * @Description: file content
 */
import { ENVIRONMENT_STATUS } from '@/constant/environment';

// ssh
export interface ISsh {
  enable: boolean;
  'id_rsa.pub': string;
}

// jupyter
export interface Ijupyter {
  enable: boolean;
  token: string;
}

export interface IPath {
  subpath: string;
  mountpath: string;
}

export interface IImage {
  repo: string;
  tags: any;
  size?: number;
  importing_time?: string;
  source?: 'private' | 'public';
}

export interface ISshInfo {
  ssh_ip: string;
  ssh_port: string;
}

export interface IEnvironmentConfig {
  image: IImage;
  mounts: IPath[];
  compute_unit?: string; // 算力规格
  ssh: ISsh;
  jupyter: Ijupyter;
}

export interface IStaticInfo {
  name: string;
  create_tm: string;
  environmentConfig: IEnvironmentConfig;
}

// 开发环境
export interface IEnvironment {
  state: ENVIRONMENT_STATUS;
  staticInfo: IStaticInfo;
  sshInfo: ISshInfo;
}

// 提交开发环境对象(后端不统一，所以单独定义)
export interface IEnvironmentData {
  image: IImage;
  mounts: IPath[];
  compute_unit: string;
  ssh: ISsh;
  jupyter: Ijupyter;
}

// 文件管理
export interface IFile {
  name: string;
  is_dir: boolean;
  size: number;
  modification_time: string;
}

// 算力规格
export interface IComputeUnit {
  id: string;
  price: number;
  description: string;
}

// 镜像导入的任务
export interface ITask {
  start_time: string;
  end_time: string;
  importing_id: number;
  repo: string;
  status: string;
  tag: string;
  registry: {
    id: number;
    url: string;
  };
}
// 调查问卷的配置
export interface IFormConfig {
  key: string;
  lable: string;
  type: 'INPUT' | 'SELECT' | 'RADIO';
  placeholder: string;
  rules: any[];
  maxLength?: number;
  options?: any[];
}
// 消息通知
export interface IMessage {
  /* id */
  id: number | string;
  /* 后端的标记名称 */
  name: string;
  /* 展示title */
  title: string;
  /* 说明的markdown */
  descriptionMd: any;
  /* 规则md */
  ruleMd?: any;
  /* 是否可用 */
  avl: boolean;
  /* 需要填写的表单配置 */
  formConfig?: IFormConfig[];
  /* 报名开始时间 */
  beginning?: number | string;
  /* 报名截止日期 */
  deadline?: number | string;
  /* 初始化环境开始时间 */
  initBeginning?: number | string;
  /* 初始化环境结束时间 */
  initDeadline?: number | string;
  /* 初始化环境的JSON */
  initEnvJson?: { name: string; config: any };
}

// banner对象
export interface IBanner {
  /* id */
  id: number | string;
  /* 标题 */
  title: string;
  /* 图片 */
  image: any;
}
