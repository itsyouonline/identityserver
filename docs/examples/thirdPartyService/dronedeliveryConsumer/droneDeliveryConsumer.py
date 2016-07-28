from docs.examples.thirdPartyService.dronedeliveryConsumer.droneDeliveryClient.client import Client
from clients.python import itsyouonline
import requests

itsyouonline_client = itsyouonline.Client()
itsyouonline_client_applicationid = 'Zfv5-Ew9Blv900ZyWz5BvBvuTes9'
itsyouonline_client_secret = 'CiLj3xMpBaiYzBXJ0iAe_HLPWaTi'
itsyouonline_client.oauth.LoginViaClientCredentials(client_id=itsyouonline_client_applicationid,
                                                    client_secret=itsyouonline_client_secret)


def getjwt():
    uri = "https://itsyou.online/v1/oauth/jwt?aud=dronedelivery&scope="
    r = requests.post(uri, headers=itsyouonline_client.oauth.session.headers)
    return r.text


dronedelivery_consumer_client = Client()
dronedelivery_consumer_client.url = "http://127.0.0.1:5000"
jwt = getjwt()
dronedelivery_consumer_client.set_auth_header("bearer %s" % jwt)
r = dronedelivery_consumer_client.deliveries_get()

print(r)
print(r.json())
