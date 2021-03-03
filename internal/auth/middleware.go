package auth

import (
	"context"
	"github.com/sirupsen/logrus"
	"gqlexample/v2/internal/users"
	"gqlexample/v2/pkg/jwt"
	"net/http"
)

type contextKey struct {
	name string
}

var userCtxKey = &contextKey{"athResp"}

type AthResponseWriter struct {
	http.ResponseWriter
	User             *users.User
	UserIDFromCookie string
}

func (w *AthResponseWriter) Write(b []byte) (int, error) {
	//if w.User.Username != w.UserIDFromCookie {
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    string(b),
		HttpOnly: true,
		Path:     "/query",
		Domain:   "localhost",
	})
	return w.ResponseWriter.Write([]byte(""))
}

/*
	try to get jwt from cookie or header
	cookie having priority
 */
func getJwtFromRequest(r * http.Request) (jwt string) {
	header := r.Header.Get("Authorization")
	cookie, err := r.Cookie("Authorization")
	if err != nil && err.Error() != "http: named cookie not present" {
		check(err)
	}
	if header == "" && cookie == nil {
		return ""
	}
	if cookie != nil {
		return cookie.Value
	}
	return header
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentJwt := getJwtFromRequest(r)

			// Allow unauthenticated users in but give access to writer
			if currentJwt == "" {
				authResp := AthResponseWriter{
					w,
					nil,
					"",
				}
				ctx := context.WithValue(r.Context(), userCtxKey, &authResp)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}

			//validate jwt token
			username, err := jwt.ParseToken(currentJwt)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// create user and check if user exists in db
			user := users.User{Username: username}
			id, err := users.GetUserIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.ID = id
			// put it in context
			authResp := AthResponseWriter{
				w,
				&user,
				"",
			}
			ctx := context.WithValue(r.Context(), userCtxKey, &authResp)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	athResp, _ := ctx.Value(userCtxKey).(*AthResponseWriter)
	return athResp.User
}

func SetCookieForContext(ctx context.Context, token string) {
	athResp, _ := ctx.Value(userCtxKey).(*AthResponseWriter)
	if athResp == nil {
		logrus.Info("can't setCookieForContext, athResp is nil")
		return
	}
	_, err := athResp.Write([]byte(token))
	check(err)
	logrus.Info("setting cookie for context")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
