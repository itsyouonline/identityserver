#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
if ! which go-bindata > /dev/null; then
    echo 'Installing go-bindata'
    go get -u github.com/jteeuwen/go-bindata/...
fi
pushd siteservice/website > /dev/null
echo 'Switching assets in debug'
go-bindata -debug -pkg assets -prefix assets -o ./packaged/assets/assets.go assets/...
echo 'Switching 3rd party assets in debug'
go-bindata -debug -pkg thirdpartyassets -prefix thirdpartyassets -o ./packaged/thirdpartyassets/thirdpartyassets.go thirdpartyassets/...

echo 'Switching components in debug'
go-bindata -debug -pkg components -prefix components -o ./packaged/components/components.go components/...

echo 'Switching html in debug'
go-bindata -debug -pkg html -o ./packaged/html/html.go index.html registration.html login.html privacypolicy.html error.html apidocumentation.html emailconfirmation.html smsconfirmation.html base.html
popd > /dev/null
echo "Switching templates to debug mode"
pushd templates > /dev/null
go-bindata -debug -pkg templates -prefix templates -o packaged/templates.go templates/...
popd > /dev/null
echo 'Build project with switched assets'
go build
