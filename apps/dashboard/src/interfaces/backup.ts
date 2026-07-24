import type { BaseResponse } from './base';

export type BackupConfigStatus = 'active' | 'inactive';
export type BackupRecordStatus = 'running' | 'completed' | 'failed';

export interface BackupConfig {
  id: string;
  projectId: string;
  databaseId?: string;
  s3DestinationId?: string;
  name: string;
  description: string;
  dbUser: string;
  dbPassword?: string;
  backupEnabled: boolean;
  s3Enabled: boolean;
  disableLocal: boolean;
  schedule: string;
  timezone: string;
  timeout: number;
  retentionDays: number;
  maxBackups: number;
  maxStorageGb: number;
  status: BackupConfigStatus;
  createdAt: string;
  updatedAt: string;
}

export interface BackupRecord {
  id: string;
  backupConfigId: string;
  projectId: string;
  databaseId?: string;
  status: BackupRecordStatus;
  filePath: string;
  fileSizeBytes: number;
  s3Url?: string;
  logs: string;
  startedAt: string;
  completedAt: string;
}

export interface S3Destination {
  id: string;
  projectId: string;
  name: string;
  description: string;
  provider: string;
  endpoint: string;
  bucket: string;
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
  createdAt: string;
}

export interface CreateBackupConfigRequest {
  projectId: string;
  name: string;
  description: string;
  dbUser: string;
  dbPassword?: string;
  backupEnabled: boolean;
  s3Enabled: boolean;
  disableLocal: boolean;
  schedule: string;
  timezone: string;
  timeout: number;
  retentionDays: number;
  maxBackups: number;
  maxStorageGb: number;
  databaseId?: string;
  s3DestinationId?: string;
}

export interface CreateS3DestinationRequest {
  projectId: string;
  name: string;
  description: string;
  provider: string;
  endpoint: string;
  bucket: string;
  region: string;
  accessKeyId: string;
  secretAccessKey: string;
}

export type ListBackupsResponse = BaseResponse<BackupConfig[]>;
export type GetBackupResponse = BaseResponse<BackupConfig>;
export type CreateBackupResponse = BaseResponse<BackupConfig>;
export type ListBackupRecordsResponse = BaseResponse<BackupRecord[]>;
export type ListS3DestinationsResponse = BaseResponse<S3Destination[]>;
export type CreateS3DestinationResponse = BaseResponse<S3Destination>;
