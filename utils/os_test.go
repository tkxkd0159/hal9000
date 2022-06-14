package utils_test

import (
	"errors"
	"github.com/Carina-labs/HAL9000/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckErr(t *testing.T) {
	assert.Panics(t, func() { utils.CheckErr(errors.New(""), "", 0) })
}
