package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/ArtisanCloud/PowerSocialite/v4/contracts"
	"github.com/ArtisanCloud/PowerSocialite/v4/models"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

func StartAuthHandler(p contracts.IProvider, needRedirect bool, res http.ResponseWriter, req *http.Request) (contracts.ISession, string, error) {
	iSession, authUrl, err := GetAuthURL(p, res, req)
	//fmt.Println(authUrl)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(res, err)
		return nil, authUrl, err
	}
	if needRedirect {
		http.Redirect(res, req, authUrl, http.StatusTemporaryRedirect)
	}
	return iSession, authUrl, err
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(req *http.Request) string {
	state := req.URL.Query().Get("state")
	if len(state) > 0 {
		return state
	}

	// If a state query param is not passed in, generate a random
	// base64-encoded nonce so that the state on the auth URL
	// is unguessable, preventing CSRF attacks, as described in
	//
	// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("gothic: source of randomness unavailable: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(req *http.Request) string {
	params := req.URL.Query()
	if params.Encode() == "" && req.Method == http.MethodPost {
		return req.FormValue("state")
	}
	return params.Get("state")
}

func GetAuthURL(p contracts.IProvider, res http.ResponseWriter, req *http.Request) (contracts.ISession, string, error) {

	iSess, err := p.StartAuth(SetState(req))
	if err != nil {
		return nil, "", err
	}

	authUrl, err := iSess.GetAuthURL()
	if err != nil {
		return nil, "", err
	}

	return iSess, authUrl, err
}

var CompleteUserAuth = func(p contracts.IProvider, sess contracts.ISession, res http.ResponseWriter, req *http.Request) (*models.User, error) {
	var err error
	err = ValidateState(req, sess)
	if err != nil {
		return nil, err
	}

	user, err := p.UserFromSession(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	params := req.URL.Query()
	if params.Encode() == "" && req.Method == "POST" {
		err = req.ParseForm()
		if err != nil {
			return nil, err
		}
		params = req.Form
	}

	// get new token and retry fetch
	_, err = sess.Authorize(p, params)
	if err != nil {
		return nil, err
	}

	gu, err := p.UserFromSession(sess)
	return gu, err
}

func ValidateState(req *http.Request, sess contracts.ISession) error {
	rawAuthURL, err := sess.GetAuthURL()
	if err != nil {
		return err
	}

	authURL, err := url.Parse(rawAuthURL)
	if err != nil {
		return err
	}

	reqState := GetState(req)

	originalState := authURL.Query().Get("state")
	if originalState != "" && (originalState != reqState) {
		return errors.New("state token mismatch")
	}
	return nil
}
