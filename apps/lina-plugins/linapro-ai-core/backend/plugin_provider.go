// This file adapts the plugin-local text AI provider surface to the host
// aitext provider contract published through pluginhost Declarations. The
// plugin-local aitext package keeps its historical types, so the adapter maps
// the structurally identical request and response payloads across the two
// package boundaries.

package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	coreaitext "lina-core/pkg/plugin/capability/aicap/aitext"
	"lina-plugin-linapro-ai-core/backend/cap/aicap/aitext"
)

// provideCoreAIText creates the host-registered text AI provider from the
// plugin's framework provider factory. It satisfies the provider factory type
// required by pluginhost Declarations.Providers().ProvideAIText.
func provideCoreAIText(
	ctx context.Context,
	env coreaitext.ProviderEnv,
) (coreaitext.Provider, error) {
	provider, err := provideAIText(ctx, aitext.ProviderEnv{
		PluginID: env.PluginID,
		BizCtx:   env.BizCtx,
		Cache:    env.Cache,
	})
	if err != nil {
		return nil, err
	}
	return coreAITextProvider{inner: provider}, nil
}

// coreAITextProvider projects the plugin-local provider onto the host
// aitext.Provider contract.
type coreAITextProvider struct {
	inner aitext.Provider
}

// GenerateText executes one synchronous text generation request through the
// plugin-local provider while converting payloads across package boundaries.
func (p coreAITextProvider) GenerateText(
	ctx context.Context,
	request coreaitext.ProviderRequest,
) (*coreaitext.GenerateResponse, error) {
	if p.inner == nil {
		return nil, gerror.New("linapro-ai-core text provider is not initialized")
	}
	response, err := p.inner.GenerateText(ctx, convertCoreProviderRequest(request))
	if err != nil {
		return nil, err
	}
	return convertCoreGenerateResponse(response), nil
}

// convertCoreProviderRequest copies the structurally identical host request
// into the plugin-local request type.
func convertCoreProviderRequest(request coreaitext.ProviderRequest) aitext.ProviderRequest {
	messages := make([]aitext.Message, 0, len(request.Messages))
	for _, message := range request.Messages {
		messages = append(messages, aitext.Message{
			Role:    aitext.MessageRole(message.Role),
			Content: message.Content,
		})
	}
	return aitext.ProviderRequest{
		GenerateRequest: aitext.GenerateRequest{
			Purpose:         request.Purpose,
			Tier:            aitext.Tier(request.Tier),
			Messages:        messages,
			MaxOutputTokens: request.MaxOutputTokens,
			Temperature:     request.Temperature,
			ThinkingEffort:  convertCoreThinkingEffort(request.ThinkingEffort),
			Metadata:        request.Metadata,
		},
		SourcePluginID: request.SourcePluginID,
	}
}

// convertCoreThinkingEffort maps the optional host thinking-effort pointer.
func convertCoreThinkingEffort(effort *coreaitext.ThinkingEffort) *aitext.ThinkingEffort {
	if effort == nil {
		return nil
	}
	converted := aitext.ThinkingEffort(*effort)
	return &converted
}

// convertCoreGenerateResponse copies the plugin-local response into the host
// response type.
func convertCoreGenerateResponse(response *aitext.GenerateResponse) *coreaitext.GenerateResponse {
	if response == nil {
		return nil
	}
	converted := &coreaitext.GenerateResponse{
		Text:         response.Text,
		Tier:         coreaitext.Tier(response.Tier),
		ProviderName: response.ProviderName,
		ModelName:    response.ModelName,
		Protocol:     response.Protocol,
		Usage: coreaitext.Usage{
			InputTokens:  response.Usage.InputTokens,
			OutputTokens: response.Usage.OutputTokens,
		},
		LatencyMs:   response.LatencyMs,
		GeneratedAt: response.GeneratedAt,
	}
	if response.ThinkingEffort != nil {
		effort := coreaitext.ThinkingEffort(*response.ThinkingEffort)
		converted.ThinkingEffort = &effort
	}
	return converted
}
