package common_test

import (
	"bufio"
	"github.com/Carina-labs/HAL9000/client/common"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockio struct{}

func (mockio) Write(p []byte) (n int, err error) { return len(p), nil }
func (mockio) Read(p []byte) (n int, err error)  { return len(p), nil }

func TestMakeContext(t *testing.T) {
	common.SetBechPrefix()
	encCfg := common.MakeEncodingConfig(novaapp.ModuleBasics)

	tcs := []client.Context{
		{
			Simulate:          false,
			GenerateOnly:      true,
			SkipConfirm:       true,
			SignModeStr:       "direct",
			AccountRetriever:  authtypes.AccountRetriever{},
			BroadcastMode:     "sync",
			From:              "nova1z36nmc2efth7wy3dcnjsw2tu83qn5mxyydu663",
			FromName:          "",
			NodeURI:           "tcp://127.0.0.1:7777",
			ChainID:           "nova",
			HomeDir:           "home",
			KeyringDir:        "home",
			Input:             bufio.NewReader(mockio{}),
			Output:            mockio{},
			Codec:             encCfg.Marshaler,
			InterfaceRegistry: encCfg.InterfaceRegistry,
			TxConfig:          encCfg.TxConfig,
			LegacyAmino:       encCfg.Amino,
		},
	}

	for _, tc := range tcs {
		tc.FromAddress, _ = sdktypes.AccAddressFromBech32(tc.From)

		tc.Keyring = common.MakeKeyring(tc, "test")
		tmc, err := client.NewClientFromNode(tc.NodeURI)
		assert.NoError(t, err)
		tc.Client = tmc

		got := common.MakeContext(
			novaapp.ModuleBasics,
			tc.From,
			tc.NodeURI,
			tc.ChainID,
			tc.KeyringDir,
			"test",
			mockio{},
			mockio{},
			true,
		)
		got = common.AddMoreFromInfo(got)

		assert.NotEqual(t, tc.TxConfig, got.TxConfig)
		assert.NotEqual(t, tc.Client, got.Client)

		//assert.Equal(t, tc.LegacyAmino, got.LegacyAmino)
		assert.Equal(t, tc.HomeDir, got.HomeDir)
		assert.Equal(t, tc.Keyring, got.Keyring)
		assert.Equal(t, tc.KeyringDir, got.KeyringDir)
		assert.Equal(t, tc.AccountRetriever, got.AccountRetriever)
		assert.Equal(t, tc.Codec, got.Codec)
		assert.Equal(t, tc.InterfaceRegistry, got.InterfaceRegistry)
		assert.Equal(t, tc.BroadcastMode, got.BroadcastMode)
		assert.Equal(t, tc.SignModeStr, got.SignModeStr)
		assert.Equal(t, tc.Simulate, got.Simulate)
		assert.Equal(t, tc.FromAddress, got.FromAddress)
		assert.Equal(t, tc.FromName, got.FromName)
		assert.Equal(t, tc.From, got.From)
		assert.Equal(t, tc.ChainID, got.ChainID)
		assert.Equal(t, tc.NodeURI, got.NodeURI)
		assert.Equal(t, tc.SkipConfirm, got.SkipConfirm)
		assert.Equal(t, tc.Input, got.Input)
		assert.Equal(t, tc.Output, got.Output)
		assert.Equal(t, tc.GenerateOnly, got.GenerateOnly)

	}
}
