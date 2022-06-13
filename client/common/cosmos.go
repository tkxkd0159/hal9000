package common

import (
	"context"
	"fmt"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	tendermintv1beta1 "github.com/Carina-labs/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-labs/nova/api/cosmos/staking/v1beta1"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"time"
)

func GetNodeInfo(conn *grpc.ClientConn) *tendermintv1beta1.GetNodeInfoResponse {
	c := tendermintv1beta1.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	r, err := c.GetNodeInfo(ctx, &tendermintv1beta1.GetNodeInfoRequest{})
	utils.CheckErr(err, "Can't get Node info", types.KEEP)
	// log.Printf("From gRPC srv : \n %+v", r.GetApplicationVersion())
	//log.Printf("From gRPC srv : \n %+v", r.GetNodeInfo())
	return r
}

func GetValInfo(conn *grpc.ClientConn, valAddr string) *stakingv1beta1.QueryValidatorResponse {
	c := stakingv1beta1.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Validator(ctx, &stakingv1beta1.QueryValidatorRequest{ValidatorAddr: valAddr})
	utils.CheckErr(err, fmt.Sprintf("Can't get %s info", valAddr), types.KEEP)
	// r.GetValidator().Tokens
	return r
}

// Deprecated: SendTx
// 1. Generate a TX with Msg (TxBuilder)
// 2. Sign the generated transaction with the keyring's account
// 3. Broadcast the tx using gPRC
func SendTx(ctx client.Context, keyname string) {
	encCfg := MakeEncodingConfig(novaapp.ModuleBasics)
	txBuilder := encCfg.TxConfig.NewTxBuilder()
	from := sdktypes.AccAddress(viper.GetString("nova.local_addr"))
	to := sdktypes.AccAddress(viper.GetString("nova.target_addr"))
	coin := sdktypes.Coin{Denom: "nova", Amount: sdktypes.NewInt(333)}
	msg1 := MakeSendMsg(from, to, coin)
	coin2 := sdktypes.Coin{Denom: "nova", Amount: sdktypes.NewInt(222)}
	msg2 := MakeSendMsg(from, to, coin2)

	err := txBuilder.SetMsgs(msg1, msg2)
	utils.CheckErr(err, "", 0)
	txBuilder.SetMemo("")
	priv := GetPrivKey(ctx, keyname)
	_ = priv

}

func MakeSendMsg(from sdktypes.AccAddress, to sdktypes.AccAddress, coin sdktypes.Coin) *banktypes.MsgSend {
	return banktypes.NewMsgSend(from, to, sdktypes.NewCoins(sdktypes.NewCoin(coin.Denom, coin.Amount)))
}

func GenTxWithFactory(ctx client.Context, txf tx.Factory, onlyGen bool, msgs ...sdktypes.Msg) {
	if onlyGen {
		// build unsigned tx
		ctx = ctx.WithGenerateOnly(true)
	}
	err := tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msgs...)
	utils.CheckErr(err, "something went wrong while make tx", 0)
}

func GetPrivKey(ctx client.Context, keyname string) cryptotypes.PrivKey {
	privArmor, _ := ctx.Keyring.ExportPrivKeyArmor(keyname, "")
	novaPrivRaw, _, _ := crypto.UnarmorDecryptPrivKey(privArmor, "")
	return novaPrivRaw
}
