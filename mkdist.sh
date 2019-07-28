#!/bin/bash

mkdir -p dist
rm -rf dist/*
out=dist/ud
go build -o $out ./cmd/ud
upx $out
cp CHANGELOG.md dist/
# Get the version from the binary
ver=ud-$(./dist/ud -v)
rm -f $ver.tgz
mv dist $ver
tar -c $ver | gzip - > $ver.tgz
rm -rf $ver
tar tvfz $ver.tgz
