# rtoken-relay-core

rToken relay core service to interact with Multiple Networks; StaFi Hub, Cosmos etc.

## Build

```
git clone https://github.com/stafihub/rtoken-relay-core.git

cd rtoken-relay-core

make install

```

## Run 

**commands list:**
```
relay

Usage:
  relay [command]

Available Commands:
  start             Start relay procedure
  version           Show version information
  keys              Key tool to manage keys
  multisig-transfer Tranfer token from multisig account
  help              Help about any command

Flags:
  -h, --help   help for relay

Use "relay [command] --help" for more information about a command.
```

**show version:**

```shell
relay version
```

**start relay:**

```
relay start --config ./config_template_stafihub_cosmoshub.json
```

**manage keys:**

```shell
relay keys [command]
```