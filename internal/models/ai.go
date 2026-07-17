package models

type AISettings struct {
	ID              string `json:"id"`
	DefaultProvider string `json:"defaultProvider"`
	OpenAIKey       string `json:"openAIKey,omitempty"`
	AnthropicKey    string `json:"anthropicKey,omitempty"`
	GoogleKey       string `json:"googleKey,omitempty"`
	MistralKey      string `json:"mistralKey,omitempty"`
	GroqKey         string `json:"groqKey,omitempty"`
	DeepSeekKey     string `json:"deepSeekKey,omitempty"`
	XAIKey          string `json:"xaiKey,omitempty"`
	MoonshotKey     string `json:"moonshotKey,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type UpdateAISettingsRequest struct {
	AISettings
}
