# coredns-1ns developer guide

## Quick start

```bash
# Clone the repository 
git clone https://github.com/jw-1ns/coredns-1ns.git
cd coredns-1ns

# update the dependencies
go get -u ./...
go mod tidy

# Build
go build

# Build coredns standalone
./build-standalone.sh

# Run coredns
./coredns

```

Adding a package

```bash
# Get the package
go get github.com/coredns/coredns/plugin/pkg/upstream

# Modify the code to use the package

```

## Updating configuration

There are two configurtion files needed included in the repository and copied at time of bulding the standalone test environment

* `.env` contains environment variables for go-1ns including client connection info and deployed contract addresses
* `Corefile` contains zone and server information see [here](https://coredns.io/2017/07/23/corefile-explained/) for more information.

## Running go tests

## End to End Local Testing

For the following it is assumed that ens-deployer and coredns-1ns have the same parent directory.

Start a local ganache instance

```bash
cd ens-deployer/env
./ganache-new.sh
```

Deploy ens contracts, register test.country and create config records for them

```bash
cd ens-deployer/contract
yarn deploy-sample

```

Build local coredns and copy Corefile.local to local

```bash
cd coredns-1ns
./build-standalone.sh
```

Start coredns

```bash
./coredns
```

Test we are retrieving the dns records succesfully

```bash
# test retreival of stored DNS test record
dig @127.0.0.1 -p 53 test.country

# test forwarding is working
dig @127.0.0.1 -p 53 www.example.com

```

## Publishing new versions

See [here](https://go.dev/doc/modules/publishing) and [here](https://go.dev/blog/publishing-go-modules).

```bash
git tag v0.1.2
git push origin v0.1.2
```

## Appendices

### Appendix A sample local Corefile

```bash
. {
    whoami  # coredns.io/plugins/whoami
    log     # coredns.io/plugins/log
    errors  # coredns.io/plugins/errors
    rewrite stop {
    # This rewrites any requests for *.eth.link domains to *.eth internally
    # prior to being processed by the main ONENS resolver.
    name regex (.*)\.country\.link {1}.country
    answer name (.*)\.country {1}.country.link
    }
    onens {
        # connection is the connection to an Ethereum node.  It is *highly*
        # recommended that a local node is used, as remote connections can
        # cause DNS requests to time out.
        # This can be either a path to an IPC socket or a URL to a JSON-RPC
        # endpoint.
        # connection /home/ethereum/.ethereum/geth.ipc
        connection http://127.0.0.1:8545

        # ethlinknameservers are the names of the nameservers that serve
        # EthLink domains.  This will usually be the name of this server,
        # plus potentially one or more others.
        ethlinknameservers ns1.ethdns.xyz ns2.ethdns.xyz

    }
    forward . 8.8.8.8 9.9.9.9 # Forward all requests to Google Public DNS (8.8.8.8) and Quad9 DNS (9.9.9.9).
}
```

### Appendix B sample dig commands for local testing

```bash
# Get A records (initARecFQDN, initARecAFQDN)
cd 
dig @127.0.0.1 a in -p 53 recurse=false a.test.country 

# Get CNAME record (initCnameRecOneFQDN, initCnameRecWWW)
dig @127.0.0.1  cname in -p 53 recurse=false one.test.country
dig @127.0.0.1 cname in -p 53 recurse=false www.test.country
dig @127.0.0.1  in -p 53 recurse=false one.test.country

# Get SOA record (InitialSoaRecord)
dig @127.0.0.1 soa in -p 53 recurse=false test.country

# Get TXT record
# const [initTxtRec] = dns.encodeTXTRecord({ name: FQDN, text: 'magic' })
dig @127.0.0.1 txt in -p 53 recurse=false test.country

# TestSubdomainOneInitialNSnameRecord
const [initNsRec] = dns.encodeNSRecord({ name: FQDN, nsname: 'ns3.hiddenstate.xyz' })
dig @127.0.0.1 ns in -p 53 recurse=false test.country

# InitialSrvRecord
# const [initSrvRec] = dns.encodeSRVRecord({ name: 'srv.' + FQDN, rvalue: DefaultSrv })
dig @127.0.0.1 srv in -p 53 recurse=false srv.test.country


# Not Working 

# const [initCnameRecWWW] = dns.encodeCNAMERecord({ name: 'www', cname: 'test.country' })
dig @127.0.0.1 cname in -p 53 recurse=false www

# initDnameRec
# const [initDnameRec] = dns.encodeDNAMERecord({ name: 'docs.' + FQDN, dname: 'docs.harmony.one' })
dig @127.0.0.1 dname in -p 53 recurse=false docs.test.country
```
