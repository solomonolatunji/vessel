import type { BaseResponse } from './base';

export interface GitStatus {
  provider: string;
  connected: boolean;
}

export interface GitRepo {
  id: string;
  name: string;
  fullName: string;
  cloneUrl: string;
  private: boolean;
  defaultBranch: string;
  updatedAt: string;
}

export interface GitBranch {
  name: string;
  sha: string;
}

export interface ConnectGitRequest {
  provider: string;
  code: string;
}

export type GetGitStatusResponse = BaseResponse<GitStatus[]>;
export type ListGitReposResponse = BaseResponse<GitRepo[]>;
export type ListGitBranchesResponse = BaseResponse<GitBranch[]>;
export type ConnectGitResponse = BaseResponse<GitStatus>;
