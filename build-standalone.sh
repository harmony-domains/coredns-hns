#!/bin/bash
set -e

SRCDIR=`pwd`
BUILDDIR=`pwd`/build

rm -rf ${BUILDDIR}
mkdir -p ${BUILDDIR} 2>/dev/null
cd ${BUILDDIR}
echo "Cloning coredns repo..."
git clone https://github.com/coredns/coredns.git

cd coredns
git checkout v1.10.1

echo "Patching plugin config..."
ed plugin.cfg <<EOED
/rewrite:rewrite
a
onens:github.com/jw-1ns/coredns-1ns
.
w
q
EOED

# Add our module to coredns.
echo "Patching go modules..."
ed go.mod <<EOED
a
replace github.com/jw-1ns/coredns-1ns => ../..
.
/^)
-1
a
	github.com/jw-1ns/coredns-1ns v0.1.2
.
w
q
EOED

go get github.com/jw-1ns/coredns-1ns@v0.1.2
go get
go mod download

echo "Building..."
# make SHELL='sh -x' CGO_ENABLED=1 coredns
go generate
go build
# make

cp coredns ${SRCDIR}
chmod -R 755 .git
cd ${SRCDIR}
# rm -r ${BUILDDIR}

cp Corefile.local Corefile
cp .env.local .env
