import type { BaseResponse } from '#/interfaces/base';
import type { TeamNotificationChannel } from '#/interfaces/notification';
import type {
  BitbucketApp,
  CheckUpdateResponse,
  DeployUpdateResponse,
  ExchangeGithubManifestResponse,
  GetBitbucketAppsResponse,
  GetGithubAppsResponse,
  GetGitlabAppsResponse,
  GetServerSettingsResponse,
  GetUpdateStatusResponse,
  GitAppsManifestRequest,
  GithubApp,
  GitlabApp,
  SaveBitbucketAppResponse,
  SaveGithubAppResponse,
  SaveGitlabAppResponse,
  TestNotificationRequest,
  TestNotificationResponseType,
  UpdateServerSettingsResponse,
  UpdateSettingsRequest,
} from '#/interfaces/settings';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const settingsService = {
  getSettings: async (): Promise<GetServerSettingsResponse> => {
    try {
      return await apiClient.get<GetServerSettingsResponse>('/settings');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getSetupStatus: async (): Promise<BaseResponse<{ setupRequired: boolean; cwd?: string }>> => {
    try {
      return await apiClient.get<BaseResponse<{ setupRequired: boolean; cwd?: string }>>(
        '/system/setup-status'
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getPublicSettings: async (): Promise<
    BaseResponse<{
      registrationEnabled: boolean;
      siteName?: string;
      emailEnabled: boolean;
    }>
  > => {
    try {
      return await apiClient.get<
        BaseResponse<{
          registrationEnabled: boolean;
          siteName?: string;
          emailEnabled: boolean;
        }>
      >('/system/public');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateSettings: async (payload: UpdateSettingsRequest): Promise<UpdateServerSettingsResponse> => {
    try {
      return await apiClient.put<UpdateServerSettingsResponse>('/settings', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getAISettings: async (): Promise<BaseResponse<any>> => {
    try {
      return await apiClient.get<BaseResponse<any>>('/ai');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateAISettings: async (payload: Record<string, unknown>): Promise<BaseResponse<any>> => {
    try {
      return await apiClient.put<BaseResponse<any>>('/ai', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotificationSettings: async (): Promise<BaseResponse<any>> => {
    try {
      return await apiClient.get<BaseResponse<any>>('/notifications');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateNotificationSettings: async (
    payload: Record<string, unknown>
  ): Promise<BaseResponse<any>> => {
    try {
      return await apiClient.put<BaseResponse<any>>('/notifications', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  testNotification: async (
    payload: TestNotificationRequest | Record<string, unknown>
  ): Promise<TestNotificationResponseType> => {
    try {
      return await apiClient.post<TestNotificationResponseType>(
        '/settings/notifications/test',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getNotifications: async (): Promise<BaseResponse<TeamNotificationChannel[]>> => {
    try {
      return await apiClient.get<BaseResponse<TeamNotificationChannel[]>>(
        '/settings/notifications'
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveNotificationChannel: async (
    payload: Record<string, unknown>
  ): Promise<BaseResponse<TeamNotificationChannel>> => {
    try {
      return await apiClient.post<BaseResponse<TeamNotificationChannel>>(
        '/settings/notifications',
        payload
      );
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

  getGitApps: async (
    provider: string
  ): Promise<GetGithubAppsResponse | GetGitlabAppsResponse | GetBitbucketAppsResponse> => {
    if (provider === 'gitlab') return await settingsService.getGitlabApps();
    if (provider === 'bitbucket') return await settingsService.getBitbucketApps();
    return await settingsService.getGithubApps();
  },

  // biome-ignore lint/suspicious/noExplicitAny: Provider payload and response depend dynamically on provider string
  saveGitApp: async (provider: string, payload: any): Promise<any> => {
    if (provider === 'gitlab') return await settingsService.saveGitlabApp(payload);
    if (provider === 'bitbucket') return await settingsService.saveBitbucketApp(payload);
    return await settingsService.saveGithubApp(payload);
  },

  deleteGitApp: async (provider: string, id: string): Promise<void> => {
    if (provider === 'gitlab') return await settingsService.deleteGitlabApp(id);
    if (provider === 'bitbucket') return await settingsService.deleteBitbucketApp(id);
    return await settingsService.deleteGithubApp(id);
  },

  getGithubApps: async (): Promise<GetGithubAppsResponse> => {
    try {
      return await apiClient.get<GetGithubAppsResponse>('/settings/git_apps/github');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveGithubApp: async (payload: GithubApp): Promise<SaveGithubAppResponse> => {
    try {
      return await apiClient.put<SaveGithubAppResponse>('/settings/git_apps/github', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteGithubApp: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/settings/git_apps/github/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getGitlabApps: async (): Promise<GetGitlabAppsResponse> => {
    try {
      return await apiClient.get<GetGitlabAppsResponse>('/settings/git_apps/gitlab');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveGitlabApp: async (payload: GitlabApp): Promise<SaveGitlabAppResponse> => {
    try {
      return await apiClient.put<SaveGitlabAppResponse>('/settings/git_apps/gitlab', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteGitlabApp: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/settings/git_apps/gitlab/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getBitbucketApps: async (): Promise<GetBitbucketAppsResponse> => {
    try {
      return await apiClient.get<GetBitbucketAppsResponse>('/settings/git_apps/bitbucket');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  saveBitbucketApp: async (payload: BitbucketApp): Promise<SaveBitbucketAppResponse> => {
    try {
      return await apiClient.put<SaveBitbucketAppResponse>('/settings/git_apps/bitbucket', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteBitbucketApp: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/settings/git_apps/bitbucket/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  exchangeGithubManifest: async (
    payload: GitAppsManifestRequest | Record<string, unknown>
  ): Promise<ExchangeGithubManifestResponse> => {
    try {
      return await apiClient.post<ExchangeGithubManifestResponse>(
        '/settings/git_apps/github/manifest-callback',
        payload
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getUpdateStatus: async (): Promise<GetUpdateStatusResponse> => {
    try {
      return await apiClient.get<GetUpdateStatusResponse>('/settings/updates/status');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  checkUpdate: async (): Promise<CheckUpdateResponse> => {
    try {
      return await apiClient.post<CheckUpdateResponse>('/settings/updates/check');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployUpdate: async (): Promise<DeployUpdateResponse> => {
    try {
      return await apiClient.post<DeployUpdateResponse>('/settings/updates/deploy');
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
