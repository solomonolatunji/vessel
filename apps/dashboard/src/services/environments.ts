import type { CreateEnvironmentResponse, ListEnvironmentsResponse } from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const environmentsService = {
  listByProject: async (projectId: string): Promise<ListEnvironmentsResponse> => {
    try {
      return await apiClient.get<ListEnvironmentsResponse>(`/projects/${projectId}/environments`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createEnvironment: async (
    projectId: string,
    name: string
  ): Promise<CreateEnvironmentResponse> => {
    try {
      return await apiClient.post<CreateEnvironmentResponse>(
        `/projects/${projectId}/environments`,
        { name }
      );
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteEnvironment: async (environmentId: string): Promise<void> => {
    try {
      await apiClient.delete(`/environments/${environmentId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
