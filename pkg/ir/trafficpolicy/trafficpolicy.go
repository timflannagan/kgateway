package trafficpolicy

import (
	"context"

	internal "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugins/routepolicy"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
)

// NewGatewayTranslationPass exposes the traffic policy pass.
func NewGatewayTranslationPass(ctx context.Context, tctx ir.GwTranslationCtx) ir.ProxyTranslationPass {
	return internal.NewGatewayTranslationPass(ctx, tctx)
}
