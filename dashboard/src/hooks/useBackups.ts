import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { backupsService } from '#/services/backups';

export const useList = (projectId: string) => {
  return useQuery({
    queryKey: ['backups', 'list', projectId].filter(Boolean),
    queryFn: () => backupsService.list(projectId),
  });
};

export const useGet = (id: string) => {
  return useQuery({
    queryKey: ['backups', 'get', id].filter(Boolean),
    queryFn: () => backupsService.get(id),
    enabled: !!id,
  });
};

export const useCreate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof backupsService.create>[0] }) =>
      backupsService.create(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDelete = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.delete(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useTrigger = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.trigger(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useRestore = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.restore(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useListRecords = (id: string) => {
  return useQuery({
    queryKey: ['backups', 'listRecords', id].filter(Boolean),
    queryFn: () => backupsService.listRecords(id),
    enabled: !!id,
  });
};

export const useListS3Destinations = (projectId: string) => {
  return useQuery({
    queryKey: ['backups', 'listS3Destinations', projectId].filter(Boolean),
    queryFn: () => backupsService.listS3Destinations(projectId),
  });
};

export const useCreateS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof backupsService.createS3Destination>[0] }) =>
      backupsService.createS3Destination(payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDeleteS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.deleteS3Destination(payload.id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};
