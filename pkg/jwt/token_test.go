package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenToken(t *testing.T) {
	accessSecret := "GGjahoonIZH"
	// refreshSecret := util.RandomString(6)
	fields := make(map[string]interface{})
	k := "userId"
	//v, err := json.Marshal("12")
	v := 12
	fields[k] = int64(v)
	opt := TokenOption{
		AccessSecret: accessSecret,
		AccessExpire: int64(604800),
		// RefreshAfter:  int64(time.Minute * 5),
		// RefreshSecret: refreshSecret,
		// RefreshExpire: int64(time.Minute * 60),
		Fields: fields,
	}
	token, err := BuildTokens(opt)
	require.NoError(t, err)
	require.Equal(t, opt.AccessExpire/180000000000, token.AccessExpire/180000000000)
	//require.Equal(t, opt.RefreshExpire/3600000000000, token.RefreshExpire/3600000000000)
	// require.Equal(t, opt.RefreshAfter/300000000000, token.RefreshAfter/300000000000)
	fmt.Println(token)
}
