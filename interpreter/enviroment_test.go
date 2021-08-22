package interpreter

import (
	"testing"

	"github.com/nanmu42/bluelox/token"
	"github.com/stretchr/testify/require"
)

func TestEnvironment_Define_Get(t *testing.T) {
	var (
		env = NewEnvironment()
		err error
	)

	err = env.Define("key1", "value1")
	require.NoError(t, err)

	result, err := env.Get(&token.Token{
		Type:    token.Identifier,
		Lexeme:  "key1",
		Literal: nil,
		Line:    0,
	})
	require.NoError(t, err)

	require.Equal(t, "value1", result)

	err = env.Define("key2", "value2")
	require.NoError(t, err)

	result, err = env.Get(&token.Token{
		Type:    token.Identifier,
		Lexeme:  "key2",
		Literal: nil,
		Line:    0,
	})
	require.NoError(t, err)

	require.Equal(t, "value2", result)
}
