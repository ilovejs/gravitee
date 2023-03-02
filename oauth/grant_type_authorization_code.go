package oauth

import (
	"errors"
	"net/http"

	"github.com/spring2go/gravitee/models"
	"github.com/spring2go/gravitee/oauth/tokentypes"
)

var (
	// ErrInvalidRedirectURI ...
	ErrInvalidRedirectURI = errors.New("invalid redirect URI")
)

func (s *Service) authorizationCodeGrant(r *http.Request, client *models.OauthClient) (*AccessTokenResponse, error) {
	// Fetch the authorization code
	authorizationCode, err := s.getValidAuthorizationCode(
		r.Form.Get("code"),
		r.Form.Get("redirect_uri"),
		client,
	)
	if err != nil {
		return nil, err
	}

	// Log in the user
	accessToken, refreshToken, err := s.Login(
		authorizationCode.Client,
		authorizationCode.User,
		authorizationCode.Scope,
	)
	if err != nil {
		return nil, err
	}

	// Delete the authorization code
	s.db.Unscoped().Delete(&authorizationCode)

	// Create response
	accessTokenResponse, err := NewAccessTokenResponse(
		accessToken,
		refreshToken,
		s.cfg.Oauth.AccessTokenLifetime,
		tokentypes.Bearer,
	)
	if err != nil {
		return nil, err
	}

	return accessTokenResponse, nil
}
