import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { BaseResponse, PaginatedData } from '#/interfaces/base';
import type { User } from '#/interfaces/users';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

const usersService = {
  list: async (): Promise<BaseResponse<PaginatedData<User>>> => {
    try {
      return await apiClient.get<BaseResponse<PaginatedData<User>>>('/users');
    } catch (err) {
      throw handleApiError(err);
    }
  },
  delete: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/users/${id}`);
    } catch (err) {
      throw handleApiError(err);
    }
  },
  invite: async (payload: { email: string; role: string }): Promise<User> => {
    try {
      const res = await apiClient.post<BaseResponse<User>>('/users/invite', payload);
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  },
};

export const useListUsers = () =>
  useQuery({ queryKey: ['users'], queryFn: () => usersService.list() });

export const useDeleteUser = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => usersService.delete(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['users'] }),
  });
};

export const useInviteUser = () => {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (payload: { email: string; role: string }) => usersService.invite(payload),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['users'] }),
  });
};
