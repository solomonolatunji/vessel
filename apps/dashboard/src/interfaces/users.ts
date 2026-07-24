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
  lastLogin?: string;
  projectsCount?: number;
  servicesCount?: number;
  apiKeysCount?: number;
}

export interface PersonalAccessToken {
  id: string;
  userId: string;
  name: string;
  prefix: string;
  accessLevel: 'read' | 'read_write';
  projectScope: 'all' | 'specific';
  allowedProjects: string[] | null;
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
  accessLevel: 'read' | 'read_write';
  projectScope: 'all' | 'specific';
  allowedProjects: string[];
  expiresAt?: string;
}

export type GetUsersResponse = BaseResponse<User[]>;
export type GetUserResponse = BaseResponse<User>;
export type CreatePATResponseType = BaseResponse<CreatePATResponse>;
export type ListPATsResponse = BaseResponse<PersonalAccessToken[]>;
