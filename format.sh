#!/usr/bin/env bash
/usr/local/go/bin/gofmt -w $(find * -type d | grep -v '^vendor\|^.git\|^packaged\|^specifications')