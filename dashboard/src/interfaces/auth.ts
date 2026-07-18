import type { BaseResponse } from './base';
import type { User } from './users';

export interface AuthResponse {
  token: string;
  user: User;
}

export interface ApiErrorResponse {
  error: string;
}

export interface AuthCredentials {
  email: string;
  password: string;
}

export interface RegisterCredentials extends AuthCredentials {
  name: string;
}

export interface SetupEnv {
  jwtSecret?: string;
  dataDir: string;
  dashboardUrl: string;
  port: number;
}

export interface SetupCredentials extends RegisterCredentials {
  env: SetupEnv;
  defaultWildcardDomain?: string;
  s3Skip?: boolean;
  s3AccountId?: string;
  s3Bucket?: string;
  s3AccessKeyId?: string;
  s3SecretAccessKey?: string;
}

export interface Setup2FAResponse {
  secret: string;
  qrCodeUri: string;
  recoveryCodes: string[];
}

export interface Verify2FARequest {
  passcode: string;
}

export type LoginResponse = BaseResponse<AuthResponse>;
export type RegisterResponse = BaseResponse<AuthResponse>;
export type SetupResponse = BaseResponse<AuthResponse>;
export type ForgotPasswordResponse = BaseResponse<void>;
export type ResetPasswordResponse = BaseResponse<void>;
export type Setup2FAResponseType = BaseResponse<Setup2FAResponse>;
export type Verify2FAResponse = BaseResponse<void>;
export type Disable2FAResponse = BaseResponse<void>;
