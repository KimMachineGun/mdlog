package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bl "google.golang.org/api/blogger/v3"
)

func getTokenSourceWithCredentialFile(ctx context.Context, credentialFile string, cacheFile string) (oauth2.TokenSource, error) {
	b, err := ioutil.ReadFile(credentialFile)
	if err != nil {
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, bl.BloggerScope)
	if err != nil {
		return nil, err
	}

	token, err := getTokenFromCache(cacheFile)
	if err != nil {
		return nil, err
	} else if token != nil {
		log.Println("get token from cache")

		return config.TokenSource(ctx, token), nil
	}

	authURL := config.AuthCodeURL("blogger-state", oauth2.AccessTypeOffline)
	fmt.Printf("URL: %s\nCode: ", authURL)

	var code string
	fmt.Scanln(&code)

	token, err = config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	log.Println("get token from web")

	err = writeTokenToCache(cacheFile, token)
	if err != nil {
		return nil, err
	}

	return config.TokenSource(ctx, token), nil
}

func getTokenFromCache(cacheFile string) (*oauth2.Token, error) {
	b, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(b, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func writeTokenToCache(cacheFile string, token *oauth2.Token) error {
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cacheFile, b, 0644)
}
