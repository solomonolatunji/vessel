import type { BaseResponse } from '#/interfaces/base';
import type { ServerSettings, UpdateSettingsRequest } from '#/interfaces/settings';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const settingsService = {
  getSettings: async (): Promise<BaseResponse<ServerSettings>> => {
    try {
      return await apiClient.get<BaseResponse<ServerSettings>>('/settings');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateSettings: async (payload: UpdateSettingsRequest): Promise<BaseResponse<ServerSettings>> => {
    try {
      return await apiClient.put<BaseResponse<ServerSettings>>('/settings', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotifications: async (): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.get<BaseResponse<unknown>>('/settings/notifications');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // TODO: Update exact shape
  saveNotificationChannel: async (payload: unknown): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.put<BaseResponse<unknown>>('/settings/notifications', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // TODO: Update exact shape
  testNotification: async (payload: unknown): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>('/settings/notifications/test', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteNotificationChannel: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/settings/notifications/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // --- GIT APPS ---

  getGitApps: async (provider: string): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.get<BaseResponse<unknown>>(`/settings/git_apps/${provider}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // TODO: Update exact shape
  saveGitApp: async (provider: string, payload: unknown): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.put<BaseResponse<unknown>>(`/settings/git_apps/${provider}`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteGitApp: async (provider: string, id: string): Promise<void> => {
    try {
      await apiClient.delete(`/settings/git_apps/${provider}/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // TODO: Update exact shape
  exchangeGithubManifest: async (payload: unknown): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>(
        '/settings/git_apps/github/manifest-callback',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // --- UPDATES & SYSTEM ---

  checkUpdate: async (): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>('/settings/updates/check');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployUpdate: async (): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>('/settings/updates/deploy');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  // TODO: Update exact shape
  activateLicense: async (payload: unknown): Promise<BaseResponse<unknown>> => {
    try {
      return await apiClient.post<BaseResponse<unknown>>('/settings/license', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
