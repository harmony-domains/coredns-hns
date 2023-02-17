#!/bin/bash
set -e

SRCDIR=~/coredns-1n
BUILDDIR=~/coredns-1n/build

mkdir -p ${BUILDDIR} 2>/dev/null
cd ${BUILDDIR}
echo "Cloning coredns repo..."
git clone https://github.com/coredns/coredns.git

cd coredns
git checkout v1.10.1

echo "Patching plugin config..."
cp ../coredns-plugin.cfg plugin.cfg

# Add our module to coredns.
echo "YOU NEED to add this line to go.mod"
echo "	github.com/jw-1ns/coredns-1ns v0.1.2"
cp ../coredns-go.mod go.mod

go get github.com/jw-1ns/coredns-1ns@v0.1.2
go get
go mod download

echo "Building..."
# make SHELL='sh -x' CGO_ENABLED=1 coredns
go generate
go build

cd ../..
cp ./build/coredns/coredns .
cp ./Corefile.mainnet Corefile
rm -r ${BUILDDIR}

echo "==== NEXT STEPS ====="
echo "ensure you have a .env file"
echo "run coredns using sudo ./coredns"
