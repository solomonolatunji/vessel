import type { BaseResponse } from './base';

export interface TeamAISettings {
  id: string;
  provider: string;
  apiKey?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TeamEmailSettings {
  id: string;
  smtpHost?: string;
  smtpPort?: number;
  smtpUser?: string;
  smtpPassword?: string;
  smtpFromName?: string;
  smtpFromAddress?: string;
  resendApiKey?: string;
  useResend: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface ServerSettings {
  id: string;
  traefikWildcardIp: string;
  discordWebhookUrl?: string;
  discordPingEnabled: boolean;
  discordEnabled: boolean;
  slackWebhookUrl?: string;
  slackEnabled: boolean;
  telegramBotToken?: string;
  telegramChatId?: string;
  telegramEnabled: boolean;
  smtpHost?: string;
  smtpPort?: number;
  smtpUser?: string;
  smtpPassword?: string;
  smtpFromName?: string;
  smtpFromAddress?: string;
  smtpEnabled: boolean;
  resendApiKey?: string;
  resendEnabled: boolean;
  pushoverUserKey?: string;
  pushoverApiToken?: string;
  pushoverEnabled: boolean;
  genericWebhookUrl?: string;
  genericWebhookEnabled: boolean;
  notificationAlerts: boolean;
  registrationEnabled: boolean;
  registrationDomainAllowlist?: string;
  customDnsResolvers: string;
  dnsValidationEnabled: boolean;
  ipAllowlist: string;
  mcpServerEnabled: boolean;
  defaultWildcardDomain?: string;
  dashboardDomain?: string;
  siteName?: string;
  publicIpv4?: string;
  publicIpv6?: string;
  showSponsorshipPopup: boolean;
  disableTwoStepConfirmation: boolean;
  defaultOpenAIKey?: string;
  defaultAnthropicKey?: string;
  updateCheckCron: string;
  autoUpdateEnabled: boolean;
  telemetryEnabled: boolean;
  concurrentBuilds: number;
  deploymentTimeout: number;
  serverTimezone: string;
  dockerCleanupCron: string;
  diskUsageThreshold: number;
  diskUsageCron: string;
  currentVersion: string;
  latestVersion: string;
  lastUpdateCheck: string;
  cloudflareApiToken?: string;
  namecheapApiUser?: string;
  namecheapApiKey?: string;
  namecheapClientIp?: string;
  spaceshipApiKey?: string;
  updatedAt: string;
}

export type UpdateSettingsRequest = Partial<ServerSettings>;

export interface PruneResponse {
  status: string;
  message: string;
  spaceReclaimedBytes: number;
}

export interface MCPResponse {
  jsonrpc: string;
  id?: unknown;
  result?: unknown;
  error?: MCPError;
  server?: Record<string, unknown>;
  tools?: Record<string, unknown>[];
  capabilities?: Record<string, unknown>;
}

export interface MCPError {
  code: number;
  message: string;
}

export interface TestNotificationRequest {
  provider: string;
}

export interface GithubApp {
  id: string;
  name: string;
  appId: string;
  installationId: string;
  clientId: string;
  clientSecret: string;
  webhookSecret: string;
  privateKey: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface GitlabApp {
  id: string;
  name: string;
  appId: string;
  appSecret: string;
  webhookSecret: string;
  apiUrl: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface BitbucketApp {
  id: string;
  name: string;
  owner: string;
  clientId: string;
  clientSecret: string;
  webhookSecret: string;
  isPublic: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UpdateInfo {
  currentVersion: string;
  latestVersion: string;
  hasUpdate: boolean;
  releaseNotes: string;
  downloadUrl: string;
  lastChecked: string;
  autoUpdate: boolean;
  updateCheckCron: string;
}

export interface GitAppsManifestRequest {
  code: string;
}

export type GetServerSettingsResponse = BaseResponse<ServerSettings>;
export type UpdateServerSettingsResponse = BaseResponse<ServerSettings>;
export type TestNotificationResponseType = BaseResponse<void>;
export type GetGithubAppsResponse = BaseResponse<GithubApp[]>;
export type SaveGithubAppResponse = BaseResponse<GithubApp>;
export type GetGitlabAppsResponse = BaseResponse<GitlabApp[]>;
export type SaveGitlabAppResponse = BaseResponse<GitlabApp>;
export type GetBitbucketAppsResponse = BaseResponse<BitbucketApp[]>;
export type SaveBitbucketAppResponse = BaseResponse<BitbucketApp>;
export type ExchangeGithubManifestResponse = BaseResponse<GithubApp>;
export type GetUpdateStatusResponse = BaseResponse<UpdateInfo>;
export type CheckUpdateResponse = BaseResponse<UpdateInfo>;
export type DeployUpdateResponse = BaseResponse<void>;
export type SystemPruneResponse = BaseResponse<PruneResponse>;
