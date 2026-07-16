import type { BaseResponse, PaginatedData } from './base';
import type { Database } from './database';
import type { AppService } from './deployment';
import type { Storage } from './storage';

export type MemberPermission = 'admin' | 'member' | 'viewer';
export type MemberStatus = 'pending' | 'accepted';
export type SSLCertStatus = 'pending' | 'issued' | 'failed';

export interface ProjectConfig {
  id: string;
  name: string;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentConfig {
  id: string;
  projectId: string;
  name: string;
  isDefault: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface DomainConfig {
  id: string;
  projectId: string;
  domainName: string;
  redirectTo?: string;
  sslCertStatus: SSLCertStatus;
  pathPrefix: string;
  createdAt: string;
  updatedAt: string;
}

export interface ServerlessFunctionCode {
  id: string;
  serviceId: string;
  runtime: string;
  codeContent: string;
  createdAt: string;
  updatedAt: string;
}

export interface CanvasSummary {
  id: string;
  name: string;
  description?: string;
  environmentsCount: number;
  appsCount: number;
  databasesCount: number;
  storageCount: number;
  onlineServices: number;
  totalServices: number;
  serviceIcons: string[];
  defaultEnvironment?: EnvironmentConfig;
  createdAt: string;
  updatedAt: string;
}

export interface EnvironmentCanvas {
  environment: EnvironmentConfig;
  apps: AppService[];
  databases: Database[];
  storage: Storage[];
}

export interface CreateProjectRequest {
  name: string;
  description?: string;
  repositoryUrl?: string;
  branch?: string;
  internalPort?: number;
  domain?: string;
}

export interface ProjectToken {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  tokenPrefix: string;
  scopes: string[];
  ipAllowlist: string[];
  expiresAt?: string;
  createdAt: string;
}

export interface ProjectMember {
  id: string;
  projectId: string;
  userId?: string;
  email: string;
  permission: MemberPermission;
  status: MemberStatus;
  invitedAt: string;
  acceptedAt?: string;
}

export interface CreateWebhookRequest {
  url: string;
  eventTypes: string[];
  includePrEnvironments: boolean;
}

export interface CreateTokenRequest {
  name: string;
  environmentId: string;
  scopes: string[];
  ipAllowlist?: string[];
  expiresAt?: string;
}

export interface AddMemberRequest {
  email: string;
  permission: MemberPermission;
}

export type ListProjectsResponse = BaseResponse<PaginatedData<ProjectConfig>>;
export type GetProjectResponse = BaseResponse<ProjectConfig>;
export type CreateProjectResponse = BaseResponse<ProjectConfig>;
export type ListEnvironmentsResponse = BaseResponse<EnvironmentConfig[]>;
export type CreateEnvironmentResponse = BaseResponse<EnvironmentConfig>;
export type GetCanvasResponse = BaseResponse<CanvasSummary[]>;
export type GetEnvironmentCanvasResponse = BaseResponse<EnvironmentCanvas>;
export type ListDomainsResponse = BaseResponse<DomainConfig[]>;
export type CreateDomainResponse = BaseResponse<DomainConfig>;
export type ListMembersResponse = BaseResponse<ProjectMember[]>;
export type ListTokensResponse = BaseResponse<ProjectToken[]>;
