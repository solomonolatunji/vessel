import { useQuery } from '@tanstack/react-query';
import { canvasService } from '#/services/canvas';

export const useListCanvasSummaries = () => {
  return useQuery({
    queryKey: ['canvas', 'listCanvasSummaries'].filter(Boolean),
    queryFn: () => canvasService.listCanvasSummaries(),
  });
};

export const useGetCanvasSummary = (projectId: string) => {
  return useQuery({
    queryKey: ['canvas', 'getCanvasSummary', projectId].filter(Boolean),
    queryFn: () => canvasService.getCanvasSummary(projectId),
  });
};

export const useGetEnvironmentCanvas = (envId: string) => {
  return useQuery({
    queryKey: ['canvas', 'getEnvironmentCanvas', envId].filter(Boolean),
    queryFn: () => canvasService.getEnvironmentCanvas(envId),
  });
};
