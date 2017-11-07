package validation

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
	"github.com/itsyouonline/identityserver/tools"
)

const (
	emailWithButtonTemplateName = "emailwithbutton.html"
)

type EmailWithButtonTemplateParams struct {
	UrlCaption string
	Url        string
	Username   string
	Title      string
	Text       string
	ButtonText string
	Reason     string
	LogoUrl    string
}

//EmailService is the interface for an email communication channel, should be used by the IYOEmailAddressValidationService
type EmailService interface {
	Send(recipients []string, subject string, message string) (err error)
}

//IYOEmailAddressValidationService is the itsyou.online implementation of a EmailAddressValidationService
type IYOEmailAddressValidationService struct {
	EmailService EmailService
}

//RequestValidation validates the email address by sending an email
func (service *IYOEmailAddressValidationService) RequestValidation(request *http.Request, username string, email string, confirmationurl string, langKey string) (key string, err error) {
	valMngr := validation.NewManager(request)
	info, err := valMngr.NewEmailAddressValidationInformation(username, email)
	if err != nil {
		log.Error(err)
		return
	}
	err = valMngr.SaveEmailAddressValidationInformation(info)
	if err != nil {
		log.Error(err)
		return
	}

	translationValues := tools.TranslationValues{
		"emailvalidation_title":      nil,
		"emailvalidation_text":       struct{ Email string }{Email: email},
		"emailvalidation_buttontext": nil,
		"emailvalidation_reason":     nil,
		"emailvalidation_subject":    nil,
		"emailvalidation_urlcaption": nil,
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	validationurl := fmt.Sprintf("%s?c=%s&k=%s&l=%s", confirmationurl, url.QueryEscape(info.Secret), url.QueryEscape(info.Key), langKey)
	templateParameters := EmailWithButtonTemplateParams{
		UrlCaption: translations["emailvalidation_urlcaption"],
		Url:        validationurl,
		Username:   username,
		Title:      translations["emailvalidation_title"],
		Text:       translations["emailvalidation_text"],
		ButtonText: translations["emailvalidation_buttontext"],
		Reason:     translations["emailvalidation_reason"],
		LogoUrl:    fmt.Sprintf("https://%s/assets/img/its-you-online.png", request.Host),
	}
	message, err := tools.RenderTemplate(emailWithButtonTemplateName, templateParameters)
	if err != nil {
		return
	}

	go service.EmailService.Send([]string{email}, translations["emailvalidation_subject"], message)
	key = info.Key
	return
}

//RequestPasswordReset Request a password reset
func (service *IYOEmailAddressValidationService) RequestPasswordReset(request *http.Request, username string, emails []string, langKey string) (key string, err error) {
	pwdMngr := password.NewManager(request)
	token, err := pwdMngr.NewResetToken(username)
	if err != nil {
		return
	}
	if err = pwdMngr.SaveResetToken(token); err != nil {
		return
	}

	translationValues := tools.TranslationValues{
		"passwordreset_title":      nil,
		"passwordreset_text":       nil,
		"passwordreset_buttontext": nil,
		"passwordreset_reason":     nil,
		"passwordreset_subject":    nil,
		"passwordreset_urlcaption": nil,
	}

	translations, err := tools.ParseTranslations(langKey, translationValues)
	if err != nil {
		log.Error("Failed to parse translations: ", err)
		return
	}

	passwordreseturl := fmt.Sprintf("https://%s/login?lang=%s#/resetpassword/%s", request.Host, langKey, url.QueryEscape(token.Token))
	templateParameters := EmailWithButtonTemplateParams{
		UrlCaption: translations["passwordreset_urlcaption"],
		Url:        passwordreseturl,
		Username:   username,
		Title:      translations["passwordreset_title"],
		Text:       translations["passwordreset_text"],
		ButtonText: translations["passwordreset_buttontext"],
		Reason:     translations["passwordreset_reason"],
		LogoUrl:    fmt.Sprintf("https://%s/assets/img/its-you-online.png", request.Host),
	}
	message, err := tools.RenderTemplate(emailWithButtonTemplateName, templateParameters)
	if err != nil {
		return
	}
	go service.EmailService.Send(emails, translations["passwordreset_subject"], message)
	key = token.Token
	return
}

//SendOrganizationInviteEmail Sends an organization invite email
func (service *IYOEmailAddressValidationService) SendOrganizationInviteEmail(request *http.Request, invite *invitations.JoinOrganizationInvitation) (err error) {
	InviteURL := fmt.Sprintf(invitations.InviteURL, request.Host, url.QueryEscape(invite.Code))
	templateParameters := EmailWithButtonTemplateParams{
		Url:        InviteURL,
		Username:   invite.EmailAddress,
		Title:      "It's You Online organization invitation",
		Text:       fmt.Sprintf("You have been invited to the %s organization on It's You Online. Click the button below to accept the invitation.", invite.Organization),
		ButtonText: "Accept invitation",
		Reason:     "You’re receiving this email because someone invited you to an organization at ItsYou.Online. If you think this was a mistake please ignore this email.",
		LogoUrl:    fmt.Sprintf("https://%s/assets/img/its-you-online.png", request.Host),
	}
	message, err := tools.RenderTemplate(emailWithButtonTemplateName, templateParameters)
	if err != nil {
		return
	}
	subject := fmt.Sprintf("You have been invited to the %s organization", invite.Organization)
	recipients := []string{invite.EmailAddress}
	go service.EmailService.Send(recipients, subject, message)
	return
}

//ExpireValidation removes a pending validation
func (service *IYOEmailAddressValidationService) ExpireValidation(request *http.Request, key string) (err error) {
	if key == "" {
		return
	}
	valMngr := validation.NewManager(request)
	err = valMngr.RemoveEmailAddressValidationInformation(key)
	return
}

func (service *IYOEmailAddressValidationService) getEmailAddressValidationInformation(request *http.Request, key string) (info *validation.EmailAddressValidationInformation, err error) {
	if key == "" {
		return
	}
	valMngr := validation.NewManager(request)
	info, err = valMngr.GetByKeyEmailAddressValidationInformation(key)
	return
}

//IsConfirmed checks wether a validation request is already confirmed
func (service *IYOEmailAddressValidationService) IsConfirmed(request *http.Request, key string) (confirmed bool, err error) {
	info, err := service.getEmailAddressValidationInformation(request, key)
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
func (service *IYOEmailAddressValidationService) ConfirmValidation(request *http.Request, key, secret string) (err error) {
	info, err := service.getEmailAddressValidationInformation(request, key)
	if err != nil {
		return
	}
	if info == nil {
		err = ErrInvalidOrExpiredKey
		return
	}
	if info.Secret != secret {
		err = ErrInvalidCode
		return
	}
	valMngr := validation.NewManager(request)
	p := valMngr.NewValidatedEmailAddress(info.Username, info.EmailAddress)
	err = valMngr.SaveValidatedEmailAddress(p)
	if err != nil {
		return
	}
	err = valMngr.UpdateEmailAddressValidationInformation(key, true)
	if err != nil {
		return
	}
	return
}
