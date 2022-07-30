package types_test

import (
	"github.com/Carina-labs/HAL9000/client/base/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvertMainPrefix(t *testing.T) {
	old := "nova"
	new := "newprefix"
	got := types.ConvertMainPrefix(types.Bech32MainPrefix, old, new)
	got2 := types.ConvertMainPrefix(types.Bech32PrefixAccAddr, old, new)
	got3 := types.ConvertMainPrefix(types.Bech32PrefixValAddr, old, new)
	got4 := types.ConvertMainPrefix(types.Bech32PrefixConsAddr, old, new)
	require.Equal(t, "newprefix", got, "")
	require.Equal(t, "newprefix", got2, "")
	require.Equal(t, "newprefixvaloper", got3, "")
	require.Equal(t, "newprefixvalcons", got4, "")
}
