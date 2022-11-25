package msgs

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func Validate(msgs ...sdktypes.Msg) error {
	for _, m := range msgs {
		if err := m.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
