import type { BaseResponse } from '#/interfaces/base';
import type { CreateJobRequest, Job } from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const jobsService = {
  listJobs: async (): Promise<BaseResponse<Job[]>> => {
    try {
      return await apiClient.get<BaseResponse<Job[]>>(`/scheduled-tasks`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getJob: async (id: string): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.get<BaseResponse<Job>>(`/scheduled-tasks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createJob: async (payload: CreateJobRequest): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.post<BaseResponse<Job>>('/scheduled-tasks', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteJob: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/scheduled-tasks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  triggerJob: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/scheduled-tasks/${id}/trigger`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
