import type {
  CreateDatabaseRequest,
  Database,
  DatabaseQueryRequest,
  DatabaseQueryResponse,
} from '#/interfaces/database';
import { apiClient } from './instance';

export const databasesService = {
  listDatabases: async (): Promise<Database[]> => {
    const { data } = await apiClient.get<Database[]>('/databases');
    return data;
  },

  getDatabase: async (id: string): Promise<Database> => {
    const { data } = await apiClient.get<Database>(`/databases/${id}`);
    return data;
  },

  createDatabase: async (payload: CreateDatabaseRequest): Promise<Database> => {
    const { data } = await apiClient.post<Database>('/databases', payload);
    return data;
  },

  deleteDatabase: async (id: string): Promise<void> => {
    await apiClient.delete(`/databases/${id}`);
  },

  startDatabase: async (id: string): Promise<void> => {
    await apiClient.post(`/databases/${id}/start`);
  },

  stopDatabase: async (id: string): Promise<void> => {
    await apiClient.post(`/databases/${id}/stop`);
  },

  queryDatabase: async (
    id: string,
    payload: DatabaseQueryRequest
  ): Promise<DatabaseQueryResponse> => {
    const { data } = await apiClient.post<DatabaseQueryResponse>(`/databases/${id}/query`, payload);
    return data;
  },
};
