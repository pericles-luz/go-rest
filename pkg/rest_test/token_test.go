package rest_test

import (
	"testing"
	"time"

	"github.com/pericles-luz/go-rest/pkg/rest"
	"github.com/stretchr/testify/require"
)

func TestTokenWithFutureValidityMustBeValid(t *testing.T) {
	token := rest.NewToken()
	token.SetKey("1234567890")
	token.SetValidity(time.Now().UTC().Add(time.Second * 1).Format("2006-01-02 15:04:05"))
	require.True(t, token.IsValid())
}

func TestTokenWithPastValidityMustBeInvalid(t *testing.T) {
	token := rest.NewToken()
	token.SetKey("1234567890")
	token.SetValidity(time.Now().UTC().Add(time.Second * -1).Format("2006-01-02 15:04:05"))
	require.False(t, token.IsValid())
}
