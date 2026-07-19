import type { BaseResponse } from './base';

export type OAuthProviderConfig = {
  id: string;
  providerName: string;
  clientId: string;
  clientSecret?: string;
  baseUrl?: string;
  redirectUri?: string;
  tenant?: string;
  enabled: boolean;
  createdAt: string;
  updatedAt: string;
};

export type ListOAuthProvidersResponse = BaseResponse<OAuthProviderConfig[]>;
export type ListEnabledProvidersResponse = BaseResponse<OAuthProviderConfig[]>;

export type SaveOAuthProviderRequest = {
  id?: string;
  providerName: string;
  clientId: string;
  clientSecret?: string;
  baseUrl?: string;
  redirectUri?: string;
  tenant?: string;
  enabled: boolean;
};

export type SaveOAuthProviderResponse = BaseResponse<OAuthProviderConfig>;
