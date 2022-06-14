package common

import (
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/Carina-labs/HAL9000/utils/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

// MakeKeyring returns keystore from cosmos-sdk/crypto/keyring
// If you want to reset keyring, delete keyring-* dir
func MakeKeyring(ctx client.Context, backend string) keyring.Keyring {
	kb, err := newKeyringFromBackend(ctx, backend)
	utils.CheckErr(err, "Cannot generate keyring instance", types.EXIT)
	return kb
}

func newKeyringFromBackend(ctx client.Context, backend string) (keyring.Keyring, error) {
	if ctx.GenerateOnly || ctx.Simulate {
		return keyring.New(sdktypes.KeyringServiceName(), keyring.BackendMemory, ctx.KeyringDir, ctx.Input, ctx.KeyringOptions...)
	}

	return keyring.New(sdktypes.KeyringServiceName(), backend, ctx.KeyringDir, ctx.Input, ctx.KeyringOptions...)
}
