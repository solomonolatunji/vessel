import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { deploymentsService } from '#/services/deployments';

export const useListByService = (serviceId: string) => {
  return useQuery({
    queryKey: ['deployments', 'listByService', serviceId].filter(Boolean),
    queryFn: () => deploymentsService.listByService(serviceId),
  });
};

export const useTrigger = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      payload?: Parameters<typeof deploymentsService.trigger>[1];
    }) => deploymentsService.trigger(payload.serviceId, payload.payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['deployments'] });
    },
  });
};

export const useRollback = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { deploymentId: string }) =>
      deploymentsService.rollback(payload.deploymentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['deployments'] });
    },
  });
};

export const useGetLogs = (deploymentId: string) => {
  return useQuery({
    queryKey: ['deployments', 'getLogs', deploymentId].filter(Boolean),
    queryFn: () => deploymentsService.getLogs(deploymentId),
  });
};

export const useGetMetrics = (serviceId: string) => {
  return useQuery({
    queryKey: ['deployments', 'getMetrics', serviceId].filter(Boolean),
    queryFn: () => deploymentsService.getMetrics(serviceId),
  });
};

export const useDiagnostics = () => {
  return useMutation({
    mutationFn: (deploymentId: string) => deploymentsService.diagnostics(deploymentId),
  });
};
