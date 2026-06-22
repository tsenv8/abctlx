package abctlx

import (
	"abctlx/internal/abctlx/cmd/destination"
	"abctlx/internal/abctlx/cmd/sources"
	"abctlx/internal/airbyte"
	"context"
)

type CmdHandler interface{}
type cmdHandler struct {
	AirbyteSvc airbyte.AirbyteService
	DestSvc    destination.Service
	SrcSvc     sources.Service
}

func NewCmdHandler(ctx context.Context) *cmdHandler {
	abSvc := airbyte.NewAirbyteService(ctx)
	destSvc := destination.NewService(abSvc.GetClient())
	srcSvc := sources.NewService(abSvc.GetClient())
	return &cmdHandler{
		AirbyteSvc: abSvc,
		DestSvc:    destSvc,
		SrcSvc:     srcSvc,
	}
}
