package oauth

import (
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	oauth2.HandleRequest(w, r)
}
