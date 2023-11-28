package jwt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenToken(t *testing.T) {
	accessSecret := "GGjahoonIZH"
	accessExpire := 604800
	// refreshSecret := util.RandomString(6)
	// fields := make(map[string]interface{})
	// k := "userId"
	// //v, err := json.Marshal("12")
	// v := int64(12)
	// fields[k] = int64(v)
	// fmt.Println(fields[k])
	// opt := TokenOption{
	// 	AccessSecret: accessSecret,
	// 	AccessExpire: int64(604800),
	// 	// RefreshAfter:  int64(time.Minute * 5),
	// 	// RefreshSecret: refreshSecret,
	// 	// RefreshExpire: int64(time.Minute * 60),
	// 	Fields: fields,
	// }
	//token, err := BuildTokens(opt)
	jwtToken, err := BuildTokens(TokenOption{
		AccessSecret: accessSecret,
		AccessExpire: int64(accessExpire),
		Fields: map[string]interface{}{
			"userId": int64(31),
		},
	})
	require.NoError(t, err)
	//require.Equal(t, opt.AccessExpire/180000000000, token.AccessExpire/180000000000)
	//require.Equal(t, opt.RefreshExpire/3600000000000, token.RefreshExpire/3600000000000)
	// require.Equal(t, opt.RefreshAfter/300000000000, token.RefreshAfter/300000000000)
	fmt.Println(jwtToken)
}
