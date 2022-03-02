# rtoken-relay-core

rToken relay core service to interact with Multiple Networks; StaFi Hub, Cosmos etc.

## build

```
git clone https://github.com/stafihub/rtoken-relay-core.git

cd rtoken-relay-core

make build

```

## run 

```
./build/relay --config ./config_template_stafihub_cosmoshub.json
```

## keytool 

```
 ./build/keytool add keyname --prefix stafi --home ./keys/stafi

 ./build/keytool list --prefix stafi --home ./keys/stafi
```