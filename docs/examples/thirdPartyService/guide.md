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


Create a new file called droneDeliveryConsumer.py in thirdPartyService/droneDeliveryConsumer and copy this code:

```
from dronedeliveryConsumer.droneDeliveryClient.client import Client

client = Client()
client.url = "http://127.0.0.1:5000"

print(client.deliveries_get().content)
```

Start up the server and then run this script in order to get data back.
