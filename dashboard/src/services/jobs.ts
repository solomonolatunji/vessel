import type { CreateJobRequest, Job, UpdateJobRequest } from '#/interfaces/deployment';
import { apiClient } from './instance';

export const jobsService = {
  listJobs: async (): Promise<Job[]> => {
    const { data } = await apiClient.get<Job[]>('/jobs');
    return data;
  },

  getJob: async (id: string): Promise<Job> => {
    const { data } = await apiClient.get<Job>(`/jobs/${id}`);
    return data;
  },

  createJob: async (payload: CreateJobRequest): Promise<Job> => {
    const { data } = await apiClient.post<Job>('/jobs', payload);
    return data;
  },

  deleteJob: async (id: string): Promise<void> => {
    await apiClient.delete(`/jobs/${id}`);
  },

  triggerJob: async (id: string): Promise<void> => {
    await apiClient.post(`/jobs/${id}/trigger`);
  },
};
