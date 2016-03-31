package oauth2

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

const (
	INVALID_GRANT             = "invalid_grant"
	INVALID_REQUEST           = "invalid_request"
	INVALID_CLIENT            = "invalid_client"
	ACCESS_DENIED             = "access_denied"
	UNAUTHORIZED_CLIENT       = "unauthorized_client"
	UNSUPPORTED_RESPONSE_TYPE = "unsupported_response_type"
	INVALID_SCOPE             = "invalid_scope"
	SERVER_ERROR              = "server_error"
	TEMPORARILY_UNAVAILABLE   = "temporarily_unavailable"
	UNSUPPORTED_GRANT_TYPE    = "unsupported_grant_type"
)

type Connector interface {
	GetClient(key, secret string) (interface{}, error)
	GetUserFromAccessToken(accessToken string) (interface{}, error)
	AuthenticateUser(username, password string) (interface{}, error)
	GetAccessToken(interface{}, interface{}, *http.Request) (interface{}, error)
}

var kConnector Connector

type OAuthError struct {
	HTTPStatusCode int
	Err            string
	Description    string
}

func (e *OAuthError) Error() (err string) {
	err = "OAuth Error: " + e.Err

	if len(e.Description) > 0 {
		err += " " + e.Description
	}
	return
}

func (e *OAuthError) ToJSON() (rt []byte, err error) {
	obj := make(map[string]string)

	if len(e.Err) > 0 {
		obj["error"] = e.Err
	} else {
		obj["error"] = "unknown"
	}

	if len(e.Description) > 0 {
		obj["error_description"] = e.Description
	}

	rt, err = json.Marshal(obj)
	return
}

func GetAccessToken(req *http.Request) (accessToken string, err *OAuthError) {
	/*
	 * Check Query String
	 */
	query := req.URL.Query()
	accessTokenQuery, exists := query["access_token"]
	if exists && len(accessTokenQuery[0]) > 0 {
		accessToken = accessTokenQuery[0]
		return
	}

	/*
	 * Check Authorization Header
	 */

	accessToken, err = GetAuthorizationHeaderValue(req, "Bearer")
	return
}

func SetConnector(connector Connector) {
	kConnector = connector
}

func GetUser(res http.ResponseWriter, req *http.Request) (interface{}, *OAuthError) {
	accessToken, err := GetAccessToken(req)
	if err != nil {
		return nil, err
	}

	user, fail := kConnector.GetUserFromAccessToken(accessToken)
	if fail != nil {
		log.Error("[OAuth] Cannot retreive user form access token: " + fail.Error())
		return nil, &OAuthError{500, SERVER_ERROR, "Internal Server Error"}
	}

	if user != nil {
		return user, nil
	}

	return nil, &OAuthError{http.StatusUnauthorized, ACCESS_DENIED, "Invalid access token"}
}

func GetAuthorizationHeaderValue(req *http.Request, authType string) (string, *OAuthError) {
	rawHeader := req.Header.Get("Authorization")
	if len(rawHeader) < 1 {
		return "", &OAuthError{http.StatusUnauthorized, INVALID_REQUEST, "Authorization header is missing"}
	}

	splt := strings.SplitN(rawHeader, " ", 2)
	if len(splt) != 2 {
		return "", &OAuthError{http.StatusUnauthorized, INVALID_REQUEST, "Invalid Authorization header"}
	}

	if splt[0] != authType {
		return "", &OAuthError{http.StatusUnauthorized, INVALID_REQUEST, "Invalid authorization type"}
	}

	token := splt[1]
	return token, nil
}

func clientBasicAuth(req *http.Request) (interface{}, *OAuthError) {
	rawAuthToken, err := GetAuthorizationHeaderValue(req, "Basic")
	if err != nil {
		return nil, err
	}

	log.Debug("[OAuth] token = " + rawAuthToken)

	bToken, fail := base64.StdEncoding.DecodeString(rawAuthToken)
	if fail != nil {
		log.Warn("[Oauth] Unable to parse base64 auth basic string: " + rawAuthToken)
		return nil, &OAuthError{http.StatusUnauthorized, INVALID_REQUEST, "Invalid Authorization header"}
	}

	token := string(bToken)
	splt := strings.SplitN(token, ":", 2)
	if len(splt) != 2 {
		return nil, &OAuthError{http.StatusUnauthorized, INVALID_REQUEST, "Invalid Authorization header"}
	}

	clientKey := splt[0]
	clientSecret := splt[1]

	client, fail := kConnector.GetClient(clientKey, clientSecret)
	if fail != nil {
		log.Error("[Oauth] Unable to retreive client: " + fail.Error())
		return nil, &OAuthError{500, SERVER_ERROR, "Internal Server Error"}
	}

	return client, nil
}

