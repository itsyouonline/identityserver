package main

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/dgrijalva/jwt-go"
	"github.com/itsyouonline/identityserver/communication"
	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/globalconfig"
	"github.com/itsyouonline/identityserver/https"
	"github.com/itsyouonline/identityserver/identityservice"
	"github.com/itsyouonline/identityserver/identityservice/security"
	"github.com/itsyouonline/identityserver/oauthservice"
	"github.com/itsyouonline/identityserver/routes"
	"github.com/itsyouonline/identityserver/siteservice"
)

var version string

func main() {
	if version == "" {
		version = "Dev"
	}
	app := cli.NewApp()
	app.Name = "Identity server"
	app.Version = version

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	// Set log output to stdout so we can pipe it
	log.SetOutput(os.Stdout)

	var debugLogging, ignoreDevcert, testEnv bool
	var bindAddress, dbConnectionString string
	var tlsCert, tlsKey string
	var twilioAccountSID, twilioAuthToken, twilioMessagingServiceSID string
	var smtpserver, smtpuser, smtppassword string
	var smtpport int

	var smsAeroUser, smsAeroPassword, smsAeroSenderId string

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
		cli.StringFlag{
			Name:        "bind, b",
			Usage:       "Bind address",
			Value:       ":8443",
			Destination: &bindAddress,
		},
		cli.StringFlag{
			Name:        "connectionstring, c",
			Usage:       "Mongodb connection string",
			Value:       "127.0.0.1:27017",
			Destination: &dbConnectionString,
		},
		cli.StringFlag{
			Name:        "cert, s",
			Usage:       "TLS certificate path",
			Value:       "",
			Destination: &tlsCert,
		},
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "TLS private key path",
			Value:       "",
			Destination: &tlsKey,
		},
		cli.BoolFlag{
			Name:        "ignore-devcert, i",
			Usage:       "Ignore default devcert even if exists",
			Destination: &ignoreDevcert,
		},
		cli.StringFlag{
			Name:        "twilio-AccountSID",
			Usage:       "Twilio AccountSID",
			Destination: &twilioAccountSID,
		},
		cli.StringFlag{
			Name:        "twilio-AuthToken",
			Usage:       "Twilio AuthToken",
			Destination: &twilioAuthToken,
		},
		cli.StringFlag{
			Name:        "twilio-MsgSvcSID",
			Usage:       "Twilio MessagingServiceSID",
			Destination: &twilioMessagingServiceSID,
		},
		cli.StringFlag{
			Name:        "smtp-server",
			Usage:       "Host of smtp server",
			Destination: &smtpserver,
		},
		cli.StringFlag{
			Name:        "smtp-user",
			Usage:       "User to login smtp server",
			Destination: &smtpuser,
		},
		cli.StringFlag{
			Name:        "smtp-password",
			Usage:       "Password of smtp server",
			Destination: &smtppassword,
		},
		cli.IntFlag{
			Name:        "smtp-port",
			Usage:       "Port of smtp server",
			Destination: &smtpport,
			Value:       587,
		},
		cli.StringFlag{
			Name:        "SmsAeroUser",
			Usage:       "User for SmsAero",
			Destination: &smsAeroUser,
		},
		cli.StringFlag{
			Name:        "SmsAeroPassword",
			Usage:       "Password in md5 format or an apikey for SmsAero",
			Destination: &smsAeroPassword,
		},
		cli.StringFlag{
			Name:        "SmsAeroSenderId",
			Usage:       "The sender Id for SmsAero, filled in in the from field",
			Destination: &smsAeroSenderId,
		},
		cli.BoolFlag{
			Name:        "testEnv",
			Usage:       "Designate if this is a production environment",
			Destination: &testEnv,
		},
	}

	app.Before = func(c *cli.Context) error {
		if debugLogging {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
		}
		return nil
	}

	app.Action = func(c *cli.Context) {

		log.Infoln(app.Name, "version", app.Version)
		// Connect to DB!
		go db.Connect(dbConnectionString)
		defer db.Close()

		cookieSecret := identityservice.GetCookieSecret()
		var smsService communication.SMSService
		var emailService communication.EmailService
		var mailService communication.MailService
		if twilioAccountSID != "" && smsAeroPassword != "" {
			smsService = &communication.SMSServiceProxySeparateRussia{
				DefaultSMSService: &communication.TwilioSMSService{
					AccountSID:          twilioAccountSID,
					AuthToken:           twilioAuthToken,
					MessagingServiceSID: twilioMessagingServiceSID,
				},
				RussianSMSService: &communication.SmsAeroSMSService{
					Username: smsAeroUser,
					Password: smsAeroPassword,
					SenderId: smsAeroSenderId,
				},
			}
		} else if twilioAccountSID != "" {
			smsService = &communication.TwilioSMSService{
				AccountSID:          twilioAccountSID,
				AuthToken:           twilioAuthToken,
				MessagingServiceSID: twilioMessagingServiceSID,
			}
		} else if smsAeroUser != "" {
			smsService = &communication.SmsAeroSMSService{
				Username: smsAeroUser,
				Password: smsAeroPassword,
			}
		} else {
			log.Warn("============================================================================")
			log.Warn("No valid Twilio Account provided, falling back to development implementation")
			log.Warn("============================================================================")
			smsService = &communication.DevSMSService{}
		}

		if smtpserver == "" {
			log.Warn("============================================================================")
			log.Warn("No valid SMTP server provided, falling back to development implementation")
			log.Warn("============================================================================")
			emailService = &communication.DevEmailService{}

		} else {
			emailService = communication.NewSMTPEmailService(smtpserver, smtpport, smtpuser, smtppassword)
		}

		// TODO: implement actual snailmail provider
		if true {
			log.Warn("============================================================================")
			log.Warn("No valid mail service provided, falling back to development implementation")
			log.Warn("============================================================================")
			mailService = &communication.DevMailService{}
		} else {
			//TODO implement actual snailmail provider
		}

		is := identityservice.NewService(smsService, emailService, mailService)
		sc := siteservice.NewService(cookieSecret, smsService, emailService, is, version, testEnv)

		config := globalconfig.NewManager()

		var jwtKey []byte
		var err error
		exists, err := config.Exists("jwtkey")
		if err == nil && exists {
			var jwtKeyConfig *globalconfig.GlobalConfig
			jwtKeyConfig, err = config.GetByKey("jwtkey")
			jwtKey = []byte(jwtKeyConfig.Value)
		} else {
			if err == nil {
				if _, e := os.Stat("devcert/jwt_key.pem"); e == nil {
					log.Warning("===============================================================================")
					log.Warning("This instance uses a development JWT signing key, don't do this in production !")
					log.Warning("===============================================================================")

					jwtKey, err = ioutil.ReadFile("devcert/jwt_key.pem")
				}
			}
		}
		if err != nil {
			log.Fatal("Unable to load a valid key for signing JWT's: ", err)
		}
		ecdsaKey, err := jwt.ParseECPrivateKeyFromPEM(jwtKey)
		if err != nil {
			log.Fatal("Unable to load a valid key for signing JWT's: ", err)
		}
		security.JWTPublicKey = &ecdsaKey.PublicKey
		oauthsc, err := oauthservice.NewService(sc, is, ecdsaKey)
		if err != nil {
			log.Fatal("Unable to create the oauthservice: ", err)
		}

		r := routes.GetRouter(sc, is, oauthsc)

		server := https.PrepareHTTP(bindAddress, r)
		https.PrepareHTTPS(server, tlsCert, tlsKey, ignoreDevcert)

		if testEnv {
			log.Warn("Running in test environment - forget account endpoints enabled")
		}

		// Go make magic over HTTPS
		log.Info("Listening (https) on ", bindAddress)
		log.Fatal(server.ListenAndServeTLS("", ""))
	}

	app.Run(os.Args)
}
