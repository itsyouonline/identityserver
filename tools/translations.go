package tools

import (
	"bytes"
	"encoding/json"
	"html/template"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/templates/packaged"
)

// translationTemplates is a map representing the translations as loaded from the
// translation files
type translationTemplates map[string]string

// TranslationValues is a map, with the keys beign the requested translations,
// and the values being a struct with the values to use in the template of said
// translation
type TranslationValues map[string]interface{}

// Translations map the requested keys to their final rendered values
type Translations map[string]string

// loadTranslationTemplates tries to load a translation file
func loadTranslationTemplates(rawKey string) (tt translationTemplates, err error) {
	langKey := parseLangKey(rawKey)
	assetName := "i18n/" + langKey + ".json"
	translationFile, err := templates.Asset(assetName)
	// translation file doesn't exis, or there is an error loading it
	if err != nil {
		// try and use the default translations
		translationFile, err = templates.Asset(defaultTranslations)
		if err != nil {
			log.Error("Error while loading translations: ", err)
		}
	}
	err = json.NewDecoder(bytes.NewReader(translationFile)).Decode(&tt)
	return
}

// Parse translations loads the translation templates and renders them with the
// provided values
func ParseTranslations(rawKey string, tv TranslationValues) (translations Translations, err error) {
	// Load the translation templates from the file
	tt, err := loadTranslationTemplates(rawKey)
	if err != nil {
		log.Error("Error while parsing translations - Failed to load translations: ", err)
		return
	}
	translations = make(Translations)
	templateEngine := template.New("translations")
	buf := new(bytes.Buffer)
	// Iterate over the keys and values provided to render
	for key, value := range tv {
		// If a translation key is provided that does not exist in the file consider it an error
		template, exists := tt[key]
		if !exists {
			log.Error("Error while parsing translations - Trying to render an unexisting key: ", err)
			return
		}
		// If no translation values are provided, store the raw template (could be that there
		// are no values required for this template).
		if value == nil {
			translations[key] = template
			continue
		}
		// Make sure the buffer is empty
		buf.Reset()
		// Parse the template string in the template engine and render it in the buffer
		templateEngine.Parse(template)
		err = templateEngine.Execute(buf, value)
		if err != nil {
			log.Error("Error while parsing translations - Failed to render template: ", err)
			return
		}
		translations[key] = buf.String()
	}
	return
}

// parseLangKey return the first 2 characters of a string in lowercase.
// If the string is empty or has only 1 character, and empty string is returned
func parseLangKey(rawKey string) string {
	if len(rawKey) < 2 {
		return ""
	}
	return strings.ToLower(string(rawKey[0:2]))
}
