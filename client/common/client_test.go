package common_test

import (
	"bufio"
	"github.com/Carina-labs/HAL9000/client/common"
	novaapp "github.com/Carina-labs/nova/app"
	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockio struct{}

func (mockio) Write(p []byte) (n int, err error) { return len(p), nil }
func (mockio) Read(p []byte) (n int, err error)  { return len(p), nil }

func TestMakeContext(t *testing.T) {
	encCfg := common.MakeEncodingConfig(novaapp.ModuleBasics)

	tcs := []client.Context{
		{
			Simulate:          false,
			From:              "nova1z36nmc2efth7wy3dcnjsw2tu83qn5mxyydu663",
			Input:             bufio.NewReader(mockio{}),
			Output:            mockio{},
			NodeURI:           "tcp://127.0.0.1:7777",
			JSONCodec:         encCfg.Marshaler,
			Codec:             encCfg.Marshaler,
			InterfaceRegistry: encCfg.InterfaceRegistry,
			TxConfig:          encCfg.TxConfig,
			LegacyAmino:       encCfg.Amino,
			SignModeStr:       "direct",
			AccountRetriever:  authtypes.AccountRetriever{},
			BroadcastMode:     "sync",
			KeyringDir:        "home",
			SkipConfirm:       true,
			ChainID:           "nova",
		},
	}

	for _, tc := range tcs {
		tc.Keyring = common.MakeKeyring(tc, "test")
		tmc, err := client.NewClientFromNode(tc.NodeURI)
		assert.NoError(t, err)
		tc.Client = tmc

		got := common.MakeContext(
			novaapp.ModuleBasics,
			"nova1z36nmc2efth7wy3dcnjsw2tu83qn5mxyydu663",
			"tcp://127.0.0.1:7777",
			"nova",
			"home",
			"test",
			mockio{},
			mockio{},
		)
		assert.Equal(t, tc.ChainID, got.ChainID)
		assert.Equal(t, tc.From, got.From)
		assert.Equal(t, tc.NodeURI, got.NodeURI)
		assert.Equal(t, tc.KeyringDir, got.KeyringDir)
		assert.Equal(t, tc.SkipConfirm, got.SkipConfirm)
		assert.Equal(t, tc.Input, got.Input)
		assert.Equal(t, tc.Output, got.Output)

	}
}
