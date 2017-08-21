#!/usr/bin/python
from optparse import OptionParser

try:
    import configparser
except:
    from six.moves import configparser

parser = OptionParser()
parser.add_option('--env', help='itsyouonline url', dest='itsyouonline_url', default='', action='store')
parser.add_option('--vmail', help='validation mail', dest='validation_email', default='', action='store')
parser.add_option('--vmail-pasword', help='validation mail password', dest='validation_email_password', default='', action='store')
parser.add_option('--user1-authsecret', help='user_1 auth secret', dest='user1_totp_secret', default='', action='store')
parser.add_option('--user1', help='username_1', dest='user1_username', default='', action='store')
parser.add_option('--user1-password', help='user_1 password', dest='user1_password', default='', action='store')
parser.add_option('--user1-appid', help='user_1 applicationid', dest='user1_applicationid', default='', action='store')
parser.add_option('--user1-appsecret', help='user_1 clien secret', dest='user1_secret', default='', action='store')
parser.add_option('--user2', help='username_2', dest='user2_username', default='', action='store')
parser.add_option('--user2-password', help='user_2 password', dest='user2_password', default='', action='store')
parser.add_option('--user2-appid', help='user_2 applicationid', dest='user2_applicationid', default='', action='store')
parser.add_option('--user2-appsecret', help='user_2 clien secret', dest='user2_secret', default='', action='store')

(options, args) = parser.parse_args()

config = configparser.ConfigParser()
config.read('config.ini')
config.set('main', 'itsyouonline_url', options.itsyouonline_url)
config.set('main', 'validation_email', options.validation_email)
config.set('main', 'validation_email_password', options.validation_email_password)
config.set('main', 'user1_totp_secret', options.user1_totp_secret)
config.set('main', 'user1_username', options.user1_username)
config.set('main', 'user1_password', options.user1_password)
config.set('main', 'user1_applicationid', options.user1_applicationid)
config.set('main', 'user1_secret', options.user1_secret)
config.set('main', 'user2_username', options.user2_username)
config.set('main', 'user2_password', options.user2_password)
config.set('main', 'user2_applicationid', options.user2_applicationid)
config.set('main', 'user2_secret', options.user2_secret)

with open('config.ini', 'w') as configfile:
    config.write(configfile)
