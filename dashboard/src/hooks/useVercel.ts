import { useQuery } from '@tanstack/react-query';
import { vercelService } from '#/services/vercel';

export const useListProjects = (enabled = true) => {
  return useQuery({
    queryKey: ['vercel', 'listProjects'].filter(Boolean),
    queryFn: () => vercelService.listProjects(),
    enabled,
  });
};

export const useGetProjectEnv = (id: string, enabled = true) => {
  return useQuery({
    queryKey: ['vercel', 'getProjectEnv', id].filter(Boolean),
    queryFn: () => vercelService.getProjectEnv(id),
    enabled: enabled && !!id,
  });
};
