import { Brain, Check, ChevronDown, Star } from 'lucide-react';
import React, { useState } from 'react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '#/components/ui/dropdown-menu';
import { Input } from '#/components/ui/input';
import { useGetAISettings, useUpdateAISettings } from '#/hooks/useSettings';
import { aiProviderCatalog } from '#/lib/ai-providers';

const PROVIDERS = [
  {
    id: 'openai',
    name: 'OpenAI',
    models: aiProviderCatalog.find((c) => c.id === 'openai')?.models || [],
    icon: '/ai-providers/openai.svg',
    keyField: 'openAIKey' as const,
    modelField: 'openAIModel' as const,
  },
  {
    id: 'anthropic',
    name: 'Anthropic',
    models: aiProviderCatalog.find((c) => c.id === 'anthropic')?.models || [],
    icon: '/ai-providers/anthropic.svg',
    keyField: 'anthropicKey' as const,
    modelField: 'anthropicModel' as const,
  },
  {
    id: 'google',
    name: 'Google',
    models: aiProviderCatalog.find((c) => c.id === 'google')?.models || [],
    icon: '/ai-providers/google.svg',
    keyField: 'googleKey' as const,
    modelField: 'googleModel' as const,
  },
  {
    id: 'mistral',
    name: 'Mistral',
    models: aiProviderCatalog.find((c) => c.id === 'mistral')?.models || [],
    icon: '/ai-providers/mistral.svg',
    keyField: 'mistralKey' as const,
    modelField: 'mistralModel' as const,
  },
  {
    id: 'groq',
    name: 'Groq',
    models: aiProviderCatalog.find((c) => c.id === 'groq')?.models || [],
    icon: '/ai-providers/groq.svg',
    keyField: 'groqKey' as const,
    modelField: 'groqModel' as const,
  },
  {
    id: 'deepseek',
    name: 'DeepSeek',
    models: aiProviderCatalog.find((c) => c.id === 'deepseek')?.models || [],
    icon: '/ai-providers/deepseek.svg',
    keyField: 'deepSeekKey' as const,
    modelField: 'deepSeekModel' as const,
  },
  {
    id: 'xai',
    name: 'xAI',
    models: aiProviderCatalog.find((c) => c.id === 'xai')?.models || [],
    icon: '/ai-providers/xai.svg',
    keyField: 'xaiKey' as const,
    modelField: 'xaiModel' as const,
  },
  {
    id: 'moonshot',
    name: 'Moonshot',
    models: aiProviderCatalog.find((c) => c.id === 'moonshot')?.models || [],
    icon: '/ai-providers/moonshot.svg',
    keyField: 'moonshotKey' as const,
    modelField: 'moonshotModel' as const,
  },
] as const;

export function AISettings() {
  const { data: settings } = useGetAISettings();
  const updateSettings = useUpdateAISettings();
  const [editingId, setEditingId] = useState<string | null>(null);

  const pendingSettings = React.useRef<Record<string, unknown> | null>(null);

  React.useEffect(() => {
    if (settings?.data && !updateSettings.isPending) {
      pendingSettings.current = settings.data;
    }
  }, [settings?.data, updateSettings.isPending]);

  const defaultProvider = (settings?.data?.defaultProvider as string) || 'none';

  const handleSetDefault = (id: string) => {
    const payload = { ...(pendingSettings.current || settings?.data || {}), defaultProvider: id };
    pendingSettings.current = payload;
    updateSettings.mutate(payload);
  };

  const handleUpdateKey = (keyField: string, value: string) => {
    const payload = { ...(pendingSettings.current || settings?.data || {}), [keyField]: value };
    pendingSettings.current = payload;
    updateSettings.mutate(payload);
  };

  return (
    <div className="space-y-6">
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <Brain className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">AI</h1>
            <p className="text-muted-foreground text-sm">
              Configure built-in AI models and providers for your Codedock instance.
            </p>
          </div>
        </div>
        <div className="flex shrink-0 items-center gap-4">
          <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-widest">
            DEFAULT <span className="text-foreground">{defaultProvider}</span>
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
        {PROVIDERS.map((provider) => {
          const isDefault = defaultProvider === provider.id;
          const currentKey = (settings?.data?.[provider.keyField] as string) || '';
          const currentModel =
            (settings?.data?.[provider.modelField] as string) || provider.models[0]?.id || '';
          const isSet = currentKey.length > 0;
          const isEditing = editingId === provider.id;

          return (
            <button
              type="button"
              key={provider.id}
              className="relative flex flex-col justify-between space-y-4 rounded-xl border border-border/50 bg-card/40 p-6 text-left transition-colors hover:border-border"
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
                    <div className="flex items-center text-muted-foreground">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <button
                            type="button"
                            className="flex max-w-37.5 cursor-pointer appearance-none items-center gap-1 truncate bg-transparent font-medium text-xs outline-none hover:text-foreground/80"
                            onClick={(e) => e.stopPropagation()}
                          >
                            <span className="truncate">
                              {provider.models.find((m) => m.id === currentModel)?.name ||
                                currentModel}
                            </span>
                            <ChevronDown className="h-3 w-3 shrink-0" />
                          </button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="start" className="w-50">
                          {provider.models.map((m) => (
                            <DropdownMenuItem
                              key={m.id}
                              onSelect={() => handleUpdateKey(provider.modelField, m.id)}
                              className="flex cursor-pointer items-center justify-between py-2"
                            >
                              <div className="flex flex-col">
                                <span className="font-medium text-sm">{m.name}</span>
                                <span className="font-mono text-[10px] text-muted-foreground">
                                  {m.id}
                                </span>
                              </div>
                              {currentModel === m.id && (
                                <Check className="h-4 w-4 text-emerald-500" />
                              )}
                            </DropdownMenuItem>
                          ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleSetDefault(isDefault ? 'none' : provider.id);
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
                    className="h-10 font-mono text-sm"
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
