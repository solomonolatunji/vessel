import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { serverlessService } from '#/services/serverless';

export const useGetCode = (projectId: string | undefined, serviceId: string) => {
  return useQuery({
    queryKey: ['serverless', 'getCode', projectId, serviceId].filter(Boolean),
    queryFn: () => serverlessService.getCode(projectId!, serviceId),
    enabled: !!projectId,
  });
};

export const useSaveCode = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      projectId: string;
      serviceId: string;
      payload: { codeContent: string; runtime?: string };
    }) => serverlessService.saveCode(payload.projectId, payload.serviceId, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['serverless'] });
    },
  });
};
