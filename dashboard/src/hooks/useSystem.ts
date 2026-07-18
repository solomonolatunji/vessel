import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { RailwayImportRequest } from '#/interfaces/system';
import { systemService } from '#/services/system';

export const useGetSystemStats = () => {
  return useQuery({
    queryKey: ['system', 'stats'],
    queryFn: () => systemService.getSystemStats(),
    refetchInterval: 30_000,
  });
};

export const useExportSystem = () => {
  return useMutation({
    mutationFn: (passphrase: string) => systemService.exportSystem(passphrase),
  });
};

export const useImportSystem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: FormData) => systemService.importSystem(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};

export const useRestartSystem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => systemService.restartSystem(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};

export const useCleanupSystem = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => systemService.cleanupSystem(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};

export const useGetRailwayProjects = (token: string, enabled = false) => {
  return useQuery({
    queryKey: ['system', 'migration', 'railway', 'projects', token],
    queryFn: () => systemService.getRailwayProjects(token),
    enabled: enabled && !!token,
  });
};

export const useImportRailwayProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: RailwayImportRequest) => systemService.importRailwayProject(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['system'] });
    },
  });
};
