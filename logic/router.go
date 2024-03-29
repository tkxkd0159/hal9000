package logic

import (
	"github.com/tkxkd0159/HAL9000/client/base/query"
	basetypes "github.com/tkxkd0159/HAL9000/client/base/types"
	novaq "github.com/tkxkd0159/HAL9000/client/nova/query"
	"github.com/tkxkd0159/HAL9000/config"
	"github.com/tkxkd0159/HAL9000/utils"
)

func RouteBotAction(b *basetypes.Bot, cni *config.ChainNetInfo, hci *config.HostChainInfo) {
	initialBanner(b.Type)

	switch b.Type {
	case config.ActOracle:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr, cni.Secure)
		defer utils.CloseGrpc(cq.ClientConn)
		UpdateChainState(cq, b, hci)
	case config.ActStake:
		nq := novaq.NewNovaQueryClient(cni.GRPC.Host, cni.Secure)
		defer utils.CloseGrpc(nq.ClientConn)
		IcaStake(nq, b, hci)
	case config.ActAutoStake:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr, cni.Secure)
		nq := novaq.NewNovaQueryClient(cni.GRPC.Host, cni.Secure)
		defer utils.CloseGrpc(cq.ClientConn)
		defer utils.CloseGrpc(nq.ClientConn)
		IcaAutoStake(cq, nq, b, hci)
	case config.ActWithdraw:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr, cni.Secure)
		nq := novaq.NewNovaQueryClient(cni.GRPC.Host, cni.Secure)
		defer utils.CloseGrpc(cq.ClientConn)
		defer utils.CloseGrpc(nq.ClientConn)
		UndelegateAndWithdraw(cq, nq, b, hci)
	case config.ActAutoClaim:
		ClaimAllSnAsset(b, hci)
	default:
		panic("This type cannot handle at this action router")
	}
}
