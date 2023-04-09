package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/wecoding/iam/internal/apiserver/config"
	"github.com/wecoding/iam/internal/apiserver/domain/model"
	"github.com/wecoding/iam/internal/apiserver/domain/repository"
	convert "github.com/wecoding/iam/internal/apiserver/interfaces/api/convert/v1alpha1"
	"github.com/wecoding/iam/internal/pkg/code"
	"github.com/wecoding/iam/pkg/api/apiserver/v1alpha1"
)

const (
	jwtIssuer = "iam-issuer"
	audience  = "iam.api.wecoding.top"

	// GrantTypeAccess is the grant type for access token
	GrantTypeAccess = "access"
	// GrantTypeRefresh is the grant type for refresh token
	GrantTypeRefresh = "refresh"
)

// AuthenticationService authentication service
type AuthenticationService interface {
	Authenticate(ctx context.Context, loginReq v1alpha1.AuthenticateRequest) (*v1alpha1.AuthenticateResponse, error)
}

type authenticationServiceImpl struct {
	cfg         config.Config
	Store       repository.Factory `inject:"repository"`
	UserService UserService        `inject:""`
}

// NewAuthenticationService new authentication service
func NewAuthenticationService(c config.Config) AuthenticationService {
	return &authenticationServiceImpl{cfg: c}
}

type authHandler interface {
	authenticate(ctx context.Context) (*v1alpha1.UserBase, error)
}

type localHandlerImpl struct {
	store       repository.Factory
	userService UserService
	username    string
	password    string
}

func (a *authenticationServiceImpl) newLocalHandler(loginReq v1alpha1.AuthenticateRequest) (*localHandlerImpl, error) {
	if loginReq.Username == "" || loginReq.Password == "" {
		return nil, errors.WithCode(code.ErrMissingLoginValues, "Missing Username or Password")
	}

	return &localHandlerImpl{
		store:       a.Store,
		userService: a.UserService,
		username:    loginReq.Username,
		password:    loginReq.Password,
	}, nil
}

func (a *authenticationServiceImpl) Authenticate(ctx context.Context, loginReq v1alpha1.AuthenticateRequest) (*v1alpha1.AuthenticateResponse, error) {
	var handler authHandler
	var err error
	handler, err = a.newLocalHandler(loginReq)
	if err != nil {
		return nil, err
	}
	userBase, err := handler.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	accessToken, err := a.generateJWTToken(userBase.Name, GrantTypeAccess, a.cfg.JwtOptions.Timeout)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateJWTToken(userBase.Name, GrantTypeRefresh, a.cfg.JwtOptions.MaxRefresh)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.AuthenticateResponse{
		User:         userBase,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authenticationServiceImpl) generateJWTToken(username, grantType string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)
	claims := model.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Audience:  jwt.ClaimStrings{audience},
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expire),
		},
		Username:  username,
		GrantType: grantType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(a.cfg.JwtOptions.Key))
}

func (l *localHandlerImpl) authenticate(ctx context.Context) (*v1alpha1.UserBase, error) {
	user, err := l.userService.Get(ctx, l.username, metav1alpha1.GetOptions{})
	if err != nil {
		if errors.IsCode(err, code.ErrUserNotFound) {
			return nil, errors.WithCode(code.ErrPasswordIncorrect, "Password was incorrect")
		}
		return nil, err
	}
	if err := passwordVerify(user.Password, l.password); err != nil {
		return nil, err
	}

	return convert.ConvertUserModelToBase(user), nil
}

func passwordVerify(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return errors.WithCode(code.ErrPasswordIncorrect, "Password was incorrect")
	}

	return err
}
