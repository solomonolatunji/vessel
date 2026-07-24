import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { scheduledTasksService } from '#/services/scheduled-tasks';

export const useListScheduledTasks = (serviceId: string) => {
  return useQuery({
    queryKey: ['scheduled-tasks', serviceId],
    queryFn: () => scheduledTasksService.listScheduledTasks(serviceId),
  });
};

export const useGetScheduledTask = (id: string) => {
  return useQuery({
    queryKey: ['scheduled-tasks', 'getScheduledTask', id].filter(Boolean),
    queryFn: () => scheduledTasksService.getScheduledTask(id),
  });
};

export const useCreateScheduledTask = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      payload: Parameters<typeof scheduledTasksService.createScheduledTask>[0];
    }) => scheduledTasksService.createScheduledTask(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['scheduled-tasks'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteScheduledTask = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => scheduledTasksService.deleteScheduledTask(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['scheduled-tasks'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useTriggerScheduledTask = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => scheduledTasksService.triggerScheduledTask(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['scheduled-tasks'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};
