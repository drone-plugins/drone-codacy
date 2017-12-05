package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
	"github.com/schrej/godacov/coverage"
	"golang.org/x/tools/cover"
)

const (
	RequestURI = "https://api.codacy.com/2.0/coverage/%s/%s"
)

type (
	Build struct {
		Commit string `json:"commit"`
	}

	Config struct {
		Token    string
		Pattern  string
		Language string
		Debug    bool
	}

	Internal struct {
		Report  []byte
		Matches []string
		Merged  Packages
	}

	Plugin struct {
		Build    Build
		Config   Config
		Internal Internal
	}
)

func (p *Plugin) Exec() error {
	if p.Config.Token == "" {
		return errors.New("you must provide a token")
	}

	if err := p.match(); err != nil {
		return err
	}

	if err := p.merge(); err != nil {
		return err
	}

	if err := p.build(); err != nil {
		return err
	}

	if err := p.submit(); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) match() error {
	log.Printf("searching coverage files: %s", p.Config.Pattern)
	matches, err := zglob.Glob(p.Config.Pattern)

	if err != nil {
		return errors.Wrap(err, "failed to match files")
	}

	if len(matches) > 0 {
		log.Printf("found coverage files: %s", strings.Join(matches, ", "))
	} else {
		log.Printf("no coverage files found")
	}

	p.Internal.Matches = matches

	return nil
}

func (p *Plugin) merge() error {
	if len(p.Internal.Matches) == 0 {
		return nil
	}

	for _, f := range p.Internal.Matches {
		profiles, err := cover.ParseProfiles(f)

		if err != nil {
			return errors.Wrap(err, "failed to parse profile")
		}

		for _, profile := range profiles {
			p.Internal.Merged.Add(profile)
		}
	}

	return nil
}

func (p *Plugin) build() error {
	if len(p.Internal.Matches) == 0 {
		return nil
	}

	if err := os.MkdirAll("/tmp", os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create tempdir")
	}

	tmpfile, err := ioutil.TempFile("", "codacy-")

	if err != nil {
		return errors.Wrap(err, "failed to create tempfile")
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(p.Internal.Merged.Dump()); err != nil {
		return errors.Wrap(err, "failed to write tempfile")
	}

	stat, err := tmpfile.Stat()

	if err != nil {
		return errors.Wrap(err, "failed to retrieve stats")
	}

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "failed to close temp file")
	}

	if stat.Size() > 0 {
		report, err := coverage.GenerateCoverageJSON(tmpfile.Name())

		if err != nil {
			return errors.Wrap(err, "failed to generate report")
		}

		p.Internal.Report = report
	}

	return nil
}

func (p *Plugin) submit() error {
	if len(p.Internal.Matches) == 0 {
		return nil
	}

	if len(p.Internal.Report) == 0 {
		log.Printf("skipping submission of empty report")
		return nil
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(RequestURI, p.Build.Commit, p.Config.Language),
		bytes.NewBuffer(p.Internal.Report),
	)

	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("project_token", p.Config.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return errors.Wrap(err, "failed to process request")
	}

	defer res.Body.Close()

	if res.Status == "200 OK" {
		log.Printf("successfully uploaded coverage report")
	} else {
		response := struct {
			Error string
		}{}

		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return errors.Wrap(err, "failed to parse response")
		}

		return errors.Errorf("failed to submit request: %s", response.Error)
	}

	return nil
}
