from docs.examples.thirdPartyService.dronedeliveryConsumer.droneDeliveryClient.client import Client
from clients.python import itsyouonline
import requests

itsyouonline_client = itsyouonline.Client()
itsyouonline_client_applicationid = 'FILL_IN_APPLICATIONID'
itsyouonline_client_secret = 'FILL_IN_CLIENTSECRET'
itsyouonline_client.oauth.LoginViaClientCredentials(client_id=itsyouonline_client_applicationid,
                                                    client_secret=itsyouonline_client_secret)


def getjwt():
    uri = "https://itsyou.online/v1/oauth/jwt?aud=dronedelivery&scope="
    r = requests.post(uri, headers=itsyouonline_client.oauth.session.headers)
    if r.status_code != 200:
        return r.text
    else:
        return r.status_code


dronedelivery_consumer_client = Client()
dronedelivery_consumer_client.url = "http://127.0.0.1:5000"
jwt = getjwt()
dronedelivery_consumer_client.set_auth_header("bearer %s" % jwt)
r = dronedelivery_consumer_client.deliveries_get()

print(r)
print(r.json())
