package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Carina-labs/HAL9000/utils"
)

func TestSetJSONHeader(t *testing.T) {
	expected := http.Header{
		"Content-Type": {"application/json"},
	}

	w := httptest.NewRecorder()
	utils.SetJSONHeader(w)
	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, expected, res.Header, "SetJSONHeader not working")
}
