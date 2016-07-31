# Creating a rest service secured by itsyou.online

## Introduction

This tutorial is going to create a simple restservice, 'DroneDelivery', that is defined using raml and secured using itsyou.online.


## Step 1, create a Flask/Python Server using go-raml


- First, make sure you have [Jumpscale go-raml](https://github.com/Jumpscale/go-raml) installed.

- Clone this repository
```
git clone https://github.com/itsyouonline/identityserver
```

- Generate the server code
```
cd identityserver/docs/examples/thirdPartyService
go-raml server -l python --dir ./dronedeliveryService --ramlfile api.raml
```

This will result in a new directory with this structure:

```
dronedeliveryService
├── apidocs
│   └── ...
├── app.py
├── deliveries.py
├── drones.py
├── index.html
├── input_validators.py
└── requirements.txt
```
Go into dronedeliveryService/deliveries.py and add to the method deliveries_get the following:

```
data = {
        "id": "4",
        "at": "Tue, 08 Jul 2014 13:00:00 GMT",
        "toAddressId": "gi6w4fgi",
        "orderItemId": "6782798",
        "status": "completed",
        "droneId": "f"
    }
    return Response(json.dumps(data), mimetype='application/json')
```
This can be found in the RAML file

To launch the server in this directory, go to the terminal and enter:

`python3 app.py`

To view the RAML specs, open your browser and go to http://127.0.0.1:5000/apidocs/index.html?raml=api.raml


## Step 2, create client

```
cd identityserver/docs/examples/thirdPartyService
go-raml client --language python --dir ./dronedeliveryConsumer/droneDeliveryClient --ramlfile api.raml
```

A python 3.5 compatible client is generated in thirdPartyService/droneDeliveryConsumer directory.


## Step 3, secure your application with itsyou.online

Make an account on https://www.itsyou.online.

go to settings, add API keys.

Make a new file in dronedeliveryConsumer folder called dronedeliveryConsumer.py

In python 3.4 you can use this code:
```
from docs.examples.thirdPartyService.dronedeliveryConsumer.droneDeliveryClient.client import Client
from clients.python import itsyouonline
import requests

itsyouonline_client = itsyouonline.Client()
itsyouonline_client_applicationid = 'FILL_IN_APPLICATIONID'
itsyouonline_client_secret = 'FILL_IN_CLIENTSECRET'
itsyouonline_client.oauth.LoginViaClientCredentials(client_id=itsyouonline_client_applicationid,client_secret=itsyouonline_client_secret)


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
```

You can find the application_id and secret in the API key you made on itsyou.online

itsyouonline documentation : https://www.gitbook.com/book/gig/itsyouonline/details


To decode the jwt you install pyJwt : https://github.com/jpadilla/pyjwt

Then we add this code to dronedeliveryService/deliveries.py deliveries_get:

`import jwt`

```
pubkey = "-----BEGIN PUBLIC KEY-----\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2\n7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6\n6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv\n-----END PUBLIC KEY-----"

    header = request.headers.get("Authorization")
    if not header or not header.startswith("bearer "):
        return 'Unauthorized', 401
    webtoken = header.split()
    decoded_jwt = jwt.decode(webtoken[1], pubkey, algorithms=["ES384"], audience="dronedelivery")
    print(decoded_jwt)
    if decoded_jwt["iss"] != "itsyouonline":
        return 'Unauthorized', 401
    else:
        data = {
            "id": "4",
            "at": "Tue, 08 Jul 2014 13:00:00 GMT",
            "toAddressId": "gi6w4fgi",
            "orderItemId": "6782798",
            "status": "completed",
            "droneId": "f"
        }
    return Response(json.dumps(data), mimetype='application/json')
```

Start up the server, run dronedeliveryConsumer.py and your application has been secured by itsyou.online!
