package validation

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/itsyouonline/identityserver/credentials/password"
	"github.com/itsyouonline/identityserver/db/validation"
	"github.com/itsyouonline/identityserver/tools"
	"net/http"
	"net/url"
)

//EmailService is the interface for an email communication channel, should be used by the IYOEmailAddressValidationService
type EmailService interface {
	Send(receipients []string, subject string, message string) (err error)
}

//IYOEmailAddressValidationService is the itsyou.online implementation of a EmailAddressValidationService
type IYOEmailAddressValidationService struct {
	EmailService EmailService
}

//RequestValidation validates the email address by sending an EMail
func (service *IYOEmailAddressValidationService) RequestValidation(request *http.Request, username string, email string, confirmationurl string) (key string, err error) {
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
	validationurl := fmt.Sprintf("%s?c=%s&k=%s", confirmationurl, url.QueryEscape(info.Secret), url.QueryEscape(info.Key))
	templateParameters := struct {
		Url        string
		Username   string
		Title      string
		Text       string
		ButtonText string
		Reason     string
	}{
		Url:        validationurl,
		Username:   username,
		Title:      "It'sYou.Online email verification",
		Text:       fmt.Sprintf("To verify your email address %s on ItsYou.Online, click the button below.", email),
		ButtonText: "Verify email",
		Reason:     " You’re receiving this email because you recently created a new ItsYou.Online account or added a new email address. If this wasn’t you, please ignore this email.",
	}
	message, err := tools.RenderTemplate("templates/email/emailwithbutton.html", templateParameters)
	if err != nil {
		return
	}

	go service.EmailService.Send([]string{email}, "ItsYou.Online email verification", message)
	key = info.Key
	return
}

//RequestPassword reset
func (service *IYOEmailAddressValidationService) RequestPasswordReset(request *http.Request, username string, emails []string) (key string, err error) {
	pwdMngr := password.NewManager(request)
	token, err := pwdMngr.NewResetToken(username)
	if err != nil {
		return
	}
	if err = pwdMngr.SaveResetToken(token); err != nil {
		return
	}

	passwordreseturl := fmt.Sprintf("https://%s/login#/resetpassword?token=%s", request.Host, url.QueryEscape(token.Token))
	templateParameters := struct {
		Url        string
		Username   string
		Title      string
		Text       string
		ButtonText string
		Reason     string
	}{
		Url:        passwordreseturl,
		Username:   username,
		Title:      "It's You Online password reset",
		Text:       "To reset your ItsYou.Online password, click the button below.",
		ButtonText: "Reset password",
		Reason:     "You’re receiving this email because you recently requested to reset your password at ItsYou.Online. If this wasn’t you, please ignore this email.",
	}
	message, err := tools.RenderTemplate("templates/email/emailwithbutton.html", templateParameters)
	if err != nil {
		return
	}
	go service.EmailService.Send(emails, "ItsYou.Online password reset", message)
	key = token.Token
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
