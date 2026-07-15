import type { BaseResponse } from './base';

export type DatabaseEngine = 'postgres' | 'mysql' | 'redis' | 'mongodb' | 'mariadb';
export type DatabaseStatus = 'created' | 'running' | 'stopped' | 'error';
export type StorageType = 'minio' | 's3';
export type StorageStatus = 'running' | 'stopped' | 'error';

export interface Database {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  engine: DatabaseEngine;
  version: string;
  port: number;
  username: string;
  password: string;
  databaseName: string;
  volumePath: string;
  containerId: string;
  status: DatabaseStatus;
  internalDns: string;
  externalDns: string;
  customArgs: string;
  logicalReplication: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface CreateDatabaseRequest {
  projectId: string;
  environmentId: string;
  name: string;
  engine: DatabaseEngine;
  version: string;
  port: number;
  username: string;
  password: string;
  databaseName: string;
  volumePath: string;
  customArgs: string;
  logicalReplication: boolean;
}

export interface UpdateDatabaseRequest {
  externalDns: string;
  customArgs: string;
  logicalReplication: boolean;
}

export interface Storage {
  id: string;
  projectId: string;
  environmentId: string;
  name: string;
  type: StorageType;
  apiPort: number;
  consolePort: number;
  accessKey: string;
  secretKey?: string;
  bucketName: string;
  volumePath: string;
  containerId: string;
  status: StorageStatus;
  internalDns: string;
  externalDns: string;
  createdAt: string;
  updatedAt: string;
}

export interface DatabaseQueryRequest {
  query: string;
}

export interface DatabaseQueryResponse {
  columns?: string[];
  rows?: Record<string, unknown>[];
  result?: unknown;
}

export interface TableSchema {
  name: string;
  columns: ColumnSchema[];
}

export interface ColumnSchema {
  name: string;
  type: string;
  isNullable: boolean;
  isPrimary: boolean;
}

export interface ImportDatabaseRequest {
  sourceUrl: string;
}

export type TableRowPayload = Record<string, unknown>;

export type GetDatabasesResponse = BaseResponse<Database[]>;
export type GetDatabaseResponse = BaseResponse<Database>;
export type CreateDatabaseResponse = BaseResponse<Database>;
export type DatabaseQueryResponseType = BaseResponse<DatabaseQueryResponse>;
export type ListTablesResponse = BaseResponse<TableSchema[]>;
export type ImportDatabaseResponse = BaseResponse<void>;
export type DeleteDatabaseResponse = BaseResponse<void>;
export type GetStoragesResponse = BaseResponse<Storage[]>;
export type GetDatabaseStorageResponse = BaseResponse<Storage>;
