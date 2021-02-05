package config
import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("THISISsfdffdf!@@#@!)D_DFDFDS")

type UserInfo struct {
	Username string
	Password string
}