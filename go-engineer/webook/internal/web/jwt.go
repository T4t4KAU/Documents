package web

import "github.com/golang-jwt/jwt/v5"

// 专门用于 JWT 的代码

type UserClaims struct {
	// 我们只需要放一个 user id 就可以了
	Id        int64
	UserAgent string
	jwt.RegisteredClaims
}

// JWTKey 因为 JWT Key 不太可能变，所以可以直接写成常量
// 也可以考虑做成依赖注入
var JWTKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")
