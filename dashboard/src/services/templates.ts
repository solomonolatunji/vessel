import type { OneClickApp } from '#/interfaces/database';
import type {
  ArchiveDeployResponse,
  ComposeDeployResponse,
  OneClickDeployRequest,
  OneClickDeployResponse,
} from '#/interfaces/templates';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const templatesService = {
  listOneClickApps: async (): Promise<Record<string, OneClickApp>> => {
    try {
      const response = await apiClient.get<{ data: Record<string, OneClickApp> }>('/one-click');
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployOneClickApp: async (payload: OneClickDeployRequest): Promise<OneClickDeployResponse> => {
    try {
      return await apiClient.post<OneClickDeployResponse>('/one-click/deploy', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployCompose: async (projectId: string, file: File): Promise<ComposeDeployResponse> => {
    try {
      const formData = new FormData();
      formData.append('projectId', projectId);
      formData.append('file', file);

      return await apiClient.post<ComposeDeployResponse>('/compose/deploy', formData);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deployArchive: async (projectId: string, file: File): Promise<ArchiveDeployResponse> => {
    try {
      const formData = new FormData();
      if (projectId) {
        formData.append('projectId', projectId);
      }
      formData.append('file', file);

      return await apiClient.post<ArchiveDeployResponse>('/deploy/archive', formData);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
