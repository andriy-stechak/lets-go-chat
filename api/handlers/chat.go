package handlers

import (
	"net/http"

	"github.com/andriystech/lgc/services"
)

func WSConnectHandler(ws *services.WebSocketService, ts *services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		token, ok := q["token"]
		if !ok || len(token) < 1 {
			SendErrorJsonResponse(w, http.StatusBadRequest, "Query parameter 'token' is missing")
			return
		}

		user, err := ts.GetUserByToken(r.Context(), token[0])
		if err != nil {
			SendErrorJsonResponse(w, http.StatusForbidden, err.Error())
			return
		}
		err = ws.NewConnection(w, r, user)
		if err != nil {
			SendErrorJsonResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
