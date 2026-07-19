import type {
  CreateBackupConfigRequest,
  CreateBackupResponse,
  CreateS3DestinationRequest,
  CreateS3DestinationResponse,
  GetBackupResponse,
  ListBackupRecordsResponse,
  ListBackupsResponse,
  ListS3DestinationsResponse,
} from '#/interfaces/backup';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const backupsService = {
  list: async (): Promise<ListBackupsResponse> => {
    try {
      return await apiClient.get<ListBackupsResponse>(`/backups`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  get: async (id: string): Promise<GetBackupResponse> => {
    try {
      return await apiClient.get<GetBackupResponse>(`/backups/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  create: async (payload: CreateBackupConfigRequest): Promise<CreateBackupResponse> => {
    try {
      return await apiClient.post<CreateBackupResponse>('/backups', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  update: async (id: string, payload: CreateBackupConfigRequest): Promise<CreateBackupResponse> => {
    try {
      return await apiClient.put<CreateBackupResponse>(`/backups/${id}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/backups/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  trigger: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/backups/${id}/trigger`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  restore: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/backups/${id}/restore`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteRecord: async (configId: string, recordId: string): Promise<void> => {
    try {
      await apiClient.delete(`/backups/${configId}/records/${recordId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listRecords: async (id: string): Promise<ListBackupRecordsResponse> => {
    try {
      return await apiClient.get<ListBackupRecordsResponse>(`/backups/${id}/records`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listS3Destinations: async (): Promise<ListS3DestinationsResponse> => {
    try {
      return await apiClient.get<ListS3DestinationsResponse>(`/s3-destinations`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createS3Destination: async (
    payload: CreateS3DestinationRequest
  ): Promise<CreateS3DestinationResponse> => {
    try {
      return await apiClient.post<CreateS3DestinationResponse>('/s3-destinations', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteS3Destination: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/s3-destinations/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
