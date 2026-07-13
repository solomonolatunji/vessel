export interface TeamAISettings {
  id: string;
  workspaceId: string;
  provider: string;
  apiKey?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TeamEmailSettings {
  id: string;
  workspaceId: string;
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
  licenseKey?: string;
  plan: string;
  maxSeats: number;
  currentVersion: string;
  latestVersion: string;
  lastUpdateCheck: string;
  updatedAt: string;
}

export type UpdateSettingsRequest = Record<string, unknown>;

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
