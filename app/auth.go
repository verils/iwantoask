package app

import (
	haikunator "github.com/atrox/haikunatorgo/v2"
	"net/http"
	"strings"
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

			uname = capitalize(haikunate.Haikunate())
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

func capitalize(s string) string {
	split := strings.Split(s, " ")
	for i, str := range split {
		capitalized := strings.ToUpper(str[0:1]) + str[1:]
		split[i] = capitalized
	}
	return strings.Join(split, " ")
}
