package app

import (
	haikunator "github.com/atrox/haikunatorgo/v2"
	"net/http"
)

const CookieUname = "u_name"

func CookieAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		host := request.Header.Get("host")

		cookieUName, err := request.Cookie(CookieUname)

		var uname string
		if err == http.ErrNoCookie || cookieUName == nil || cookieUName.Value == "" {
			haikunate := haikunator.New()
			haikunate.Delimiter = " "
			haikunate.TokenLength = 0

			uname = Capitalize(haikunate.Haikunate())
		} else {
			uname = cookieUName.Value
		}

		cookieUName = &http.Cookie{
			Name:   CookieUname,
			Value:  uname,
			Domain: host,
			Path:   "/",
			MaxAge: 7 * 24 * 60 * 60,
		}
		http.SetCookie(writer, cookieUName)

		f(writer, request)
	}
}
