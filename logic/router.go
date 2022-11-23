package logic

import (
	"github.com/Carina-labs/HAL9000/client/base/query"
	basetypes "github.com/Carina-labs/HAL9000/client/base/types"
	novaq "github.com/Carina-labs/HAL9000/client/nova/query"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
)

func RouteBotAction(b *basetypes.Bot, cni *config.ChainNetInfo, hci *config.HostChainInfo) {
	initialBanner(b.Type)
	switch b.Type {
	case config.ActOracle:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		UpdateChainState(cq, b, hci)
	case config.ActStake:
		nq := novaq.NewNovaQueryClient(cni.GRPC.Host)
		defer utils.CloseGrpc(nq.ClientConn)
		IcaStake(nq, b, hci)
	case config.ActAutoStake:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr)
		nq := novaq.NewNovaQueryClient(cni.GRPC.Host)
		defer utils.CloseGrpc(cq.ClientConn)
		defer utils.CloseGrpc(nq.ClientConn)
		IcaAutoStake(cq, nq, b, hci)
	case config.ActWithdraw:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr)
		nq := novaq.NewNovaQueryClient(cni.GRPC.Host)
		defer utils.CloseGrpc(cq.ClientConn)
		defer utils.CloseGrpc(nq.ClientConn)
		UndelegateAndWithdraw(cq, nq, b, hci)
	default:
		panic("This type cannot handle at this action router")
	}
}
