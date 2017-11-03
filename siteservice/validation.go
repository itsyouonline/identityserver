package siteservice

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/tools"
	"github.com/itsyouonline/identityserver/validation"
)

//PhonenumberValidation is the page that is linked to in the SMS for phonenumbervalidation and is thus accessed on the mobile phone
func (service *Service) PhonenumberValidation(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := request.Form
	key := values.Get("k")
	smscode := values.Get("c")
	langKey := values.Get("l")

	translationValues := tools.TranslationValues{
		"invalidlink":      nil,
		"error":            nil,
		"smsconfirmed":     nil,
		"return_to_window": nil,
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	err = service.phonenumberValidationService.ConfirmValidation(request, key, smscode)
	if err == validation.ErrInvalidCode || err == validation.ErrInvalidOrExpiredKey {
		service.renderSMSConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}
	if err != nil {
		log.Error(err)
		service.renderSMSConfirmationPage(w, request, translations["error"], "")
		return
	}

	service.renderSMSConfirmationPage(w, request, translations["smsconfirmed"], translations["return_to_window"])
}

func (service *Service) EmailValidation(w http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		log.Debug(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	values := request.Form
	key := values.Get("k")
	smscode := values.Get("c")
	langKey := values.Get("l")

	translationValues := tools.TranslationValues{
		"invalidlink":      nil,
		"error":            nil,
		"emailconfirmed":   nil,
		"return_to_window": nil,
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	err = service.emailaddressValidationService.ConfirmValidation(request, key, smscode)
	if err == validation.ErrInvalidCode || err == validation.ErrInvalidOrExpiredKey {
		service.renderEmailConfirmationPage(w, request, translations["invalidlink"], "")
		return
	}
	if err != nil {
		log.Error(err)
		service.renderEmailConfirmationPage(w, request, translations["error"], "")
		return
	}

	service.renderEmailConfirmationPage(w, request, translations["emailconfirmed"], translations["return_to_window"])
}
