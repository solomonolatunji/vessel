import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { profileService } from '#/services/profile';

export const useGetProfile = () => {
  return useQuery({
    queryKey: ['profile', 'getProfile'].filter(Boolean),
    queryFn: () => profileService.getProfile(),
  });
};

export const useUpdateProfile = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof profileService.updateProfile>[0]) =>
      profileService.updateProfile(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};

export const useRequestEmailChange = () => {
  return useMutation({
    mutationFn: (payload: Parameters<typeof profileService.requestEmailChange>[0]) =>
      profileService.requestEmailChange(payload),
  });
};

export const useVerifyEmailChange = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Parameters<typeof profileService.verifyEmailChange>[0]) =>
      profileService.verifyEmailChange(payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};

export const useChangePassword = () => {
  return useMutation({
    mutationFn: (payload: Parameters<typeof profileService.changePassword>[0]) =>
      profileService.changePassword(payload),
  });
};

export const useListTokens = () => {
  return useQuery({
    queryKey: ['profile', 'listTokens'].filter(Boolean),
    queryFn: () => profileService.listTokens(),
  });
};

export const useCreateToken = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { payload: Parameters<typeof profileService.createToken>[0] }) =>
      profileService.createToken(payload.payload),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};

export const useDeleteToken = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string }) => profileService.deleteToken(payload.id),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['profile'] });
    },
  });
};
