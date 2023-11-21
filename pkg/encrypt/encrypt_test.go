package encrypt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptMobil(t *testing.T) {
	mobile := "17615164896"
	encrptyMobile, err := EncMobile(mobile)
	require.NoError(t, err)
	decMobile, err := DecMobile(encrptyMobile)
	require.NoError(t, err)
	require.Equal(t, mobile, decMobile)
}
