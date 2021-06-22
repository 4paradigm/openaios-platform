export interface LoginData {
  username?: string;
  password: string;
  email?: string;
  code?: string;
}

export interface Userinfo {
  id: number;
  username: string;
  displayName: string;
  token: string;
  createTime: string;
  updateTime: string;
}
