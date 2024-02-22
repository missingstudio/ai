package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/missingstudio/studio/backend/internal/constants"
	"github.com/missingstudio/studio/backend/internal/httputil"
	"github.com/missingstudio/studio/backend/internal/providers"
	"github.com/missingstudio/studio/backend/internal/providers/base"
	"github.com/missingstudio/studio/backend/models"
	"github.com/missingstudio/studio/backend/pkg/utils"
	"github.com/missingstudio/studio/common/errors"
	llmv1 "github.com/missingstudio/studio/protos/pkg/llm"
)

var (
	ErrChatCompletionStreamNotSupported = errors.NewBadRequest("streaming is not supported with this method, please use StreamChatCompletions")
	ErrChatCompletionNotSupported       = errors.NewInternalError("provider don't have chat Completion capabilities")
	ErrRequiredHeaderNotExit            = errors.NewBadRequest(fmt.Sprintf("either %s or %s header is required", constants.XMSProvider, constants.XMSConfig))
)

func (s *V1Handler) ChatCompletions(
	ctx context.Context,
	req *connect.Request[llmv1.ChatCompletionRequest],
) (*connect.Response[llmv1.ChatCompletionResponse], error) {
	// Check if required headers are available
	providerName := req.Header().Get(constants.XMSProvider)
	config := req.Header().Get(constants.XMSConfig)
	if providerName == "" && config == "" {
		return nil, ErrRequiredHeaderNotExit
	}

	startTime := time.Now()
	payload, err := json.Marshal(req.Msg)
	if err != nil {
		return nil, errors.New(err)
	}

	headerConfig := httputil.GetContextWithHeaderConfig(ctx)
	connectionObj := models.Connection{
		Name:    providerName,
		Headers: headerConfig,
	}

	provider, err := s.providerService.GetProvider(connectionObj)
	if err != nil {
		return nil, errors.New(err)
	}

	// Validate provider configs
	err = providers.Validate(provider, map[string]any{
		"headers": headerConfig,
	})
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	chatCompletionProvider, ok := provider.(base.ChatCompletionInterface)
	if !ok {
		return nil, ErrChatCompletionNotSupported
	}

	resp, err := chatCompletionProvider.ChatCompletion(ctx, payload)
	if err != nil {
		return nil, errors.New(err)
	}

	latency := time.Since(startTime)
	data := &llmv1.ChatCompletionResponse{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, errors.New(err)
	}

	ingesterdata := make(map[string]any)
	providerInfo := provider.Info()

	ingesterdata["provider"] = providerInfo.Name
	ingesterdata["model"] = data.Model
	ingesterdata["latency"] = latency
	ingesterdata["total_tokens"] = *data.Usage.TotalTokens
	ingesterdata["prompt_tokens"] = *data.Usage.PromptTokens
	ingesterdata["completion_tokens"] = *data.Usage.CompletionTokens

	go s.ingester.Ingest(ingesterdata, "analytics")

	response := connect.NewResponse(data)
	return utils.CopyHeaders[llmv1.ChatCompletionResponse](resp, response), nil
}
