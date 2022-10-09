package logic

import (
	"fmt"

	"github.com/Carina-labs/HAL9000/client/base/query"
	basetypes "github.com/Carina-labs/HAL9000/client/base/types"
	"github.com/Carina-labs/HAL9000/config"
	"github.com/Carina-labs/HAL9000/utils"
)

func RouteBotAction(botType string, b *basetypes.Bot, hci *config.HostChainInfo) {
	fmt.Printf("\n 🤖 %s bot has started working... 🤖 \n", botType)
	switch botType {
	case config.ActOracle:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		UpdateChainState(cq, b, hci)
	case config.ActStake:
		IcaStake(b, hci)
	case config.ActRestake:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		IcaAutoStake(cq, b, hci)
	case config.ActWithdraw:
		cq := query.NewCosmosQueryClient(hci.GrpcAddr)
		defer utils.CloseGrpc(cq.ClientConn)
		UndelegateAndWithdraw(cq, b, hci)
	default:
		panic("This type cannot handle at this action router")
	}
}
