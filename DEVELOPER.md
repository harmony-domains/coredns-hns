# coredns-1ns developer guide

## Quick start

```
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

## Updating configuration
There are two configurtion files needed included in the repository and copied at time of bulding the standalone test environment
* `.env` contains environment variables for go-1ns including client connection info and deployed contract addresses
* `Corefile` contains zone and server information see [here](https://coredns.io/2017/07/23/corefile-explained/) for more information.
## Running go tests


## End to End Local Testing

For the following it is assumed that ens-deployer and coredns-1ns have the same parent directory.

Start a local ganache instance

```
cd ens-deployer/env
./ganache-new.sh
```

Deploy ens contracts, register test.country and create config records for them

```
cd ens-deployer/contract
yarn deploy --network local

```


Build local coredns and copy Corefile.local to local

```
cd coredns-1ns
./build-standalone.sh
```

Start coredns

```
./coredns
```

Test we are retrieving the dns records succesfully

```
# test retreival of stored DNS test record
dig @127.0.0.1 -p 53 test.country

# test forwarding is working
dig @127.0.0.1 -p 53 www.example.com

```




## Publishing new versions

See [here](https://go.dev/doc/modules/publishing) and [here](https://go.dev/blog/publishing-go-modules).

```
git tag v0.1.0
git push origin v0.1.0
```


## Appendices

### Appendix A sample local Corefile

```
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

        # ipfsgatewaya is the address of an ONENS-enabled IPFS gateway.
        # This value is returned when a request for an A record of an Ethlink
        # domain is received and the domain has a contenthash record in ONENS but
        # no A record.  Multiple values can be supplied, separated by a space,
        # in which case all records will be returned.
        ipfsgatewaya 176.9.154.81

        # ipfsgatewayaaaa is the address of an ONENS-enabled IPFS gateway.
        # This value is returned when a request for an AAAA record of an Ethlink
        # domain is received and the domain has a contenthash record in ONENS but
        # no A record.  Multiple values can be supplied, separated by a space,
        # in which case all records will be returned.
        ipfsgatewayaaaa 2a01:4f8:160:4069::2
    }
    forward . 8.8.8.8 9.9.9.9 # Forward all requests to Google Public DNS (8.8.8.8) and Quad9 DNS (9.9.9.9).
}
```