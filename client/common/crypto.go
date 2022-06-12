package common

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

func MakeKeyring(ctx client.Context, backend string) keyring.Keyring {
	kb, err := client.NewKeyringFromBackend(ctx, backend)
	utils.HandleErr(err, "Cannot generate keyring instance", types.EXIT)
	return kb
}
