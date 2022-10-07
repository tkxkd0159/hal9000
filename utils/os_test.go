package utils_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Carina-labs/HAL9000/utils"
)

func TestCheckErr(t *testing.T) {
	assert.Panics(t, func() { utils.CheckErr(errors.New(""), "", 0) })
}
