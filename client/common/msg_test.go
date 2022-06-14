package common_test

import (
	"github.com/Carina-labs/HAL9000/client/common"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NovaTestSuite struct {
	suite.Suite
}

func TestNovaTestSuite(t *testing.T) {
	suite.Run(t, new(NovaTestSuite))
}

func (s *NovaTestSuite) SetupSuite() {
	common.SetBechPrefix()
}

func (s *NovaTestSuite) TestCheckAccAddr() {
	tcs := []struct {
		expected sdktypes.AccAddress
		input    any
	}{
		{
			sdktypes.AccAddress([]byte{0xfb, 0x61, 0x43, 0xb4, 0x68, 0x3f, 0xd9, 0x3e, 0x84, 0x78, 0xe6, 0x55, 0x0, 0x21, 0x2c, 0x94, 0xc5, 0x4f, 0xcb, 0xf6}),
			"nova1lds58drg8lvnaprcue2sqgfvjnz5ljlkq9lsyf",
		},
		{
			sdktypes.AccAddress("FB6143B4683FD93E8478E65500212C94C54FCBF6"),
			[]byte("FB6143B4683FD93E8478E65500212C94C54FCBF6"),
		},
		{
			nil,
			1,
		},
	}
	for _, tc := range tcs {
		got, err := common.CheckAccAddr(tc.input)
		assert.Equal(s.T(), tc.expected, got)
		if got == nil {
			assert.Errorf(s.T(), err, "cannot covert target to AccAddress")
		} else {
			assert.NoError(s.T(), err)
		}
	}
}
