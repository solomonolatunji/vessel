import type { BaseResponse } from './base';

export type UserRole = 'admin' | 'member' | 'viewer';

export interface User {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  totpEnabled: boolean;
  oauthProvider?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PersonalAccessToken {
  id: string;
  userId: string;
  name: string;
  prefix: string;
  expiresAt?: string;
  createdAt: string;
}

export interface CreatePATResponse {
  token: PersonalAccessToken;
  plain: string;
}

export interface UpdateUserRequest {
  email: string;
  name: string;
}

export interface CreatePATRequest {
  name: string;
}

export type GetUsersResponse = BaseResponse<User[]>;
export type GetUserResponse = BaseResponse<User>;
export type CreatePATResponseType = BaseResponse<CreatePATResponse>;
export type ListPATsResponse = BaseResponse<PersonalAccessToken[]>;
