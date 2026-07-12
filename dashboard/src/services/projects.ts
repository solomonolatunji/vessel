import type { CreateProjectRequest, ProjectConfig } from '#/interfaces/project';
import { apiClient } from './instance';

export const projectsService = {
  listProjects: async (): Promise<ProjectConfig[]> => {
    const { data } = await apiClient.get<ProjectConfig[]>('/projects');
    return data;
  },

  getProject: async (id: string): Promise<ProjectConfig> => {
    const { data } = await apiClient.get<ProjectConfig>(`/projects/${id}`);
    return data;
  },

  createProject: async (payload: CreateProjectRequest): Promise<ProjectConfig> => {
    const { data } = await apiClient.post<ProjectConfig>('/projects', payload);
    return data;
  },

  deleteProject: async (id: string): Promise<void> => {
    await apiClient.delete(`/projects/${id}`);
  },

  deployProject: async (id: string): Promise<void> => {
    await apiClient.post(`/projects/${id}/deploy`);
  },
};
