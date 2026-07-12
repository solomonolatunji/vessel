import type { User } from '#/interfaces/users';
import { apiClient } from './instance';

export const profileService = {
  getProfile: async (): Promise<User> => {
    const { data } = await apiClient.get<User>('/profile');
    return data;
  },

  updateProfile: async (payload: { email?: string; role?: string }): Promise<User> => {
    const { data } = await apiClient.put<User>('/profile', payload);
    return data;
  },
};
