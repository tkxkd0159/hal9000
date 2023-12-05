package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tkxkd0159/HAL9000/client/base/types"
)

func TestConvertMainPrefix(t *testing.T) {
	old := "nova"
	_new := "newprefix"
	got := types.ConvertMainPrefix(types.Bech32MainPrefix, old, _new)
	got2 := types.ConvertMainPrefix(types.Bech32PrefixAccAddr, old, _new)
	got3 := types.ConvertMainPrefix(types.Bech32PrefixValAddr, old, _new)
	got4 := types.ConvertMainPrefix(types.Bech32PrefixConsAddr, old, _new)
	require.Equal(t, "newprefix", got, "")
	require.Equal(t, "newprefix", got2, "")
	require.Equal(t, "newprefixvaloper", got3, "")
	require.Equal(t, "newprefixvalcons", got4, "")
}
