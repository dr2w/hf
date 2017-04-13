package fe

import (
    "fmt"
    "google.golang.org/appengine"
    "google.golang.org/appengine/user"
    "google.golang.org/appengine/log"
    "html/template"
    "net/http"
)

func init() {
    http.HandleFunc("/", handler)
}

type Info struct {
	User string
        Email string
        Logout string
}

var frontEndTemplate = template.Must(template.New("fe").Parse(`
<html>
  <head>
    <title>High Five</title>
  </head>
  <body>
    User: {{.User}}<br>
    Email: {{.Email}}<br>
    <a href="{{.Logout}}">sign out</a>
  </body>
</html>
`))

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-type", "text/html; charset=utf-8")
    ctx := appengine.NewContext(r);
    u := user.Current(ctx);
    if u == nil {
        url, _ := user.LoginURL(ctx, "/")
        fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
        return
    }
    url, _ := user.LogoutURL(ctx, "/")
    i := Info{u.String(), u.Email, url};
    if err := frontEndTemplate.Execute(w, i); err != nil {
        log.Errorf(ctx, err.Error());
        http.Error(w, err.Error(), http.StatusInternalServerError);
    }
}
