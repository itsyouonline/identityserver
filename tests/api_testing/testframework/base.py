from client import Client as APIClient
#from client_jwt import JWTClient
import requests


class Oauth:
    def __init__(self, env_url):
        self.session = requests.Session()
        self.env_url = env_url

    def login_via_client_credentials(self, client_id, client_secret):
        url = self.env_url + 'v1/oauth/access_token'
        params = {'grant_type': 'client_credentials',
                  'client_id': client_id,
                  'client_secret': client_secret}
        data = self.session.post(url, params=params)
        if data.status_code != 200:
            raise RuntimeError("Failed to login")
        token = data.json()['access_token']
        self.session.headers['Authorization'] = 'token {token}'.format(token=token)
        self.session.headers['Content-Type'] = "application/json"

class Client:
    def __init__(self, env_url):
        session = requests.Session()
        self.api = APIClient(env_url)
        self.api.session = session
        #self.jwt = JWTClient(env_url)
        #self.jwt.session = session
        self.oauth = Oauth(env_url)
        self.oauth.session = session
