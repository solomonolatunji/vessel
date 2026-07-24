import type { BaseResponse } from '#/interfaces/base';
import type {
  CreateProjectRequest,
  CreateProjectResponse,
  GetProjectResponse,
  ListProjectsResponse,
} from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const projectsService = {
  listProjects: async (): Promise<ListProjectsResponse> => {
    try {
      const url = '/projects';
      return await apiClient.get<ListProjectsResponse>(url);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getProject: async (id: string): Promise<GetProjectResponse> => {
    try {
      return await apiClient.get<GetProjectResponse>(`/projects/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createProject: async (payload: CreateProjectRequest): Promise<CreateProjectResponse> => {
    try {
      return await apiClient.post<CreateProjectResponse>('/projects', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteProject: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/projects/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getVars: async (id: string): Promise<BaseResponse<Record<string, string>>> => {
    try {
      return await apiClient.get<BaseResponse<Record<string, string>>>(`/projects/${id}/env`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  setVars: async (
    id: string,
    payload: { variables: Record<string, string> }
  ): Promise<BaseResponse<void>> => {
    try {
      return await apiClient.put<BaseResponse<void>>(`/projects/${id}/env`, payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
