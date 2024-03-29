package query

import (
	"context"
	"time"

	"google.golang.org/grpc"

	galtypes "github.com/Carina-labs/nova/x/gal/types"
	icatypes "github.com/Carina-labs/nova/x/icacontrol/types"

	"github.com/tkxkd0159/HAL9000/client"
	"github.com/tkxkd0159/HAL9000/client/nova/types"
	"github.com/tkxkd0159/HAL9000/utils"
)

const (
	ctxTimeout = time.Second * 10
)

var _ types.NovaQuerier = &NovaQueryClient{}

type NovaQueryClient struct {
	*grpc.ClientConn
}

func NewNovaQueryClient(grpcAddr string, secure bool) *NovaQueryClient {
	conn, err := client.OpengRPCConn(grpcAddr, secure)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	return &NovaQueryClient{conn}
}

func (nqc *NovaQueryClient) CurrentDelegateVersion(zoneid string) (*galtypes.QueryCurrentDelegateVersionResponse, error) {
	c := galtypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.DelegateCurrentVersion(ctx, &galtypes.QueryCurrentDelegateVersion{ZoneId: zoneid})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (nqc *NovaQueryClient) CurrentUndelegateVersion(zoneid string) (*galtypes.QueryCurrentUndelegateVersionResponse, error) {
	c := galtypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.UndelegateCurrentVersion(ctx, &galtypes.QueryCurrentUndelegateVersion{ZoneId: zoneid})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (nqc *NovaQueryClient) CurrentWithdrawVersion(zoneid string) (*galtypes.QueryCurrentWithdrawVersionResponse, error) {
	c := galtypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.WithdrawCurrentVersion(ctx, &galtypes.QueryCurrentWithdrawVersion{ZoneId: zoneid})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (nqc *NovaQueryClient) CurrentAutoStakingVersion(zoneid string) (*icatypes.QueryCurrentAutoStakingVersionResponse, error) {
	c := icatypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.AutoStakingCurrentVersion(ctx, &icatypes.QueryCurrentAutoStakingVersion{ZoneId: zoneid})
	if err != nil {
		return nil, err
	}
	return r, nil
}
