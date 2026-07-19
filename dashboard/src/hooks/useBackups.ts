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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useUpdate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; payload: Parameters<typeof backupsService.update>[1] }) =>
      backupsService.update(payload.id, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDelete = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; projectId: string }) =>
      backupsService.delete(payload.id, payload.projectId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useTrigger = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.trigger(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDeleteRecord = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; recordId: string }) =>
      backupsService.deleteRecord(payload.id, payload.recordId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useRestore = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => backupsService.restore(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};

export const useDeleteS3Destination = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; projectId: string }) =>
      backupsService.deleteS3Destination(payload.id, payload.projectId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['backups'] });
    },
  });
};
