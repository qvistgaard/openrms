package announcer

import (
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"math/rand"
	"strings"
	"text/template"
)

type Announcement interface {
	Get() string
}

type StringTemplateAnnouncement struct {
	Paragraph string
	Data      any
}

func (s *StringTemplateAnnouncement) Get() string {
	return processTemplate(s.Paragraph, s.Data)
}

type ReadFileTemplateAnnouncement struct {
	Fs       fs.ReadFileFS
	Filename string
	Random   bool
	Data     any
}

func (f *ReadFileTemplateAnnouncement) Get() string {
	data, err := f.Fs.ReadFile(f.Filename)
	if err != nil {
		log.Error(errors.WithMessage(err, "Unable to read announcement file: "+f.Filename))
		return ""
	}

	if f.Random {
		split := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
		line := rand.Intn(len(split) - 1)
		return processTemplate(split[line], f.Data)
	}
	return processTemplate(string(data), f.Data)
}

func processTemplate(paragraph string, data any) string {
	if data != nil {
		tmpl, err := template.New("tmpl").Parse(paragraph)

		if err != nil {
			log.Error(errors.WithMessage(err, "Unable to parse announcement template: \""+paragraph+"\""))
		}
		var tmplBytes bytes.Buffer

		err = tmpl.Execute(&tmplBytes, data)
		if err != nil {
			log.Error(errors.WithMessage(err, "Unable to process announcement template: \""+paragraph+"\""))

		}
		return tmplBytes.String()
	}
	return paragraph
}
