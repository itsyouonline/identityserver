package main

//This file contains the go:generate commands

// ## Embed the api specification
//go:generate go-bindata -pkg specifications -prefix specifications/api -o specifications/packaged.go specifications/api/...

//go-raml https://github.com/Jumpscale/go-raml server code generation from the RAML specification
//TODO: fix serverside code generation with go-raml

// ## API clients
//go:generate go-raml client -l go --dir clients/go/itsyouonline --ramlfile specifications/api/itsyouonline.raml --package itsyouonline
//go:generate go-raml client -l python --dir clients/python/itsyouonline --ramlfile specifications/api/itsyouonline.raml --package itsyouonline --python-unmarshall-response

// ## Website ##
//package the assets
//go:generate go-bindata -pkg assets -prefix siteservice/website/assets -o siteservice/website/packaged/assets/assets.go siteservice/website/assets/...

//package 3rd party assets
//go:generate go-bindata -pkg thirdpartyassets -prefix siteservice/website/thirdpartyassets -o siteservice/website/packaged/thirdpartyassets/thirdpartyassets.go siteservice/website/thirdpartyassets/...

//go:generate go-bindata -pkg components -prefix siteservice/website/components -ignore=_test.js$  -o siteservice/website/packaged/components/components.go siteservice/website/components/...

//package the html files
//go:generate go-bindata -pkg html -prefix siteservice/website -o siteservice/website/packaged/html/html.go siteservice/website/index.html siteservice/website/registration.html siteservice/website/login.html siteservice/website/base.html siteservice/website/error.html siteservice/website/apidocumentation.html siteservice/website/smsconfirmation.html siteservice/website/emailconfirmation.html

// ## Email templates ##
//go:generate go-bindata -pkg templates -prefix templates/templates -o templates/packaged/templates.go templates/templates/...
