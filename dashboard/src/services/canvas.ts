import type { CanvasSummary, EnvironmentCanvas } from '#/interfaces/project';
import { apiClient } from './instance';

export const canvasService = {
  listCanvasSummaries: async (): Promise<CanvasSummary[]> => {
    const { data } = await apiClient.get<CanvasSummary[]>('/canvas/projects');
    return data;
  },

  getCanvasSummary: async (projectId: string): Promise<CanvasSummary> => {
    const { data } = await apiClient.get<CanvasSummary>(`/projects/${projectId}/summary`);
    return data;
  },

  getEnvironmentCanvas: async (envId: string): Promise<EnvironmentCanvas> => {
    const { data } = await apiClient.get<EnvironmentCanvas>(`/environments/${envId}/canvas`);
    return data;
  },
};
