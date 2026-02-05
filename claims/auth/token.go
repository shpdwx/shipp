package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shpdwx/claims/common"
)

type AccessClaim struct {
	UserId   int64  `json:"uid"`
	Username string `json:"un"`
	Jti      string `json:"jti"`
	jwt.RegisteredClaims
}

type RefreshClaim struct {
	UserId   int64  `json:"uid"`
	Jti      string `json:"jti"`
	DeviceId string `json:"did"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	Token   string `json:"access_token"`
	Refresh string `json:"refresh_token"`
	Expires int    `json:"expires_in"`
}

type JwtToken interface {
	Gen() (*TokenPair, error)
	User(userId int64, username string)
	Device(str string)
	Validate(jwtStr string) (err error)
	Refresh(refresh string) (*TokenPair, error)
}

type defJwtToken struct {
	ctx             context.Context
	AppName         string
	UserId          int64
	Username        string
	DeviceId        string //user-agent
	JtiNum          int
	Duration        int //minute
	RefreshDuration int //day
	JwtSecret       []byte
	TokenNumber     int64
}

func NewJwtToken(ctx context.Context, app string) JwtToken {
	return &defJwtToken{
		ctx:             ctx,
		AppName:         app,
		JtiNum:          2,
		Duration:        15,
		RefreshDuration: 14,
		JwtSecret:       []byte("89d78bd0-9c48-4ab6-96cb-9d067c761164"),
		TokenNumber:     5,
	}
}

func (t *defJwtToken) User(userId int64, username string) {
	t.UserId = userId
	t.Username = username
}

func (t *defJwtToken) Device(str string) {
	if str == "" {
		return
	}
	t.DeviceId = base64.RawURLEncoding.EncodeToString([]byte(str))
}

// 登陆流程
func (t *defJwtToken) Gen() (*TokenPair, error) {
	rsp := TokenPair{}

	id := uuid.NewString()
	jti := base64.RawURLEncoding.EncodeToString([]byte(id))

	// 限制用户Token数量
	nc := NewCache(t.ctx)
	if err := nc.LimitTokens(t.TokenNumber, t.UserId, jti); err != nil {
		return nil, err
	}

	// 生成access token
	ac := AccessClaim{
		UserId:   t.UserId,
		Username: t.Username,
		Jti:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.Duration) * time.Minute)),
			Issuer:    t.AppName,
			Subject:   strconv.FormatInt(t.UserId, 10),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)
	if str, err := at.SignedString(t.JwtSecret); err != nil {
		return nil, err
	} else {
		rsp.Token = str
	}

	// 生成refresh token
	rc := RefreshClaim{
		UserId:   t.UserId,
		Jti:      jti,
		DeviceId: t.DeviceId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.RefreshDuration) * time.Hour * 24)),
			Issuer:    t.AppName,
			Subject:   t.Username,
		},
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	if str, err := rt.SignedString(t.JwtSecret); err != nil {
		return nil, err
	} else {
		rsp.Refresh = base64.StdEncoding.EncodeToString([]byte(str))
	}

	// 缓存refresh token info
	refreshKey := fmt.Sprintf("refresh_token:%s", jti)
	rdb := nc.Rdb()

	b, err := json.Marshal(rc)
	if err != nil {
		return nil, err
	}
	rdb.SetEx(t.ctx, refreshKey, string(b), time.Duration(t.RefreshDuration)*time.Hour*24)

	rsp.Expires = int(ac.ExpiresAt.Unix())
	return &rsp, nil
}

// 验证token
func (t *defJwtToken) Validate(str string) (err error) {
	if str == "" {
		return errors.New("Token不能为空")
	}

	pr := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(t.AppName),
		jwt.WithLeeway(5*time.Second),
		jwt.WithExpirationRequired(),
	)

	claims := &AccessClaim{}

	token, err := pr.ParseWithClaims(str, claims, func(j *jwt.Token) (any, error) {
		if _, ok := j.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", j.Header["alg"])
		}
		return t.JwtSecret, nil
	})

	fmt.Println(err)
	fmt.Println(token.Valid)
	return
}

// 刷新token
func (f *defJwtToken) Refresh(str string) (*TokenPair, error) {

	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	// 校验refresh token
	claim := RefreshClaim{}
	if err := f.check(string(b), &claim); err != nil {
		return nil, err
	}

	// 新 jti
	oldJti := claim.Jti
	claim.Jti = f.jti()

	ac := AccessClaim{
		UserId:           claim.UserId,
		Username:         claim.RegisteredClaims.Subject,
		Jti:              claim.Jti,
		RegisteredClaims: f.baseClaims(time.Now().Add(time.Duration(f.Duration)*time.Minute), strconv.FormatInt(f.UserId, 10)),
	}

	// 限制用户Token数量
	rdb := common.InitRedis(f.ctx)
	if err := f.limitTokens(rdb, ac, -1, oldJti); err != nil {
		return nil, err
	}

	claim.RegisteredClaims = f.baseClaims(time.Now().Add(time.Duration(f.RefreshDuration)*time.Hour*24), f.Username)

	rsp := TokenPair{}

	// 新 access token
	rsp.Token, err = f.genT(ac)

	if err != nil {
		return nil, err
	}

	// 新 refresh token
	rsp.Refresh, err = f.genRT(claim)
	if err != nil {
		return nil, err
	}

	if err := f.cacheRT(rdb, claim); err != nil {
		return nil, err
	}
	rsp.Expires = int(ac.ExpiresAt.Unix())
	return &rsp, nil
}

func (f *defJwtToken) check(token string, claims jwt.Claims) error {

	if token == "" {
		return errors.New("验签为空")
	}

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(f.AppName),
		jwt.WithLeeway(5*time.Second),
		jwt.WithExpirationRequired(),
	)

	result, err := parser.ParseWithClaims(token, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("验签算法错误: %v", t.Header["alg"])
			}
			return f.JwtSecret, nil
		})

	if err != nil {
		return fmt.Errorf("验签失败: %v", err)
	}

	if !result.Valid {
		return errors.New("验签不通过")
	}

	return nil
}

func (f *defJwtToken) jti() string {
	id := uuid.NewString()
	return base64.RawURLEncoding.EncodeToString([]byte(id))
}

func (f *defJwtToken) baseClaims(exp time.Time, sub string) jwt.RegisteredClaims {

	return jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(exp),
		Issuer:    f.AppName,
		Subject:   sub,
	}
}

func (f *defJwtToken) genT(user jwt.Claims) (token string, err error) {

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, user)
	token, err = t.SignedString(f.JwtSecret)
	if err != nil {
		err = fmt.Errorf("验签生成失败:%v", err)
		return
	}
	return
}

func (f *defJwtToken) genRT(user RefreshClaim) (str string, err error) {

	token, err := f.genT(user)
	if err != nil {
		return
	}
	str = base64.StdEncoding.EncodeToString([]byte(token))
	return
}

func (f *defJwtToken) cacheRT(rdb *redis.Client, user RefreshClaim) error {

	refreshKey := fmt.Sprintf("refresh_token:%s", user.Jti)

	b, err := json.Marshal(user)
	if err != nil {
		return err
	}
	rdb.SetEx(f.ctx, refreshKey, string(b), time.Duration(f.RefreshDuration)*time.Hour*24)
	return nil
}

func (f *defJwtToken) limitTokens(rdb *redis.Client, user AccessClaim, flush int, oldJti string) error {

	key := fmt.Sprintf("user_tokens:%d", user.UserId)

	switch flush {
	case 0:
		rdb.Del(f.ctx, key)
	case -1:
		rdb.SRem(f.ctx, key, oldJti)
	default:
		num, err := rdb.SCard(f.ctx, key).Result()
		if err != nil {
			return err
		}
		if num >= f.TokenNumber {
			return errors.New("超出限制的Token数量")
		}
	}

	if oldJti != "" {
		refreshKey := fmt.Sprintf("refresh_token:%s", oldJti)
		rdb.Del(f.ctx, refreshKey)
	}

	rdb.SAdd(f.ctx, key, user.Jti)
	return nil
}
