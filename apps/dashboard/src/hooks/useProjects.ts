import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { projectsService } from '#/services/projects';

export const useListProjects = () => {
  return useQuery({
    queryKey: ['projects', 'listProjects'],
    queryFn: () => projectsService.listProjects(),
  });
};

export const useGetProject = (id: string) => {
  return useQuery({
    queryKey: ['projects', 'getProject', id].filter(Boolean),
    queryFn: () => projectsService.getProject(id),
  });
};

export const useCreateProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof projectsService.createProject>[0] }) =>
      projectsService.createProject(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useDeleteProject = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => projectsService.deleteProject(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};

export const useGetVars = (id: string) => {
  return useQuery({
    queryKey: ['projects', 'getVars', id].filter(Boolean),
    queryFn: () => projectsService.getVars(id),
  });
};

export const useSetVars = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; payload: { variables: Record<string, string> } }) =>
      projectsService.setVars(payload.id, payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['projects'] });
      await queryClient.invalidateQueries({ queryKey: ['canvas'] });
    },
  });
};
