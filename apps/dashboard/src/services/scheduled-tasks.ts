import type { BaseResponse } from '#/interfaces/base';
import type { CreateJobRequest, Job } from '#/interfaces/deployment';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const scheduledTasksService = {
  listScheduledTasks: async (serviceId: string): Promise<BaseResponse<Job[]>> => {
    try {
      return await apiClient.get<BaseResponse<Job[]>>(`/scheduled-tasks?serviceId=${serviceId}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  getScheduledTask: async (id: string): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.get<BaseResponse<Job>>(`/scheduled-tasks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createScheduledTask: async (payload: CreateJobRequest): Promise<BaseResponse<Job>> => {
    try {
      return await apiClient.post<BaseResponse<Job>>('/scheduled-tasks', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteScheduledTask: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/scheduled-tasks/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  triggerScheduledTask: async (id: string): Promise<void> => {
    try {
      await apiClient.post(`/scheduled-tasks/${id}/trigger`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
