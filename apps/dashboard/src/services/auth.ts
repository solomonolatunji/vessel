import type {
  AuthCredentials,
  AuthResponse,
  Disable2FAResponse,
  ForgotPasswordResponse,
  LoginResponse,
  RegisterCredentials,
  RegisterResponse,
  ResetPasswordResponse,
  Setup2FAResponseType,
  SetupCredentials,
  SetupResponse,
  Verify2FARequest,
  Verify2FAResponse,
} from '#/interfaces/auth';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export const authService = {
  login: async (credentials: AuthCredentials): Promise<AuthResponse> => {
    try {
      const response = await apiClient.post<LoginResponse>('/auth/signin', credentials);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  register: async (details: RegisterCredentials): Promise<AuthResponse> => {
    try {
      const response = await apiClient.post<RegisterResponse>('/auth/signup', details);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  setup: async (details: SetupCredentials): Promise<AuthResponse> => {
    try {
      const response = await apiClient.post<SetupResponse>('/system/setup', details);
      return response.data;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  logout: async (): Promise<void> => {
    try {
      await apiClient.post('/auth/logout');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  forgotPassword: async (email: string): Promise<ForgotPasswordResponse> => {
    try {
      return await apiClient.post<ForgotPasswordResponse>('/auth/forgot-password', { email });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  resetPassword: async (token: string, newPassword: string): Promise<ResetPasswordResponse> => {
    try {
      return await apiClient.post<ResetPasswordResponse>('/auth/reset-password', {
        token,
        newPassword,
      });
    } catch (error) {
      throw handleApiError(error);
    }
  },

  setup2FA: async (): Promise<Setup2FAResponseType> => {
    try {
      return await apiClient.post<Setup2FAResponseType>('/auth/2fa/setup');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  verify2FA: async (payload: Verify2FARequest): Promise<Verify2FAResponse> => {
    try {
      return await apiClient.post<Verify2FAResponse>('/auth/2fa/verify', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  disable2FA: async (payload: Verify2FARequest): Promise<Disable2FAResponse> => {
    try {
      return await apiClient.post<Disable2FAResponse>('/auth/2fa/disable', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
