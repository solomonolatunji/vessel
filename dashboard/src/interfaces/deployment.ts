import type { BaseResponse } from './base';

export interface AppService {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  repositoryUrl: string;
  branch: string;
  rootDirectory: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: string;
  internalPort: number;
  domain: string;
  healthCheckPath: string;
  containerId: string;
  status: string;
  createdAt: string;
  updatedAt: string;
}

export interface Deployment {
  id: string;
  serviceId: string;
  environmentId: string;
  projectId: string;
  status: string;
  branch?: string;
  commitHash?: string;
  commitMessage?: string;
  trigger?: string;
  buildLogs?: string;
  containerId?: string;
  createdAt: string;
  updatedAt: string;
  finishedAt?: string;
}

export interface ServiceMetric {
  status: string;
  cpuUsagePercentage: number;
  memoryUsageBytes: number;
  memoryLimitBytes: number;
  uptimeSeconds: number;
}

export interface Variable {
  id: string;
  serviceId: string;
  projectId: string;
  environmentId: string;
  key: string;
  value: string;
  isSecret: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface Job {
  id: string;
  projectId: string;
  name: string;
  schedule: string;
  command: string;
  status: string;
  lastRunAt: string;
  lastOutput: string;
  createdAt: string;
  updatedAt: string;
}

export interface PRPreview {
  id: string;
  serviceId: string;
  projectId: string;
  prNumber: number;
  branch: string;
  commitHash: string;
  status: string;
  previewDomain: string;
  containerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAppServiceRequest {
  projectId: string;
  name: string;
  repositoryUrl: string;
  branch: string;
  rootDirectory: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: string;
  internalPort: number;
  domain: string;
  healthCheckPath: string;
}

export interface UpdateAppServiceRequest {
  name: string;
  repositoryUrl: string;
  branch: string;
  rootDirectory: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: string;
  internalPort: number;
  domain: string;
  healthCheckPath: string;
  containerId: string;
  status: string;
}

export interface TriggerDeploymentRequest {
  branch?: string;
}

export interface CreateServiceVarRequest {
  key: string;
  value: string;
  isSecret: boolean;
}

export interface UpdateServiceVarRequest {
  key: string;
  value: string;
  isSecret: boolean;
}

export interface CreateJobRequest {
  projectId: string;
  name: string;
  schedule: string;
  command: string;
}

export interface UpdateJobRequest {
  name?: string;
  schedule?: string;
  command?: string;
  status?: string;
}

// Response Types
export type ListAppsResponse = BaseResponse<AppService[]>;
export type GetAppResponse = BaseResponse<AppService>;
export type CreateAppResponse = BaseResponse<AppService>;
export type UpdateAppResponse = BaseResponse<AppService>;
