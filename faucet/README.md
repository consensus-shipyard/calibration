# Mycelium Calibration Faucet

## How to Run

- Install Golang
- 
- Build the service:
```bash
cd faucet
make build
```

- Run the service:
```bash
./faucet --web-host "127.0.0.1:80" --web-allowed-origins "http://service.name" --web-backend-host "http://service.name/fund" \
    --ethereum-private-key "key" --ethereum-url "url"
```

## Health API

- To check service readiness: `GET /readiness`
- To check service liveness: `GET /liveness`


## Configuration

### Logging

Logging can be configured through two main environment variables: `GOLOG_LOG_LEVEL` and `GOLOG_FILE`.
More information about how to configure logging you can find [here](https://github.com/ipfs/go-log#environment-variables).

For example:
```bash
GOLOG_FILE="./faucet.logs" GOLOG_LOG_LEVEL="debug" ./faucet --web-host "127.0.0.1:80" --web-allowed-origins "http://service.name" --web-backend-host "http://service.name/fund" \
--ethereum-private-key "key" --ethereum-url "url"
```

### Private Key
 - The private key can be provided directly via CLI or stored in a file. 
 - The private key must not contain "0x"
 - The private key file must not contain new line characters.

### Enabled TLS
```bash
./faucet --tls-enabled --web-allowed-origins "https://frontend" --web-backend-host "https://faucet/fund" \
    --tls-cert-file "path_to_cert.pem" --tls-key-file "path_to_key.pem" \
    --ethereum-private-key "key" --ethereum-url "url"
```
### Disabled TLS

```bash
./faucet --web-allowed-origins "http://frontend" --web-backend-host "https://faucet/fund" \
    --ethereum-private-key "key" --ethereum-url "url"
```

## Development

### TLS
To run the service in the development mode with TLS, you must provide an X509 certificate.

The easiest way to do that is to use [mkcert](https://github.com/FiloSottile/mkcert)
tool and `make cert` command.

Run `make all` to ensure that tests pass and `make demo` to run a demo accessible on localhost.

### Ganache Test Accounts

A `Ganache` node is started locally by `make node-start` and can be stopped by `make node-stop`.

These accounts and keys are used for local demos and testing.

```bash
Available Accounts
==================
(0) 0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1 (1000 ETH)
(1) 0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0 (1000 ETH)
(2) 0x22d491Bde2303f2f43325b2108D26f1eAbA1e32b (1000 ETH)
(3) 0xE11BA2b4D45Eaed5996Cd0823791E0C93114882d (1000 ETH)
(4) 0xd03ea8624C8C5987235048901fB614fDcA89b117 (1000 ETH)
(5) 0x95cED938F7991cd0dFcb48F0a06a40FA1aF46EBC (1000 ETH)
(6) 0x3E5e9111Ae8eB78Fe1CC3bb8915d5D461F3Ef9A9 (1000 ETH)
(7) 0x28a8746e75304c0780E011BEd21C72cD78cd535E (1000 ETH)
(8) 0xACa94ef8bD5ffEE41947b4585a84BdA5a3d3DA6E (1000 ETH)
(9) 0x1dF62f291b2E969fB0849d99D9Ce41e2F137006e (1000 ETH)

Private Keys
==================
(0) 0x4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d
(1) 0x6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1
(2) 0x6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c
(3) 0x646f1ce2fdad0e6deeeb5c7e8e5543bdde65e86029e2fd9fc169899c440a7913
(4) 0xadd53f9a7e588d003326d1cbf9e4a43c061aadd9bc938c843a79e7b4fd2ad743
(5) 0x395df67f0c2d2d9fe1ad08d1bc8b6627011959b79c53d7dd6a3536a33ab8a4fd
(6) 0xe485d098507f54e7733a205420dfddbe58db035fa577fc294ebd14db90767a52
(7) 0xa453611d9419d0e56f499079478fd72c37b251a94bfde4d19872c44cf65386e3
(8) 0x829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4
(9) 0xb0057716d5917badaf911b193b12b910811c1497b5bada8d7711f758981c3773
```