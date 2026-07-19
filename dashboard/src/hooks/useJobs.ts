import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { jobsService } from '#/services/jobs';

export const useListJobs = (projectId: string) => {
  return useQuery({
    queryKey: ['jobs', 'listJobs', projectId].filter(Boolean),
    queryFn: () => jobsService.listJobs(projectId),
  });
};

export const useGetJob = (id: string) => {
  return useQuery({
    queryKey: ['jobs', 'getJob', id].filter(Boolean),
    queryFn: () => jobsService.getJob(id),
  });
};

export const useCreateJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof jobsService.createJob>[0] }) =>
      jobsService.createJob(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => jobsService.deleteJob(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useTriggerJob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => jobsService.triggerJob(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['jobs'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};
