import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { appsService } from '#/services/apps';

export const useListByProject = (projectId: string) => {
  return useQuery({
    queryKey: ['apps', 'listByProject', projectId].filter(Boolean),
    queryFn: () => appsService.listByProject(projectId),
  });
};

export const useListByEnvironment = (environmentId: string) => {
  return useQuery({
    queryKey: ['apps', 'listByEnvironment', environmentId].filter(Boolean),
    queryFn: () => appsService.listByEnvironment(environmentId),
  });
};

export const useGetApp = (appId: string) => {
  return useQuery({
    queryKey: ['apps', 'getApp', appId].filter(Boolean),
    queryFn: () => appsService.getApp(appId),
  });
};

export const useCreateApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      environmentId: string;
      payload: Parameters<typeof appsService.createApp>[1];
    }) => appsService.createApp(payload.environmentId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useUpdateApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      appId: string;
      payload: Parameters<typeof appsService.updateApp>[1];
    }) => appsService.updateApp(payload.appId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string }) => appsService.deleteApp(payload.appId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useStopApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string }) => appsService.stopApp(payload.appId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useRedeployApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string }) => appsService.redeployApp(payload.appId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useRestartApp = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string }) => appsService.restartApp(payload.appId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['apps'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useListVariables = (appId: string) => {
  return useQuery({
    queryKey: ['apps', 'variables', appId],
    queryFn: () => appsService.listVariables(appId),
    enabled: !!appId,
  });
};

export const useCreateVariable = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      appId: string;
      payload: Parameters<typeof appsService.createVariable>[1];
    }) => appsService.createVariable(payload.appId, payload.payload),
    onSuccess: async (_, { appId }) => {
      await queryClient.invalidateQueries({ queryKey: ['apps', 'variables', appId] });
    },
  });
};

export const useUpdateVariable = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      appId: string;
      varId: string;
      payload: Parameters<typeof appsService.updateVariable>[2];
    }) => appsService.updateVariable(payload.appId, payload.varId, payload.payload),
    onSuccess: async (_, { appId }) => {
      await queryClient.invalidateQueries({ queryKey: ['apps', 'variables', appId] });
    },
  });
};

export const useDeleteVariable = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string; varId: string }) =>
      appsService.deleteVariable(payload.appId, payload.varId),
    onSuccess: async (_, { appId }) => {
      await queryClient.invalidateQueries({ queryKey: ['apps', 'variables', appId] });
    },
  });
};

export const useListLogDrains = (appId: string) => {
  return useQuery({
    queryKey: ['apps', 'logDrains', appId],
    queryFn: () => appsService.listLogDrains(appId),
    enabled: !!appId,
  });
};

export const useCreateLogDrain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string; payload: any }) =>
      appsService.createLogDrain(payload.appId, payload.payload),
    onSuccess: async (_, { appId }) => {
      await queryClient.invalidateQueries({ queryKey: ['apps', 'logDrains', appId] });
    },
  });
};

export const useDeleteLogDrain = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { appId: string; drainId: string }) =>
      appsService.deleteLogDrain(payload.appId, payload.drainId),
    onSuccess: async (_, { appId }) => {
      await queryClient.invalidateQueries({ queryKey: ['apps', 'logDrains', appId] });
    },
  });
};
