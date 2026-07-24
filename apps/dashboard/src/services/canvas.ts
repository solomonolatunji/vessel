import type { BaseResponse } from '#/interfaces/base';
import type { CanvasSummary, EnvironmentCanvas } from '#/interfaces/project';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const canvasService = {
  listCanvasSummaries: async (): Promise<BaseResponse<CanvasSummary[]>> => {
    try {
      return await apiClient.get<BaseResponse<CanvasSummary[]>>('/canvas/projects');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getCanvasSummary: async (projectId: string): Promise<BaseResponse<CanvasSummary>> => {
    try {
      return await apiClient.get<BaseResponse<CanvasSummary>>(`/projects/${projectId}/summary`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getEnvironmentCanvas: async (envId: string): Promise<BaseResponse<EnvironmentCanvas>> => {
    try {
      return await apiClient.get<BaseResponse<EnvironmentCanvas>>(`/environments/${envId}/canvas`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
