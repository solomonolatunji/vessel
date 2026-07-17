import { Brain, ChevronDown, Star } from 'lucide-react';
import { useState } from 'react';
import { Input } from '#/components/ui/input';
import { useGetAISettings, useUpdateAISettings } from '#/hooks/useSettings';

const PROVIDERS = [
  {
    id: 'openai',
    name: 'OpenAI',
    model: 'GPT-5.2',
    icon: '/ai-providers/openai.svg',
    keyField: 'openAIKey' as const,
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    model: 'Claude Opus 4.5',
    icon: '/ai-providers/anthropic.svg',
    keyField: 'anthropicKey' as const,
  },
  {
    id: 'google',
    name: 'Google',
    model: 'Gemini 3 Pro Preview',
    icon: '/ai-providers/google.svg',
    keyField: 'googleKey' as const,
  },
  {
    id: 'mistral',
    name: 'Mistral',
    model: 'Mistral Medium 3.1',
    icon: '/ai-providers/mistral.svg',
    keyField: 'mistralKey' as const,
  },
  {
    id: 'groq',
    name: 'Groq',
    model: 'Llama 3.3 70B Ver...',
    icon: '/ai-providers/groq.svg',
    keyField: 'groqKey' as const,
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    model: 'DeepSeek Chat',
    icon: '/ai-providers/deepseek.svg',
    keyField: 'deepSeekKey' as const,
  },
  {
    id: 'xai',
    name: 'xAI',
    model: 'Grok 4.3',
    icon: '/ai-providers/xai.svg',
    keyField: 'xaiKey' as const,
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    model: 'Kimi K2.6',
    icon: '/ai-providers/moonshot.svg',
    keyField: 'moonshotKey' as const,
  },
] as const;

export function AISettings() {
  const { data: settings } = useGetAISettings();
  const updateSettings = useUpdateAISettings();
  const [editingId, setEditingId] = useState<string | null>(null);

  const defaultProvider = settings?.data?.defaultProvider || 'none';

  const handleSetDefault = (id: string) => {
    updateSettings.mutate({ defaultProvider: id });
  };

  const handleUpdateKey = (keyField: string, value: string) => {
    updateSettings.mutate({ [keyField]: value });
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col justify-between gap-6 pb-2 md:flex-row md:items-start">
        <div className="flex items-center gap-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-xl border border-border/50 bg-card/40">
            <Brain className="h-5 w-5 text-primary" />
          </div>
          <h1 className="font-bold text-3xl tracking-tight">AI</h1>
        </div>
        <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
          DEFAULT <span className="text-foreground">{defaultProvider}</span>
        </p>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        {PROVIDERS.map((provider) => {
          const isEditing = editingId === provider.id;
          const currentKey = settings?.data?.[provider.keyField] || '';
          const isSet = currentKey !== '';
          const isDefault = defaultProvider === provider.id;

          return (
            <button
              type="button"
              key={provider.id}
              className="relative flex flex-col justify-between space-y-4 rounded-xl border border-border/50 bg-card/40 p-5 text-left transition-colors hover:border-border"
              onClick={() => {
                if (!isEditing) setEditingId(provider.id);
              }}
            >
              <div className="flex w-full items-start justify-between">
                <div className="flex items-start gap-4">
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-lg bg-black">
                    <img
                      src={provider.icon}
                      alt={provider.name}
                      className="h-6 w-6 object-contain"
                    />
                  </div>
                  <div className="space-y-1">
                    <h3 className="font-semibold text-sm leading-none">{provider.name}</h3>
                    <div className="flex items-center text-muted-foreground text-xs">
                      {provider.model}
                      <ChevronDown className="ml-1 h-3 w-3" />
                    </div>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleSetDefault(provider.id);
                  }}
                  className="flex h-8 w-8 items-center justify-center rounded-md border border-border/50 bg-background/50 text-muted-foreground hover:text-foreground"
                >
                  <Star
                    className={`h-4 w-4 ${isDefault ? 'fill-foreground text-foreground' : ''}`}
                  />
                </button>
              </div>

              <div className="mt-4 w-full">
                {isEditing || isSet ? (
                  <Input
                    autoFocus={isEditing && !isSet}
                    onClick={(e) => e.stopPropagation()}
                    type="password"
                    placeholder="sk-..."
                    defaultValue={currentKey}
                    className="h-8 bg-background/50 font-mono text-xs"
                    onBlur={(e) => {
                      if (e.target.value !== currentKey) {
                        handleUpdateKey(provider.keyField, e.target.value);
                      }
                      if (!e.target.value) {
                        setEditingId(null);
                      }
                    }}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.currentTarget.blur();
                      }
                    }}
                  />
                ) : (
                  <div className="flex h-8 items-center text-muted-foreground text-xs">
                    (API key unset)
                  </div>
                )}
              </div>
            </button>
          );
        })}
      </div>
    </div>
  );
}
