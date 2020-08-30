package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/KimMachineGun/mdlog/internal/blogger"
	"github.com/KimMachineGun/mdlog/internal/mdlog"

	"gopkg.in/yaml.v3"
)

func run(args []string) error {
	if len(args) < 2 {
		return errors.New("sub-command is not provided")
	}

	switch args[1] {
	case "init":
		return commandInit(args[2:])
	case "create":
		return commandCreate(args[2:])
	case "sync":
		return commandSync(args[2:])
	}

	return errors.New("unknown command")
}

func commandInit(args []string) error {
	command := flag.NewFlagSet("init", flag.ExitOnError)

	bloggerURL := command.String("url", "", "your blogger url")
	credentialPath := command.String("credential", defaultCredentialPath, "credential.json for google oauth2")
	cachePath := command.String("cache", defaultCachePath, "destination of access token cache file")
	postsPath := command.String("posts", defaultPostsPath, "directory for your local posts")

	err := command.Parse(args)
	if err != nil {
		return err
	}

	c := config{
		BloggerURL:     *bloggerURL,
		CredentialPath: *credentialPath,
		CachePath:      *cachePath,
		PostsPath:      *postsPath,
	}
	err = c.validate()
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("cannot marshal config to .yaml file: %v", err)
	}

	err = ioutil.WriteFile(configPath, b, 0644)
	if err != nil {
		return fmt.Errorf("cannot write config file: %v", err)
	}

	err = os.MkdirAll(c.PostsPath, 0755)
	if err != nil {
		return fmt.Errorf("cannot make local posts directory: %v", err)
	}

	return nil
}

func commandCreate(args []string) error {
	command := flag.NewFlagSet("create", flag.ExitOnError)

	fileName := command.String("name", "", "file name of new post")

	err := command.Parse(args)
	if err != nil {
		return err
	}

	if *fileName == "" {
		return errors.New("filepath is required")
	}

	if ext := filepath.Ext(*fileName); ext != ".md" {
		return fmt.Errorf("filepath should have .md extension: %s", ext)
	}

	c, err := getConfig()
	if err != nil {
		return fmt.Errorf("cannot get config: %v", err)
	}

	path := filepath.Join(c.PostsPath, *fileName)

	_, err = os.Stat(path)
	if !os.IsNotExist(err) {
		return fmt.Errorf("file already exists: %v", err)
	}

	ctx := context.Background()
	ts, err := getTokenSourceWithCredentialFile(ctx, c.CredentialPath, c.CachePath)
	if err != nil {
		return fmt.Errorf("cannot get token: %v", err)
	}

	b, err := blogger.New(ctx, c.BloggerURL, ts)
	if err != nil {
		return err
	}

	s := mdlog.NewService(b)

	post, err := s.Create()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintln("---"))
	buf.WriteString(fmt.Sprintf("id: %q\n", post.ID))
	buf.WriteString(fmt.Sprintf("title: %q\n", post.Title))
	buf.WriteString(fmt.Sprintf("tags: %v\n", post.Tags))
	buf.WriteString(fmt.Sprintf("status: %d\n", post.Status))
	buf.WriteString(fmt.Sprintln("---"))

	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("cannot write post file: %v", err)
	}

	return nil
}

func commandSync(args []string) error {
	command := flag.NewFlagSet("sync", flag.ExitOnError)

	err := command.Parse(args)

	if err != nil {
		return err
	}

	c, err := getConfig()
	if err != nil {
		return fmt.Errorf("cannot get config: %v", err)

	}

	ctx := context.Background()
	ts, err := getTokenSourceWithCredentialFile(ctx, c.CredentialPath, c.CachePath)
	if err != nil {
		return fmt.Errorf("cannot get token: %v", err)
	}

	b, err := blogger.New(ctx, c.BloggerURL, ts)
	if err != nil {
		return err
	}

	s := mdlog.NewService(b)

	var localPosts []*mdlog.Post
	err = filepath.Walk(c.PostsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(info.Name()) != ".md" {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		p, err := mdlog.ParsePost(b)
		if err != nil {
			return err
		}

		localPosts = append(localPosts, p)

		return nil
	})
	if err != nil {
		return fmt.Errorf("cannot gather local posts: %v", err)
	}

	plan, err := s.Sync(localPosts)
	if err != nil {
		return err
	}

	for _, localPost := range plan {
		log.Printf("updated: %s [%s]", localPost.Title, localPost.ID)
	}

	return nil
}
