/*
Joyrex YSM - Manager for Youtube Subscriptions
Copyright (C) 2025 Eric Stacey <ejstacey@joyrex.net>

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

*/

package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"time"

	gap "github.com/muesli/go-app-paths"
	"golang.org/x/oauth2"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func ConnectYoutube(badFile bool) *youtube.Service {
	ctx := context.Background()

	scope := []string{youtube.YoutubeReadonlyScope}

	config := &oauth2.Config{
		ClientID:     "598298615700-uvl0hgqplp15r7scibsn6mcjv2vc214v.apps.googleusercontent.com",
		ClientSecret: "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",
		Endpoint:     google.Endpoint,
		Scopes:       scope,
	}

	var client = getClient(ctx, config, badFile)

	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	HandleError(err, "Error creating YouTube client")

	return service
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config, badFile bool) *http.Client {
	cacheFile, err := tokenCacheFile()
	HandleError(err, "Error setting up token cache file")
	log.Printf("Youtube token cache file location: %s\n", cacheFile)

	var tok *oauth2.Token
	err = nil
	if !badFile {
		tok, err = tokenFromFile(cacheFile)
	} else {
		fmt.Print("Cached authentication credential file seems invalid. Re-authing.\n")
	}
	if err != nil || badFile {
		tok = tokenFromWeb(ctx, config)
		saveToken(cacheFile, tok)
	}

	return config.Client(ctx, tok)
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized. You may now close this window.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	go openURL(authURL)
	log.Printf("Check your browser. It should have opened the following page (use this link if it hasn't): %s", authURL)
	code := <-ch
	// log.Printf("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return token
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	credScope := gap.NewVendorScope(gap.User, "ysm", "credentials")
	credDirs, err := credScope.DataDirs()
	HandleError(err, "Could not determine user config path for youtube credentials!")
	tokenCacheDir := credDirs[0]
	os.MkdirAll(tokenCacheDir, 0700)
	return credScope.ConfigPath("ysm-youtube-creds.json")
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	HandleError(err, "Unable to cache oauth token")
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func openURL(url string) {
	try := []string{"xdg-open", "google-chrome", "open"}
	for _, bin := range try {
		err := exec.Command(bin, url).Run()
		if err == nil {
			return
		}
	}
	log.Printf("Error opening URL in browser.")
}
