package base

import (
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/tkxkd0159/HAL9000/utils"
	ut "github.com/tkxkd0159/HAL9000/utils/types"
)

func LoadClientPubInfo(ctx client.Context, keyname string) keyring.Info {
	accInfo, err := ctx.Keyring.Key(keyname)
	utils.CheckErr(err, "Cannot load client with key", ut.EXIT)
	return accInfo
}

func MakeClientWithNewAcc(ctx client.Context, keyname string, mnemonic string, bip44path string, algo keyring.SignatureAlgo) keyring.Info {
	accInfo := createAcc(ctx, keyname, mnemonic, bip44path, algo)
	return accInfo
}

func createAcc(ctx client.Context, keyname string, mnemonic string, bip44path string, algo keyring.SignatureAlgo) keyring.Info {
	info, err := ctx.Keyring.NewAccount(keyname, mnemonic, "", bip44path, algo)
	utils.CheckErr(err, " * cannot create account with those arguments. Check it!", ut.EXIT)
	return info
}

func CheckAccAddr(target any) (AccAddr, error) {
	switch target := target.(type) {
	case AccAddr:
		return target, nil
	case string:
		addr, err := sdktypes.AccAddressFromBech32(target)
		if err != nil {
			return nil, err
		}
		return addr, nil
	case []byte:
		return target, nil
	default:
		return nil, errors.New("cannot covert target to AccAddress")
	}
}

func GetPrivKey(ctx client.Context, keyname string) cryptotypes.PrivKey {
	privArmor, _ := ctx.Keyring.ExportPrivKeyArmor(keyname, "")
	novaPrivRaw, _, _ := crypto.UnarmorDecryptPrivKey(privArmor, "")
	return novaPrivRaw
}
