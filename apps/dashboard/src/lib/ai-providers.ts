export type AiProviderModel = {
  id: string;
  name: string;
};

export type AiProviderCatalogItem = {
  id: AiProviderId;
  name: string;
  apiKeyPlaceholder: string;
  models: AiProviderModel[];
};

export const aiProviderIdValues = [
  'openai',
  'anthropic',
  'google',
  'mistral',
  'groq',
  'deepseek',
  'xai',
  'moonshot',
] as const;
export type AiProviderId = (typeof aiProviderIdValues)[number];

export const aiProviderCatalog = [
  {
    id: 'openai',
    name: 'OpenAI',
    apiKeyPlaceholder: 'sk-...',
    models: [
      { id: 'gpt-5.2', name: 'GPT-5.2' },
      { id: 'gpt-5.2-pro', name: 'GPT-5.2 Pro' },
      { id: 'gpt-5.1', name: 'GPT-5.1' },
      { id: 'gpt-5-mini', name: 'GPT-5 mini' },
      { id: 'gpt-4.1', name: 'GPT-4.1' },
    ],
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    apiKeyPlaceholder: 'sk-ant-...',
    models: [
      { id: 'claude-opus-4-5-20251101', name: 'Claude Opus 4.5' },
      { id: 'claude-sonnet-4-5-20250929', name: 'Claude Sonnet 4.5' },
      { id: 'claude-haiku-4-5-20251001', name: 'Claude Haiku 4.5' },
    ],
  },
  {
    id: 'google',
    name: 'Google',
    apiKeyPlaceholder: 'Google AI API key',
    models: [
      { id: 'gemini-3-pro-preview', name: 'Gemini 3 Pro Preview' },
      { id: 'gemini-2.5-pro', name: 'Gemini 2.5 Pro' },
      { id: 'gemini-2.5-flash', name: 'Gemini 2.5 Flash' },
      { id: 'gemini-2.5-flash-lite', name: 'Gemini 2.5 Flash Lite' },
    ],
  },
  {
    id: 'mistral',
    name: 'Mistral',
    apiKeyPlaceholder: 'Mistral API key',
    models: [
      { id: 'mistral-medium-2508', name: 'Mistral Medium 3.1' },
      { id: 'magistral-medium-2509', name: 'Magistral Medium 1.2' },
      { id: 'codestral-2508', name: 'Codestral' },
      { id: 'mistral-large-2411', name: 'Mistral Large 2.1' },
      { id: 'ministral-8b-2410', name: 'Ministral 8B' },
    ],
  },
  {
    id: 'groq',
    name: 'Groq',
    apiKeyPlaceholder: 'gsk_...',
    models: [
      { id: 'llama-3.3-70b-versatile', name: 'Llama 3.3 70B Versatile' },
      { id: 'moonshotai/kimi-k2-instruct-0905', name: 'Kimi K2 Instruct' },
      { id: 'deepseek-r1-distill-llama-70b', name: 'DeepSeek R1 Distill Llama 70B' },
      { id: 'compound-beta', name: 'Compound Beta' },
    ],
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    apiKeyPlaceholder: 'DeepSeek API key',
    models: [
      { id: 'deepseek-chat', name: 'DeepSeek Chat' },
      { id: 'deepseek-reasoner', name: 'DeepSeek Reasoner' },
    ],
  },
  {
    id: 'xai',
    name: 'xAI',
    apiKeyPlaceholder: 'xAI API key',
    models: [
      { id: 'grok-4.3', name: 'Grok 4.3' },
      { id: 'grok-4.3-latest', name: 'Grok 4.3 Latest' },
      { id: 'grok-4-fast-reasoning', name: 'Grok 4 Fast Reasoning' },
      { id: 'grok-4-fast-non-reasoning', name: 'Grok 4 Fast Non-reasoning' },
      { id: 'grok-4-1-fast-reasoning', name: 'Grok 4.1 Fast Reasoning' },
      { id: 'grok-4-1-fast-non-reasoning', name: 'Grok 4.1 Fast Non-reasoning' },
      { id: 'grok-build-0.1', name: 'Grok Build 0.1' },
    ],
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    apiKeyPlaceholder: 'Moonshot API key',
    models: [
      { id: 'kimi-k2.6', name: 'Kimi K2.6' },
      { id: 'kimi-for-coding', name: 'Kimi for Coding' },
      { id: 'moonshot-v1-8k-vision-preview', name: 'Moonshot Vision Preview' },
    ],
  },
] as const satisfies readonly AiProviderCatalogItem[];

export const aiProviderIds = [...aiProviderIdValues];

export function isAiProviderId(value: unknown): value is AiProviderId {
  return typeof value === 'string' && aiProviderIds.includes(value as AiProviderId);
}

export function aiProviderName(providerId: AiProviderId) {
  return aiProviderCatalog.find((provider) => provider.id === providerId)?.name ?? providerId;
}

export function aiProviderModels(providerId: AiProviderId): AiProviderModel[] {
  return [...(aiProviderCatalog.find((provider) => provider.id === providerId)?.models ?? [])];
}

export function defaultAiModel(providerId: AiProviderId) {
  return aiProviderModels(providerId)[0]?.id ?? '';
}

export function isAiProviderModel(providerId: AiProviderId, modelId: unknown): modelId is string {
  return (
    typeof modelId === 'string' &&
    aiProviderModels(providerId).some((model) => model.id === modelId)
  );
}
