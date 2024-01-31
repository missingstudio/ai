package openai

import (
	"net/http"
	"strings"

	"github.com/missingstudio/studio/backend/internal/providers/base"
	"github.com/missingstudio/studio/backend/pkg/utils"
	"github.com/missingstudio/studio/common/errors"
)

type OpenAIProviderFactory struct{}

type OpenAIHeaders struct {
	APIKey string `validate:"required" json:"Authorization" error:"API key is required"`
}

func (azf OpenAIProviderFactory) Validate(headers http.Header) (*OpenAIHeaders, error) {
	var oaiHeaders OpenAIHeaders
	if err := utils.UnmarshalHeader(headers, &oaiHeaders); err != nil {
		return nil, errors.New(err)
	}

	if err := utils.ValidateHeaders(oaiHeaders); err != nil {
		return nil, err
	}

	return &oaiHeaders, nil
}

func (oaif OpenAIProviderFactory) Create(headers http.Header) (base.ProviderInterface, error) {
	oaiHeaders, err := oaif.Validate(headers)
	if err != nil {
		return nil, err
	}

	oaiHeaders.APIKey = strings.Replace(oaiHeaders.APIKey, "Bearer ", "", 1)
	openAIProvider := NewOpenAIProvider(*oaiHeaders)
	return openAIProvider, nil
}

type OpenAIProvider struct {
	Name   string
	Config base.ProviderConfig
	OpenAIHeaders
}

func NewOpenAIProvider(headers OpenAIHeaders) *OpenAIProvider {
	config := getOpenAIConfig("https://api.openai.com")

	return &OpenAIProvider{
		Name:          "Open AI",
		Config:        config,
		OpenAIHeaders: headers,
	}
}

func (oai OpenAIProvider) GetName() string {
	return oai.Name
}

func getOpenAIConfig(baseURL string) base.ProviderConfig {
	return base.ProviderConfig{
		BaseURL:         baseURL,
		ChatCompletions: "/v1/chat/completions",
	}
}
