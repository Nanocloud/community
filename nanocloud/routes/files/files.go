package files

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

var (
	InvalidDownloadToken = errors.New("Invalid Download Token")
)

var kExecutionServer string

func jsonResponse(w http.ResponseWriter, r *http.Request, statusCode int, body hash) {
	w.Header().Set("Content-Type", "application/json")

	bodyBuff, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write([]byte(bodyBuff))
	if err != nil {
		log.Error(err)
	}
}

func sha1Hash(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func oauthError(c *echo.Context, fail *oauth2.OAuthError) error {
	b, err := fail.ToJSON()
	if err != nil {
		return err
	}

	w := c.Response()
	w.WriteHeader(fail.HTTPStatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
	return nil
}

/*
 * Download tokens are needed to download a file without a OAuth access token.
 * As we do a simple HTTP request, we don't want the token to appear in the history
 * or to be shared to a third party user.
 * Download tokens are valid for one hour.
 *
 * A token is build as : {access_token_id}:{SHA1{oauth_access_token:filename:time_stone}}
 * Where:
 *  - access_token_id is the id of a OAuth access token of the file owner
 *  - oauth_access_token is the token associated to access_token_id
 *  - filename is the name of the filename associated to the token
 *  - time_stone is this (($NOW + $NOW % 3600) + 3600) where NOW is a unix timestamp.
 *    It makes a download token valid for current and the next hour.
 *    (is the token is generate at 1:55am then the token is valid from 1:00am to 2:59)
 */
func createDownloadToken(user *users.User, accessToken string, filename string) (string, error) {
	rows, err := db.Query(
		"SELECT id FROM oauth_access_tokens WHERE token = $1::varchar AND user_id = $2::varchar",
		accessToken,
		user.Id,
	)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", errors.New("unable to retreive user from access token")
	}

	var accessTokenId string
	err = rows.Scan(&accessTokenId)
	if err != nil {
		return "", err
	}

	timeStone := time.Now().Unix()
	timeStone = timeStone + (3600 - timeStone%3600) + 3600

	fmt.Printf("time stone = %d\n", timeStone)

	h := sha1Hash(fmt.Sprintf(
		"%s:%s:%d",
		accessToken,
		filename,
		timeStone,
	))

	return accessTokenId + ":" + h, nil
}

/**
 * Return the appropriate user for the download token and the filename.
 */
func checkDownloadToken(token, filename string) (*users.User, error) {
	splt := strings.SplitN(token, ":", 2)
	if len(splt) != 2 {
		return nil, InvalidDownloadToken
	}

	accessTokenId := splt[0]
	rows, err := db.Query(
		"SELECT token, user_id FROM oauth_access_tokens WHERE id = $1::varchar",
		accessTokenId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, InvalidDownloadToken
	}

	timeStone := time.Now().Unix()
	timeStone = timeStone + (3600 - timeStone%3600)

	var accessToken string
	var userId string
	err = rows.Scan(&accessToken, &userId)
	if err != nil {
		return nil, err
	}
	h := accessToken + ":" + filename + ":"

	expected := sha1Hash(h + strconv.FormatInt(timeStone, 10))
	if splt[1] != expected {
		timeStone = timeStone + 3600
		expected = sha1Hash(h + strconv.FormatInt(timeStone, 10))
		if splt[1] != expected {
			return nil, InvalidDownloadToken
		}
	}

	return users.GetUser(userId)
}

func GetDownloadToken(c *echo.Context) error {
	filename := c.Query("filename")
	if len(filename) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			hash{
				"error": "filename not specified",
			},
		)
	}

	accessToken, fail := oauth2.GetAccessToken(c.Request())
	if fail != nil {
		return oauthError(c, fail)
	}

	user := c.Get("user").(*users.User)
	token, err := createDownloadToken(user, accessToken, filename)
	if err != nil {
		return err
	}
	return c.JSON(
		http.StatusOK,
		hash{
			"token": token,
		},
	)
}

func Get(c *echo.Context) error {
	w := c.Response()
	r := c.Request()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cashe-Control", "no-store")
	w.Header().Set("Expires", "Sat, 01 Jan 2000 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")

	path := c.Query("path")

	filename := c.Query("filename")
	if len(filename) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			hash{
				"error": "Invalid Path",
			},
		)
	}

	var user *users.User
	var err error

	downloadToken := c.Query("token")
	if len(downloadToken) > 0 {
		user, err = checkDownloadToken(downloadToken, filename)
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				hash{
					"error": "Invalid Download Token",
				},
			)
		}
	}

	if user == nil {
		u, fail := oauth2.GetUser(w, r)
		if u == nil {
			return errors.New("no authenticated user")
		}

		if fail != nil {
			return oauthError(c, fail)
		}

		user = u.(*users.User)
	}

	filename = strings.Replace(filename, "/", "\\", -1)
	path = filename
	/*
		path := fmt.Sprintf(
			"C:\\Users\\%s\\Desktop\\Nanocloud%s",
			user.Sam,
			filename,
		)
	*/

	resp, err := http.Get("http://" + kExecutionServer + ":9090/files?create=true&path=" + url.QueryEscape(path))
	if err != nil {
		log.Error(err)

		return errors.New("Unable to contact the server")
	}

	if resp.StatusCode == http.StatusNotFound {
		jsonResponse(w, r, http.StatusNotFound, hash{
			"error": "File Not Found",
		})
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unable to retreive the file")
	}

	var contentType string
	rContentType, exists := resp.Header["Content-Type"]
	if !exists || len(rContentType) == 0 || len(rContentType[0]) == 0 {
		contentType = "application/octet-stream"
	} else {
		contentType = rContentType[0]
	}

	var sent int64
	var lastBuffSize int64

	contentLength := resp.ContentLength

	var f string
	splt := strings.Split(path, "\\")
	if len(splt) > 0 {
		f = splt[len(splt)-1]
	} else {
		f = path
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Disposition", "attachment; filename=\""+f+"\"")

	var buff []byte

	for sent < contentLength {
		remaining := contentLength - sent
		if remaining > 4096 {
			remaining = 4096
		}

		if buff == nil || remaining != lastBuffSize {
			buff = make([]byte, remaining)
			lastBuffSize = remaining
		}

		nRead, readErr := resp.Body.Read(buff)

		if nRead > 0 {
			nWrite, writeErr := w.Write(buff[0:nRead])
			sent = sent + int64(nWrite)

			if writeErr != nil {
				log.Errorf("Write error: %s\n", writeErr.Error())
				break
			}
		}

		if readErr != nil && readErr.Error() != "EOF" {
			log.Errorf("Read error: %s\n", readErr.Error())
			break
		}
	}
	return nil
}

func init() {
	executionServers := strings.Split(utils.Env("EXECUTION_SERVERS", ""), ",")
	if len(executionServers) < 1 || len(executionServers[0]) < 1 {
		panic(errors.New("EXECUTION_SERVERS not set"))
	}
	kExecutionServer = executionServers[0]
}
