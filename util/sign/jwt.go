package sign

import (
	"github.com/dgrijalva/jwt-go"
	"relocate/util/conf"
	"time"
)

type ClaimsType string

const (
	AdminClaimsType ClaimsType = "admin-claims-type" //管理员
	UserClaimsType             = "user-claims-type"  //用户
)

type Claims struct {
	Type ClaimsType
	jwt.StandardClaims
}

func GenerateToken(id string, issuer string, claimsType ClaimsType, minute ...int) (string, error) {
	nowTime := time.Now()
	var expireTime time.Time
	if len(minute) > 0 {
		expireTime = nowTime.Add(time.Duration(minute[0]) * time.Minute)
	} else {
		expireTime = nowTime.Add(8 * time.Hour)
	}
	claims := Claims{
		claimsType,
		jwt.StandardClaims{
			Audience:  EncodeMD5(id),      // 受众
			ExpiresAt: expireTime.Unix(),  // 失效时间
			Id:        EncodeMD5(id),      // 编号
			IssuedAt:  time.Now().Unix(),  // 签发时间
			Issuer:    issuer,             // 签发人
			NotBefore: time.Now().Unix(),  // 生效时间
			Subject:   string(claimsType), // 主题
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString([]byte(conf.Data.App.JwtSecret))
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Data.App.JwtSecret), nil
	})
	if jwtToken != nil {
		if claims, ok := jwtToken.Claims.(*Claims); ok && jwtToken.Valid {
			return claims, nil
		}
	}
	return nil, err
}
