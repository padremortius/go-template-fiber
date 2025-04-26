#!/bin/sh
mkdir spec
go-swagger3 --module-path . --output spec/docs.json --schema-without-pkg
