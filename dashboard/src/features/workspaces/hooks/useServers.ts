import { useQuery } from '@tanstack/react-query';
import { env } from '#/env';
import { cloudService } from '#/services/cloud';

export function useServers() {
  return useQuery({
    queryKey: ['cloud_servers'],
    queryFn: async () => {
      const profile = await cloudService.getProfile();
      // Flatten all servers from all teams
      const allServers = profile.teams.flatMap((team) => team.servers || []);
      return allServers;
    },
    enabled: env.VITE_IS_CLOUD,
  });
}
