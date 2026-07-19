import type { BaseResponse } from '#/interfaces/base';
import type { CreateJobRequest, Job } from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const jobsService = {
  listJobs: async (): Promise<BaseResponse<Job[]>> => {
    try {
      return await apiClient.get<BaseResponse<Job[]>>(`/jobs`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getJob: async (id: string): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.get<BaseResponse<Job>>(`/jobs/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createJob: async (payload: CreateJobRequest): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.post<BaseResponse<Job>>('/jobs', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteJob: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/jobs/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  triggerJob: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/jobs/${id}/trigger`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
