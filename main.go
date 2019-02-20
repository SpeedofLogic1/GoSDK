ackage main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"golang.org/x/oauth2"
)

var (
	azureAuthConfig  *oauth2.Config
	oauthStateString = "pseudo-random"
	// AzureEndpoint is a representation of a AzureEndpoint
	AzureEndpoint = oauth2.Endpoint{
		AuthURL:   "https://login.microsoftonline.com/zzz/oauth2/v2.0/authorize",
		TokenURL:  "https://login.microsoftonline.com/zz/oauth2/v2.0/token",
		AuthStyle: oauth2.AuthStyleInParams,
	}
)

// called automaticly I believe
func init() {
	azureAuthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/client",
		ClientID:     "zzz",
		ClientSecret: "zz",
		Scopes:       []string{"User.Read.All"},
		Endpoint:     AzureEndpoint,
	}
}
func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html><body><a href="/login">Azure Log In</a></body></html>`
	fmt.Fprintf(w, htmlIndex)
}
func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := azureAuthConfig.AuthCodeURL(oauthStateString)
	log.Print(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := azureAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	req.Header.Set("Authorization", token.AccessToken)
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func main() {
	log.Print("Logging in Go!")
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/client", handleCallback)
	http.HandleFunc("/login", handleLogin)
	http.ListenAndServe(":8080", nil)
}
 
