import requests


class JwtClient:
    def __init__(self, env_url):
        self.url = env_url + 'v1/oauth/jwt/'
        self.session = requests.Session()

    def GetScope(self, scope, headers=None):
        """
        Get Scope
        It is method for POST /scope
        """
        params = {'scope': scope}
        return self.session.post(self.env_url, params=params, headers=headers)