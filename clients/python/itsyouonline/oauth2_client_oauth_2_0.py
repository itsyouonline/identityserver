import requests


class Oauth2ClientOauth_2_0():
    def __init__(self, access_token_uri='https://itsyou.online/v1/oauth/access_token'):
        self.access_token_uri = access_token_uri

    def get_access_token(self, client_id, client_secret, scopes=[], audiences=[]):
        params = {
            'grant_type': 'client_credentials',
            'client_id': client_id,
            'client_secret': client_secret
        }
        if len(scopes) > 0:
            params['scope'] = ",".join(scopes)
        if len(audiences) > 0:
            params['aud'] = ",".join(audiences)

        return requests.post(self.access_token_uri, params=params)
