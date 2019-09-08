package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Luzifer/rconfig/v2"
	"github.com/bogem/id3v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	cfg = struct {
		Album          string `flag:"album" default:"" description:"Set album tag"`
		Artist         string `flag:"artist" default:"" description:"Set artist tag"`
		File           string `flag:"file,f" default:"" description:"File to read / write" validate:"nonzero"`
		LogLevel       string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		Title          string `flag:"title" default:"" description:"Set title tag"`
		VersionAndExit bool   `flag:"version" default:"false" description:"Prints current version and exits"`
		Year           string `flag:"year" default:"" description:"Set year tag"`
	}{}

	version = "dev"
)

func init() {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("id3patch %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	var (
		logger     = log.WithField("file", cfg.File)
		needsWrite bool
	)

	tag, err := id3v2.Open(cfg.File, id3v2.Options{Parse: true})
	switch err {

	case nil:
		// Everything fine

	case id3v2.ErrUnsupportedVersion:
		logger.Warn("No supported ID3v2 tags found")

	default:
		logger.WithError(err).Fatal("Unable to open file")

	}
	defer tag.Close()

	logger.WithFields(log.Fields{
		"album":       tag.Album(),
		"artist":      tag.Artist(),
		"title":       tag.Title(),
		"tag_version": tag.Version(),
		"year":        tag.Year(),
	}).Info("File opened successfully")

	needsWrite = modTag(cfg.Album, tag.Album, tag.SetAlbum, needsWrite)
	needsWrite = modTag(cfg.Artist, tag.Artist, tag.SetArtist, needsWrite)
	needsWrite = modTag(cfg.Title, tag.Title, tag.SetTitle, needsWrite)
	needsWrite = modTag(cfg.Year, tag.Year, tag.SetYear, needsWrite)

	if !needsWrite {
		logger.Info("No tags changed, no write needed")
		return
	}

	if err := save(tag); err != nil {
		logger.WithError(err).Fatal("Unable to save tags")
	}

	logger.Info("Tags written successfully")
}

func modTag(content string, contentFunc func() string, setFunc func(string), needsWrite bool) bool {
	if content == "" || content == contentFunc() {
		return needsWrite
	}

	setFunc(content)
	return true
}

func save(tag *id3v2.Tag) error {
	var err = tag.Save()
	switch err {

	case id3v2.ErrNoFile:
		// We need to do the save ourselves

	default:
		return err

	}

	// Library does not supporting initial tag addition, assemble the
	// newly added tag ourselves...
	var buf = new(bytes.Buffer)

	oStat, err := os.Stat(cfg.File)
	if err != nil {
		return errors.Wrap(err, "Unable to determine original file stats")
	}

	oFile, err := os.Open(cfg.File)
	if err != nil {
		return errors.Wrap(err, "Unable to open original file")
	}

	if _, err = tag.WriteTo(buf); err != nil {
		return errors.Wrap(err, "Unable to write tag header to buffer")
	}

	if _, err = io.Copy(buf, oFile); err != nil {
		return errors.Wrap(err, "Unable to copy original file contents to buffer")
	}

	if err = oFile.Close(); err != nil {
		return errors.Wrap(err, "Unable to close original file")
	}

	return errors.Wrap(ioutil.WriteFile(cfg.File, buf.Bytes(), oStat.Mode()), "Unable to write file contents")
}
