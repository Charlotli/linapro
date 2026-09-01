// Package bridge provides the linapro-ai-core dynamic-plugin guest SDK for
// owner-aware AI host services. It owns AI request codecs, method declaration
// helpers, and typed guest clients while all provider implementations and
// business storage remain inside the linapro-ai-core backend internals.
package bridge

import (
	"lina-plugin-linapro-ai-core/backend/cap/aicap/spi"
)

// AllMethods returns every AI method currently published by the owner.
func AllMethods() []string {
	return spi.PublishedMethods()
}

// TextMethods returns the text AI methods normally needed by a dynamic plugin
// that only generates text and reads text method status.
func TextMethods() []string {
	return []string{
		spi.MethodTextGenerate,
		spi.MethodTextStatusGet,
	}
}

// StatusMethods returns cross-capability AI status methods for dynamic plugins
// that need method-level degradation checks.
func StatusMethods() []string {
	return []string{
		spi.MethodTextStatusGet,
		spi.MethodStatusesBatchGet,
	}
}
