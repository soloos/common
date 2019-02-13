package tinyiron

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	TitlePrefix        string `json:"TitlePrefix"`
	RunMode            string `json:"RunMode"`
	ServeType          string `json:"ServeType"`
	ListenStr          string `json:"ListenStr"`
	LogPath            string `json:"LogPath"`
	AccessWhiteListStr string `json:"AccessWhiteList"`

	SiteViewDir              string `json:"SiteViewDir"`
	SiteStaticBaseUrl        string `json:"SiteStaticBaseUrl"`
	SiteStaticBasePath       string `json:"SiteStaticBasePath"`
	SiteStaticUploadBaseUrl  string `json:"SiteStaticUploadBaseUrl"`
	SiteStaticUploadBasePath string `json:"SiteStaticUploadBasePath"`

	BaseDir           string   `json:"-"`
	AccessWhiteList   []string `json:"-"`
	Log               *os.File `json:"-"`
	IsTMPLAutoRefresh bool     `json:"-"`
}

func (p *Server) LoadOptions(options Options) error {
	var err error
	if err = p.SanitizeOptions(&options); err != nil {
		return err
	}
	p.Options = options
	return nil
}

func (p *Server) LoadOptionsFile(optionsFilepath string) error {
	var (
		err     error
		content []byte
		options Options
	)

	if content, err = ioutil.ReadFile(optionsFilepath); err != nil {
		return err
	}

	if err = json.Unmarshal(content, &options); err != nil {
		return err
	}

	if err = p.LoadOptions(options); err != nil {
		return err
	}

	return nil
}

func (p *Server) SanitizeOptions(options *Options) error {
	var err error

	if options.BaseDir, err = os.Getwd(); err != nil {
		return err
	}

	if options.SiteViewDir == "" {
		options.SiteViewDir = filepath.Join(options.BaseDir, "view")
	}

	switch options.RunMode {
	case "dev", "test", "proc":
		break
	default:
		options.RunMode = "dev"
	}
	switch options.RunMode {
	case "dev":
		options.IsTMPLAutoRefresh = true
	case "test":
		options.IsTMPLAutoRefresh = true
	case "proc":
		options.IsTMPLAutoRefresh = false
	}

	switch options.ServeType {
	case "fcgi", "server":
		break
	default:
		options.ServeType = "server"
	}

	if options.LogPath != "" {
		if nil != options.Log {
			options.Log.Close()
		}
		options.Log, err = os.OpenFile(options.LogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_SYNC, 0755)
		if err != nil {
			return err
		}
		log.SetOutput(options.Log)
	}

	options.AccessWhiteList = nil
	if options.AccessWhiteListStr != "" {
		options.AccessWhiteList = strings.Split(options.AccessWhiteListStr, ",")
		for i := range options.AccessWhiteList {
			options.AccessWhiteList[i] = strings.TrimSpace(options.AccessWhiteList[i])
		}
	}

	return nil
}
