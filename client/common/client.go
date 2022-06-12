package common

import (
	"github.com/Carina-labs/HAL9000/utils"
	ut "github.com/Carina-labs/HAL9000/utils/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

func MakeClientWithNewAcc(ctx client.Context, keyname string, mnemonic string, bip44path string, algo keyring.SignatureAlgo) keyring.Info {
	accInfo := CreateAcc(ctx, keyname, mnemonic, bip44path, algo)
	return accInfo
}

func LoadClient(ctx client.Context, keyname string) keyring.Info {
	accInfo, err := ctx.Keyring.Key(keyname)
	utils.HandleErr(err, "Cannot load client with key", ut.EXIT)
	return accInfo
}

func CreateAcc(ctx client.Context, keyname string, mnemonic string, bip44path string, algo keyring.SignatureAlgo) keyring.Info {
	info, err := ctx.Keyring.NewAccount(keyname, mnemonic, "", bip44path, algo)
	utils.HandleErr(err, "Cannot create account with those arguments. Check it!", ut.EXIT)
	return info
}
