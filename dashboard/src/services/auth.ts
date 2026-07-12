import type { AuthCredentials, AuthResponse, RegisterCredentials } from '#/interfaces/auth';
import type { User } from '#/interfaces/users';
import { apiClient } from './instance';

export const authService = {
  login: async (credentials: AuthCredentials): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/signin', credentials);
    return data;
  },

  register: async (details: RegisterCredentials): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/signup', details);
    return data;
  },

  logout: async (): Promise<void> => {
    await apiClient.post('/auth/logout');
  },
};
