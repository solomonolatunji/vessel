import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { serverlessService } from '#/services/serverless';

export const useGetCode = (serviceId: string) => {
  return useQuery({
    queryKey: ['serverless', 'getCode', serviceId].filter(Boolean),
    queryFn: () => serverlessService.getCode(serviceId),
  });
};

export const useSaveCode = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      serviceId: string;
      payload: { codeContent: string; runtime?: string };
    }) => serverlessService.saveCode(payload.serviceId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['serverless'] });
    },
  });
};
