import type { BaseResponse } from './base';

export type DatabaseEngine =
  | 'postgres'
  | 'postgresql'
  | 'mysql'
  | 'redis'
  | 'mongodb'
  | 'mongo'
  | 'mariadb'
  | 'clickhouse'
  | 'kafka'
  | 'rabbitmq'
  | 'nats'
  | 'dragonfly'
  | 'keydb'
  | 'timescaledb';
export type DatabaseStatus = 'created' | 'running' | 'stopped' | 'error';

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

export interface OneClickApp {
  id: string;
  name: string;
  description: string;
  port: number;
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
