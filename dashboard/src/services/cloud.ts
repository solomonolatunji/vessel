import { apiClient } from '#/lib/apiClient';
import { handleApiError } from '#/lib/error';

export interface CloudServer {
  id: string;
  name: string;
}

export interface CloudTeam {
  id: string;
  name: string;
  plan: string;
  custom_domain: string;
  logo_url: string;
  servers: CloudServer[];
}

export interface CloudProfileResponse {
  teams: CloudTeam[];
  // Other fields omitted for brevity
}

export const cloudService = {
  /**
   * Retrieves the current user's cloud profile, including teams and servers.
   */
  getProfile: async (): Promise<CloudProfileResponse> => {
    try {
      // Endpoint is bypassed by proxy interceptor because it starts with /users/
      return await apiClient.get<CloudProfileResponse>('/users/me');
    } catch (error) {
      throw handleApiError(error);
    }
  },
};
