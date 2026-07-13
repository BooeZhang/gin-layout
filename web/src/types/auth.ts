export interface TokenPayload {
  accessToken?: string;
  refreshToken?: string;
  access_token?: string;
  refresh_token?: string;
  expiresIn?: number;
  tokenType?: string;
}

// 登录
export interface LoginPayload {
  account: string; // 账号
  password: string; // 密码
  captcha?: string; // 验证码
}

export interface PasswordPayload {
  oldPassword?: string;
  newPassword?: string;
  password?: string;
  [key: string]: unknown;
}
