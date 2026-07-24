import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { serviceVarsService } from '#/services/services';

export const useList = (serviceId: string) => {
  return useQuery({
    queryKey: ['serviceVars', 'list', serviceId].filter(Boolean),
    queryFn: () => serviceVarsService.list(serviceId),
  });
};

export const useEnvSuggestions = (serviceId: string, enabled: boolean) => {
  return useQuery({
    queryKey: ['serviceVars', 'suggestions', serviceId].filter(Boolean),
    queryFn: () => serviceVarsService.getEnvSuggestions(serviceId),
    enabled: enabled && !!serviceId,
  });
};

export const useCreate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      payload: Parameters<typeof serviceVarsService.create>[1];
    }) => serviceVarsService.create(payload.serviceId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['serviceVars'] });
    },
  });
};

export const useUpdate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      id: string;
      payload: Parameters<typeof serviceVarsService.update>[2];
    }) => serviceVarsService.update(payload.serviceId, payload.id, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['serviceVars'] });
    },
  });
};

export const useDelete = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { serviceId: string; id: string }) =>
      serviceVarsService.delete(payload.serviceId, payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['serviceVars'] });
    },
  });
};
