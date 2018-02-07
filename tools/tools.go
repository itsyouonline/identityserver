package tools

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"html/template"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/templates/packaged"
)

const (
	defaultTranslations = "i18n/en.json"
)

func GenerateRandomString() (randomString string, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	// Use URLencoding to avoid '/' characters. The generated string it then safe to use
	// in URLs
	randomString = base64.URLEncoding.EncodeToString(b)
	return
}

func RenderTemplate(templateName string, data interface{}) (message string, err error) {
	log.Debug("Email template: ", templateName)
	log.Debug("Email Data: ", data)
	htmlData, err := templates.Asset(templateName)
	if err != nil {
		log.Error("Could not get email asset: ", err)
		return
	}
	templateEngine := template.New("template")
	templateEngine, err = templateEngine.Parse(string(htmlData))
	if err != nil {
		log.Error("Could parse template: ", err)
		return
	}
	buf := new(bytes.Buffer)
	templateEngine.Execute(buf, data)
	message = buf.String()
	return
}
