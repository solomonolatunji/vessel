import type { BaseResponse } from '#/interfaces/base';
import type { SystemStats } from '#/interfaces/system';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const systemService = {
  getSystemStats: async (): Promise<BaseResponse<SystemStats>> => {
    try {
      return await apiClient.get<BaseResponse<SystemStats>>('/system/stats');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  exportSystem: async (passphrase: string): Promise<Blob> => {
    try {
      return await apiClient.postBlob('/system/export', { passphrase });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  importSystem: async (payload: FormData): Promise<void> => {
    try {
      await apiClient.post('/system/import', payload, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  restartSystem: async (): Promise<void> => {
    try {
      await apiClient.post('/system/restart');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  cleanupSystem: async (): Promise<void> => {
    try {
      await apiClient.post('/system/maintenance/cleanup');
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
