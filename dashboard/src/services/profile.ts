import type { BaseResponse } from '#/interfaces/base';
import type {
  ChangePasswordRequest,
  GetProfileResponse,
  RequestEmailChangeRequest,
  UpdateProfileRequest,
  VerifyEmailChangeRequest,
} from '#/interfaces/profile';
import type { CreatePATRequest, CreatePATResponse, PersonalAccessToken } from '#/interfaces/users';
import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

type ListPATsResponse = BaseResponse<PersonalAccessToken[]>;

export const profileService = {
  getProfile: async (): Promise<GetProfileResponse> => {
    try {
      const response = await apiClient.get<GetProfileResponse>('/profile');
      return response;
    } catch (error) {
      throw handleApiError(error);
    }
  },

  updateProfile: async (payload: UpdateProfileRequest): Promise<void> => {
    try {
      await apiClient.put('/profile', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  requestEmailChange: async (payload: RequestEmailChangeRequest): Promise<void> => {
    try {
      await apiClient.post('/profile/email/request', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  verifyEmailChange: async (payload: VerifyEmailChangeRequest): Promise<void> => {
    try {
      await apiClient.post('/profile/email/verify', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  changePassword: async (payload: ChangePasswordRequest): Promise<void> => {
    try {
      await apiClient.put('/profile/password', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  listTokens: async (): Promise<ListPATsResponse> => {
    try {
      return await apiClient.get<ListPATsResponse>('/profile/tokens');
    } catch (error) {
      throw handleApiError(error);
    }
  },

  createToken: async (payload: CreatePATRequest): Promise<CreatePATResponse> => {
    try {
      return await apiClient.post<CreatePATResponse>('/profile/tokens', payload);
    } catch (error) {
      throw handleApiError(error);
    }
  },

  deleteToken: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/profile/tokens/${id}`);
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
