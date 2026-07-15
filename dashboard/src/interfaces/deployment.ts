import type { BaseResponse, PaginatedData } from './base';

export type RuntimeMode = 'web' | 'worker';
export type BuildEngine =
  | 'nixpacks'
  | 'dockerfile'
  | 'buildpacks'
  | 'static'
  | 'railpack'
  | 'serverless';
export type ServiceStatus = 'created' | 'building' | 'running' | 'stopped' | 'error';
export type DeploymentStatus =
  | 'pending'
  | 'BUILDING'
  | 'PULLING'
  | 'CLONING'
  | 'FAILED'
  | 'SUCCESS';
export type JobStatus = 'active' | 'inactive';

export interface AppService {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  repositoryUrl: string;
  imageRef?: string;
  branch: string;
  rootDirectory: string;
  runtimeMode: RuntimeMode;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: BuildEngine;
  internalPort: number;
  domain: string;
  staticOutput: string;
  healthCheckPath: string;
  containerId: string;
  status: ServiceStatus;
  createdAt: string;
  updatedAt: string;
}

export interface Deployment {
  id: string;
  serviceId: string;
  environmentId: string;
  projectId: string;
  status: DeploymentStatus;
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
  timestamp: string;
  cpuPercent: number;
  memoryMB: number;
  networkRxKB: number;
  networkTxKB: number;
  status?: string;
  cpuUsagePercentage?: number;
  memoryUsageBytes?: number;
  memoryLimitBytes?: number;
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
  status: JobStatus;
  lastRunAt?: string;
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
  status: DeploymentStatus;
  previewDomain: string;
  containerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAppServiceRequest {
  projectId: string;
  name: string;
  repositoryUrl: string;
  imageRef?: string;
  branch: string;
  rootDirectory: string;
  runtimeMode: RuntimeMode;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: BuildEngine;
  internalPort: number;
  domain: string;
  staticOutput: string;
  healthCheckPath: string;
}

export interface UpdateAppServiceRequest {
  name: string;
  repositoryUrl: string;
  imageRef?: string;
  branch: string;
  rootDirectory: string;
  runtimeMode: RuntimeMode;
  installCommand: string;
  buildCommand: string;
  startCommand: string;
  dockerfilePath: string;
  buildEngine: BuildEngine;
  internalPort: number;
  domain: string;
  staticOutput: string;
  healthCheckPath: string;
  containerId: string;
  status: ServiceStatus;
}

export interface TriggerDeploymentRequest {
  commitId?: string;
}

export type DiagnosticsResponse = Record<string, Record<string, unknown>>;

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
  status?: JobStatus;
}

export type ListAppsResponse = BaseResponse<AppService[]>;
export type GetAppResponse = BaseResponse<AppService>;
export type CreateAppResponse = BaseResponse<AppService>;
export type UpdateAppResponse = BaseResponse<AppService>;

export type ListDeploymentsResponse = BaseResponse<PaginatedData<Deployment>>;
export type TriggerDeploymentResponse = BaseResponse<Deployment>;
export type RollbackDeploymentResponse = BaseResponse<Deployment>;
export type GetDeploymentLogsResponse = BaseResponse<string>;
export type GetServiceMetricsResponse = BaseResponse<ServiceMetric>;
export type GetDiagnosticsResponse = BaseResponse<DiagnosticsResponse>;

export type ListVariablesResponse = BaseResponse<Variable[]>;
export type SetVariablesResponse = BaseResponse<Variable[]>;
export type ListJobsResponse = BaseResponse<Job[]>;
export type GetJobResponse = BaseResponse<Job>;
export type CreateJobResponse = BaseResponse<Job>;
export type ListPRPreviewsResponse = BaseResponse<PRPreview[]>;
