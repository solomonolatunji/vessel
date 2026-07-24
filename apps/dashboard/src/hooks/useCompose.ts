import { useMutation, useQueryClient } from '@tanstack/react-query';
import { type ComposeAnalyzeRequest, composeService } from '#/services/compose';

export const useAnalyzeCompose = () => {
  return useMutation({
    mutationFn: (req: ComposeAnalyzeRequest) => composeService.analyze(req),
  });
};

export const useDeployCompose = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ projectId, composeContent }: { projectId: string; composeContent: string }) =>
      composeService.deploy(projectId, composeContent),
    onSuccess: async (_, { projectId }) => {
      // Invalidate projects to fetch updated services
      await queryClient.invalidateQueries({ queryKey: ['projects', 'getProject', projectId] });
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
    },
  });
};
