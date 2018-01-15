package validation

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
	"github.com/itsyouonline/identityserver/tools"
)

//SMSService is the interface an sms communication channel should have to be used by the IYOPhonenumberValidationService
type SMSService interface {
	Send(phonenumber string, message string) (err error)
}

//IYOPhonenumberValidationService is the itsyou.online implementation of a PhonenumberValidationService
type IYOPhonenumberValidationService struct {
	SMSService SMSService
}

//RequestValidation validates the phonenumber by sending an SMS
func (service *IYOPhonenumberValidationService) RequestValidation(request *http.Request, username string, phonenumber user.Phonenumber, confirmationurl string, langKey string) (key string, err error) {
	valMngr := validation.NewManager(request)
	info, err := valMngr.NewPhonenumberValidationInformation(username, phonenumber)
	if err != nil {
		return
	}
	err = valMngr.SavePhonenumberValidationInformation(info)
	if err != nil {
		return
	}

	translationValues := tools.TranslationValues{
		"smsconfirmation": struct {
			Code string
			Link string
		}{
			Code: info.SMSCode,
			Link: fmt.Sprintf("%s?c=%s&k=%s&l=%s", confirmationurl, info.SMSCode, url.QueryEscape(info.Key), langKey),
		},
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	go service.SMSService.Send(phonenumber.Phonenumber, translations["smsconfirmation"])
	key = info.Key
	return
}

//ExpireValidation removes a pending validation
func (service *IYOPhonenumberValidationService) ExpireValidation(request *http.Request, key string) (err error) {
	if key == "" {
		return
	}
	valMngr := validation.NewManager(request)
	err = valMngr.RemovePhonenumberValidationInformation(key)
	return
}

func (service *IYOPhonenumberValidationService) getPhonenumberValidationInformation(request *http.Request, key string) (info *validation.PhonenumberValidationInformation, err error) {
	if key == "" {
		return
	}
	valMngr := validation.NewManager(request)
	info, err = valMngr.GetByKeyPhonenumberValidationInformation(key)
	return
}

//IsConfirmed checks wether a validation request is already confirmed
func (service *IYOPhonenumberValidationService) IsConfirmed(request *http.Request, key string) (confirmed bool, err error) {
	info, err := service.getPhonenumberValidationInformation(request, key)
	if err != nil {
		return
	}
	if info == nil {
		err = ErrInvalidOrExpiredKey
		return
	}
	confirmed = info.Confirmed
	return
}

//ConfirmValidation checks if the supplied code matches the username and key
func (service *IYOPhonenumberValidationService) ConfirmValidation(request *http.Request, key, code string) (err error) {
	info, err := service.getPhonenumberValidationInformation(request, key)
	if err != nil {
		return
	}
	if info == nil {
		err = ErrInvalidOrExpiredKey
		return
	}
	if info.SMSCode != code {
		err = ErrInvalidCode
		return
	}
	valMngr := validation.NewManager(request)
	p := valMngr.NewValidatedPhonenumber(info.Username, info.Phonenumber)
	err = valMngr.SaveValidatedPhonenumber(p)
	if err != nil {
		return
	}
	err = valMngr.UpdatePhonenumberValidationInformation(key, true)
	if err != nil {
		return
	}
	return
}

// ConfirmRegistrationValidation confirms a validation in the registartion flow. It does not add an entry in the validated
// phone numbers collection
func (service *IYOPhonenumberValidationService) ConfirmRegistrationValidation(r *http.Request, key, code string) (err error) {
	info, err := service.getPhonenumberValidationInformation(r, key)
	if err != nil {
		return
	}
	if info == nil {
		err = ErrInvalidOrExpiredKey
		return
	}
	if info.SMSCode != code {
		err = ErrInvalidCode
		return
	}
	return validation.NewManager(r).UpdatePhonenumberValidationInformation(key, true)
}

//SendOrganizationInviteSms Sends an organization invite SMS
func (service *IYOPhonenumberValidationService) SendOrganizationInviteSms(request *http.Request, invite *invitations.JoinOrganizationInvitation) (err error) {
	link := fmt.Sprintf(invitations.InviteURL, request.Host, url.QueryEscape(invite.Code))
	// todo: perhaps this should be shorter but that might be confusing for the end user
	message := fmt.Sprintf("You have been invited to the %s organization on It's You Online. Click the following link to accept it. %s", invite.Organization, link)
	go service.SMSService.Send(invite.PhoneNumber, message)
	return
}
