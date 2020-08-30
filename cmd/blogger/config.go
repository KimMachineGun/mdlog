package main

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configPath            = "./blogger.yaml"
	defaultCredentialPath = "./credentials.json"
	defaultCachePath      = "./.blogger.token"
	defaultPostsPath      = "./posts"
)

type config struct {
	BloggerURL     string `yaml:"blogger_url"`
	CredentialPath string `yaml:"credential_path"`
	CachePath      string `yaml:"cache_path"`
	PostsPath      string `yaml:"posts_path"`
}

func (c *config) validate() error {
	if c.BloggerURL == "" {
		return errors.New("blogger_url is required")
	}
	if c.CredentialPath == "" {
		return errors.New("credential_path is required")
	}
	if c.CachePath == "" {
		return errors.New("cache_path is required")
	}
	if c.PostsPath == "" {
		return errors.New("posts_path is required")
	}
	return nil
}

func getConfig() (*config, error) {
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("not initialized directory")
		}
		return nil, err
	}

	c := config{
		CredentialPath: defaultCredentialPath,
		CachePath:      defaultCachePath,
		PostsPath:      defaultPostsPath,
	}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, err
	}

	return &c, nil
}
