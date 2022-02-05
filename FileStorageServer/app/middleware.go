package app

import (
	"FileStorageServer/token"

	"github.com/sirupsen/logrus"
	"strings"
	"net/http"
	"time"
)

func (app *App) LogMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

                start := time.Now()

                next.ServeHTTP(w, r)

                app.MyLogger.WithFields(logrus.Fields{
                        "method":      r.Method,
                        "remote_addr": r.RemoteAddr,
                        "work_time":   time.Since(start),
                }).Info(r.URL.Path)
        })
}

func (app *App) CheckAuthMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

                HeaderAuth := strings.Split(r.Header.Get("Authorization"), " ")

                flag, err := token.CheckToken(HeaderAuth , app.Config.KeyToken)
                if err != nil && flag == false {
                        http.Error(w, "INCORRECT TOKEN SIGNATURE", 400)
                        return //error handling
                }
                if err == nil && flag == false {
                        http.Error(w, "THE TOKEN IS NOT VALID", 401)
                        return //error handling
                }

                next.ServeHTTP(w, r)
        })
}
