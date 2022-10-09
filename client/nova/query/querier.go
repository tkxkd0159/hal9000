package query

import (
	"context"
	"time"

	galtypes "github.com/Carina-labs/nova/x/gal/types"
	icatypes "github.com/Carina-labs/nova/x/icacontrol/types"
	"google.golang.org/grpc"

	"github.com/Carina-labs/HAL9000/client/nova/types"
	"github.com/Carina-labs/HAL9000/utils"
	utiltypes "github.com/Carina-labs/HAL9000/utils/types"
)

type NovaQueryClient struct {
	*grpc.ClientConn
}

func NewNovaQueryClient(grpcAddr string) *NovaQueryClient {
	conn, err := grpc.Dial(
		grpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	return &NovaQueryClient{conn}
}

var (
	_ types.NovaQuerier = &NovaQueryClient{}
)

const (
	ctxTimeout = time.Second * 5
)

func (nqc *NovaQueryClient) CurrentDelegateVersion(zoneid string) *galtypes.QueryCurrentDelegateVersionResponse {
	c := galtypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.DelegateCurrentVersion(ctx, &galtypes.QueryCurrentDelegateVersion{ZoneId: zoneid})
	utils.CheckErr(err, "", utiltypes.KEEP)
	return r
}
func (nqc *NovaQueryClient) CurrentUndelegateVersion(zoneid string) *galtypes.QueryCurrentUndelegateVersionResponse {
	c := galtypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.UndelegateCurrentVersion(ctx, &galtypes.QueryCurrentUndelegateVersion{ZoneId: zoneid})
	utils.CheckErr(err, "", utiltypes.KEEP)
	return r
}
func (nqc *NovaQueryClient) CurrentWithdrawVersion(zoneid string) *galtypes.QueryCurrentWithdrawVersionResponse {
	c := galtypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.WithdrawCurrentVersion(ctx, &galtypes.QueryCurrentWithdrawVersion{ZoneId: zoneid})
	utils.CheckErr(err, "", utiltypes.KEEP)
	return r

}
func (nqc *NovaQueryClient) CurrentAutoStakingVersion(zoneid string) *icatypes.QueryCurrentAutoStakingVersionResponse {
	c := icatypes.NewQueryClient(nqc)
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	r, err := c.AutoStakingCurrentVersion(ctx, &icatypes.QueryCurrentAutoStakingVersion{ZoneId: zoneid})
	utils.CheckErr(err, "", utiltypes.KEEP)
	return r
}
