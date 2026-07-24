import type { BaseResponse } from './base';
import type { UserRole } from './users';

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  avatarUrl?: string;
  totpEnabled: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateProfileRequest {
  name?: string;
  email?: string;
  avatarUrl?: string;
}

export interface ChangePasswordRequest {
  oldPassword?: string;
  newPassword?: string;
}

export interface RequestEmailChangeRequest {
  newEmail: string;
}

export interface VerifyEmailChangeRequest {
  otp: string;
}

export type GetProfileResponse = BaseResponse<UserProfile>;
export type UpdateProfileResponse = BaseResponse<UserProfile>;