func oauthErrorReply(res http.ResponseWriter, oauthErr OAuthError) error {
	res.Header().Set("Content-Type", "application/json;charset=UTF-8")
	ret, err := oauthErr.ToJSON()
	if err != nil {
		log.Error("[OAuth] Cannot write JSON error: " + err.Error())
		return err
	}

	res.WriteHeader(oauthErr.HTTPStatusCode)
	res.Write(ret)
	return nil
}

func HandleRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Cache-Control", "no-store")
	res.Header().Add("Pragma", "no-cache")

	if req.Method == "POST" && req.URL.Path == "/oauth/token" {

		client, err := clientBasicAuth(req)
		if err != nil {
			oauthErrorReply(res, *err)
			return
		}

		if client == nil {
			oauthErrorReply(res, OAuthError{http.StatusUnauthorized, INVALID_CLIENT, "Invalid OAuth Client Credentials"})
			return
		}

		error := req.ParseForm()

		if error != nil {
			oauthErrorReply(res, OAuthError{http.StatusBadRequest, INVALID_REQUEST, "Unable to parse the request body"})
			return
		}

		// grant_type
		grantType := req.FormValue("grant_type")
		if grantType == "" {
			oauthErrorReply(res, OAuthError{http.StatusBadRequest, INVALID_REQUEST, "grant_type is missing"})
			return
		}

		if grantType != "password" {
			oauthErrorReply(res, OAuthError{http.StatusBadRequest, INVALID_REQUEST, "Invalid grant_type"})
			return
		}

		// username
		username := req.FormValue("username")
		if username == "" {
			oauthErrorReply(res, OAuthError{http.StatusBadRequest, INVALID_REQUEST, "username is missing"})
			return
		}

		// password
		password := req.FormValue("password")
		if password == "" {
			oauthErrorReply(res, OAuthError{http.StatusBadRequest, INVALID_REQUEST, "password is missing"})
			return
		}

		user, fail := kConnector.AuthenticateUser(username, password)

		if fail != nil {
			if fail.Error() == "invalid credentials" || fail.Error() == "user not found" {
				oauthErrorReply(res, OAuthError{http.StatusUnauthorized, ACCESS_DENIED, "Invalid User Credentials"})
				return
			}
			log.Error("[OAuth] Cannot Authenticate User: " + fail.Error())
			oauthErrorReply(res, OAuthError{http.StatusInternalServerError, SERVER_ERROR, "Internal Server Error"})
			return
		}

		accessToken, fail := kConnector.GetAccessToken(user, client, req)
		if fail != nil {
			log.Error("[OAuth] Cannot Get Access Token: " + fail.Error())
			oauthErrorReply(res, OAuthError{http.StatusInternalServerError, SERVER_ERROR, "Internal Server Error"})
			return
		}

		if accessToken == nil {
			oauthErrorReply(res, OAuthError{http.StatusUnauthorized, ACCESS_DENIED, "Access token request denied for the given client"})
			return
		}

		rt, fail := json.Marshal(accessToken)
		if fail != nil {
			log.Error("[OAuth] Unable to serialize access token: " + fail.Error())
			oauthErrorReply(res, OAuthError{http.StatusInternalServerError, SERVER_ERROR, "Internal Server Error"})
			return
		}

		res.Header().Set("Content-Type", "application/json;charset=UTF-8")
		res.Write(rt)

		return
	}

	oauthErrorReply(res, OAuthError{http.StatusNotFound, INVALID_REQUEST, "Invalid Endpoint"})
}

type dummyConnector struct{}

func (c dummyConnector) GetClient(key, secret string) (interface{}, error) {
	return nil, errors.New("GetClient is not implemented")
}

func (c dummyConnector) GetUserFromAccessToken(accessToken string) (interface{}, error) {
	return nil, errors.New("GetUserFromAccessToken is not implemented")
}

func (c dummyConnector) AuthenticateUser(username, password string) (interface{}, error) {
	return nil, errors.New("AuthenticateUser is not implemented")
}

func (c dummyConnector) GetAccessToken(user, client interface{}, req *http.Request) (interface{}, error) {
	return nil, errors.New("GetAccessToken is not implemented")
}

func init() {
	SetConnector(dummyConnector{})
}
