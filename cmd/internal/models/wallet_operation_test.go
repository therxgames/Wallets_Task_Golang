package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOperationTypeIsValid(t *testing.T) {
	require.True(t, OperationType("DEPOSIT").IsValid())
	require.True(t, OperationType("WITHDRAW").IsValid())
	require.False(t, OperationType("REFUND").IsValid())
	require.False(t, OperationType("random").IsValid())
}
