# Creating a rest service secured by itsyou.online

## Introduction

This tutorial is going to create a simple restservice, 'DroneDelivery', that is defined using raml and secured using itsyou.online.


## Step 1, generate a Flask/Python Server using go-raml


- First, make sure you have [Jumpscale go-raml](https://github.com/Jumpscale/go-raml) installed.

- Clone this repository
```
git clone https://github.com/itsyouonline/identityserver
```

- Generate the server code
```
cd identityserver/docs/examples/thirdPartyService
go-raml server -l python --dir ./dronedelivery --ramlfile api.raml
```

This will result in a new directory with this structure:

```
result-directory
├── apidocs
│   └── ...
├── app.py
├── deliveries.py
├── drones.py
├── index.html
├── input_validators.py
└── requirements.txt
```

To launch the server in this directory, go to the terminal and enter:

`python3 app.py`

Open your browser and go to http://127.0.0.1:5000/ , then click on API docs to see the RAML specs.
