package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	_ "github.com/lalyos/httptrace"
	gitter "github.com/sromku/go-gitter"
	"golang.org/x/oauth2"
)

func statusHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "ok\n")
}

func indexHandler(w http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, "session-name")
	t, ok := session.Values["token"]
	if !ok {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		//fmt.Printf("Visit the URL for the auth dialog: %v", url)
		fmt.Fprintf(w, `<a href="%s">login to gitter </a>`, url)
		return
	}
	token := t.(string)
	log.Println("---> token from session: %s", token)
	api := gitter.New(token)
	user, err := api.GetUser()
	if err != nil {
		log.Println("[ERROR] couldnt get gitter user: ", err)
	}
	fmt.Fprintf(w, `<h2>Welcome: %s</h2><img src="%s" />`, user.Username, user.AvatarURLSmall)
}
func setupHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
	 <ul>
	 <li><a href="https://developer.gitter.im/apps/new">Create a Gitter APP </a>
	 with redirect url <span style='background-color: lightblue;'>http://%v/login/callback</span> </li>
	 <li><a href="https://gitter.im/#createcommunity">create a Gitter Room</a>
	 </ul>

	 <form method="POST">
	  <br>oauth key:<input name="key" />
	  <br>oauth secret:<input name="secret" />
	  <br>room name:<input name="room" />
	  <input type="submit" value="setup" />
	 </form>
	`, req.Host)
}

func run(args ...string) string {
	cmd := exec.Command(args[0], args[1:]...)
	stde := new(strings.Builder)

	cmd.Stderr = stde
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("ERROR:", err, stde.String())
	}
	return string(out)
}

func authCalbackHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	code := req.Form.Get("code")

	if code == "" {
		log.Println("empty 'code' in oauth calback")
		return
	}
	ctx := context.Background()
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Println("[ERROR] oauth2 code exchange failed:", err)
	}
	log.Println("token:", tok)
	session, _ := store.Get(req, "session-name")
	log.Println("---> save token in session")
	session.Values["token"] = tok
	session.Save(req, w)

	api := gitter.New(tok.AccessToken)
	user, err := api.GetUser()
	if err != nil {
		log.Println("[ERROR] couldnt get gitter user: ", err)
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<h2>Welcome: %s</h2><img src="%s" />`, user.Username, user.AvatarURLSmall)

	out := run("./getsession.sh", user.Username, tok.AccessToken, os.Getenv("GITTER_ROOM_NAME"))
	log.Println("[OUTPUT]", out)
	fmt.Fprintf(w, "<h2>script output</h2> %s", out)
}

var conf *oauth2.Config
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	oauthID := os.Getenv("GITTER_OAUTH_KEY")
	oauthSecret := os.Getenv("GITTER_OAUTH_SECRET")
	domain := os.Getenv("DOMAIN")

	if domain == "" {
		log.Println("[WARNING] please set env var: DOMAIN (ex: gitter.k8z.eu)")
	} else {
		log.Println("DOMAIN:", domain)
	}
	if oauthID == "" || oauthSecret == "" {
		log.Println("[WARNING] please set env vars: GITTER_OAUTH_KEY, GITTER_OAUTH_SECRET")
	}
	if os.Getenv("GITTER_ROOM_NAME") == "" {
		log.Println("[WARNING] please set env var: GITTER_ROOM_NAME")
	} else {
		log.Println("GITTER_ROOM_NAME:", os.Getenv("GITTER_ROOM_NAME"))
	}

	fmt.Println("starting server on port:", port)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/setup", setupHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/login/callback", authCalbackHandler)

	conf = &oauth2.Config{
		ClientID:     oauthID,
		ClientSecret: oauthSecret,
		Scopes:       []string{""},
		RedirectURL:  fmt.Sprintf("http://%s/login/callback", domain),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://gitter.im/login/oauth/authorize",
			TokenURL: "https://gitter.im/login/oauth/token",
		},
	}

	log.Fatal(http.ListenAndServe(":"+port, gcontext.ClearHandler(http.DefaultServeMux)))
}
