import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { serviceVariablesService } from '#/services/service-variables';

export const serviceVariablesKeys = {
  all: ['service-variables'] as const,
  byService: (serviceId: string) => ['service-variables', serviceId] as const,
};

export const useListServiceVariables = (serviceId: string) => {
  return useQuery({
    queryKey: serviceVariablesKeys.byService(serviceId),
    queryFn: () => serviceVariablesService.list(serviceId),
    enabled: !!serviceId,
  });
};

export const useCreateServiceVariable = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ serviceId, payload }: { serviceId: string; payload: unknown }) =>
      serviceVariablesService.create(serviceId, payload),
    onSuccess: (_, { serviceId }) => {
      queryClient.invalidateQueries({
        queryKey: serviceVariablesKeys.byService(serviceId),
      });
    },
  });
};

export const useUpdateServiceVariable = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ serviceId, id, payload }: { serviceId: string; id: string; payload: unknown }) =>
      serviceVariablesService.update(serviceId, id, payload),
    onSuccess: (_, { serviceId }) => {
      queryClient.invalidateQueries({
        queryKey: serviceVariablesKeys.byService(serviceId),
      });
    },
  });
};

export const useDeleteServiceVariable = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ serviceId, id }: { serviceId: string; id: string }) =>
      serviceVariablesService.delete(serviceId, id),
    onSuccess: (_, { serviceId }) => {
      queryClient.invalidateQueries({
        queryKey: serviceVariablesKeys.byService(serviceId),
      });
    },
  });
};
