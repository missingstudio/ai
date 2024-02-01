package togetherai

import (
	"net/http"
	"strings"

	"github.com/missingstudio/studio/backend/internal/providers/base"
	"github.com/missingstudio/studio/backend/pkg/utils"
	"github.com/missingstudio/studio/common/errors"
)

type TogetherAIProviderFactory struct{}

type TogetherAIHeaders struct {
	APIKey string `validate:"required" json:"Authorization" error:"API key is required"`
}

func (ta TogetherAIProviderFactory) GetHeaders(headers http.Header) (*TogetherAIHeaders, error) {
	var togetherAIHeaders TogetherAIHeaders
	if err := utils.UnmarshalHeader(headers, &togetherAIHeaders); err != nil {
		return nil, errors.New(err)
	}

	return &togetherAIHeaders, nil
}

func (ta TogetherAIProviderFactory) Create(headers http.Header) (base.ProviderInterface, error) {
	togetherAIHeaders, err := ta.GetHeaders(headers)
	if err != nil {
		return nil, err
	}

	togetherAIHeaders.APIKey = strings.Replace(togetherAIHeaders.APIKey, "Bearer ", "", 1)
	openAIProvider := NewTogetherAIProvider(*togetherAIHeaders)
	return openAIProvider, nil
}

type TogetherAIProvider struct {
	Name   string
	Config base.ProviderConfig
	TogetherAIHeaders
}

func NewTogetherAIProvider(headers TogetherAIHeaders) *TogetherAIProvider {
	config := getTogetherAIConfig("https://api.together.xyz")

	return &TogetherAIProvider{
		Name:              "TogetherAI",
		Config:            config,
		TogetherAIHeaders: headers,
	}
}

func (togetherAI TogetherAIProvider) GetName() string {
	return togetherAI.Name
}

func (togetherAI TogetherAIProvider) Validate() error {
	return utils.ValidateHeaders(togetherAI.TogetherAIHeaders)
}

func getTogetherAIConfig(baseURL string) base.ProviderConfig {
	return base.ProviderConfig{
		BaseURL:         baseURL,
		ChatCompletions: "/v1/chat/completions",
	}
}
