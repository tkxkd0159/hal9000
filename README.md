![cross compile](https://github.com/Carina-labs/HAL9000/actions/workflows/build.yml/badge.svg)
![golangci-lint](https://github.com/Carina-labs/HAL9000/actions/workflows/lint.yml/badge.svg)
![LoC](https://img.shields.io/badge/line%20of%20codes-3094-informational)

<!-- TOC -->
* [HAL9000](#hal9000)
  * [Bot types](#bot-action-types)
  * [Event](#event)
* [Run the HAL9000](#run-the-hal9000)
<!-- TOC -->

# HAL9000

```bash
GOPRIVATE=github.com/Carina-labs go get -u github.com/Carina-labs/nova@<tag>
```

## Bot action types
* **oracle** : Update host's base token price every 15 minutes.
* **stake** : Delegate the tokens sent by the user to the host chain via IBC to the a4x validator through the controller account every 10 mintues.
* **restake** : Automatically re-stake the host account's rewards through IBC. The amount to be re-deposited is inquired from the distribution module of the host chain every 6 hours.
* **Withdraw** : Undelegate and withdraw token from host account to nova. The interval depends on the rules of the host chain.

## Event
* With `EmitTypedEvent`
    * event type == proto package name + message name (e.g. nova.oracle.v1.ChainInfo)
    * event attribute key : proto field name
    * event attribute value : proto field's value. type is depend on proto field
*
```sh
# query 시 {eventType}.{eventAttribute}={eventValue}
curl "localhost:26657/tx_search?query=\"message.sender='cosmos1...'\"&prove=true"
```

```json
{
   "jsonrpc": "2.0",
   "method": "subscribe",
   "id": "0",
   "params": {
      "query": "tm.event='eventCategory' AND eventType.eventAttribute='attributeValue'"
   }
}
```


# Run the HAL9000
## 1. Set keyring if you need
```shell
make run FLAGS="-display -new -name=<keyname>"
```
## 2. Make `.secret.yml` to pipe keyring password
```yaml
pw: <your_keyring_password>
```
## 3. Build all bots & Run
```shell
make all [ARCH=<arm64|amd64>] # if you don't set ARCH, it follows GOARCH
./out/hal <action> [flags]  # use --help to show usages
                            # e.g. ./out/hal oracle -display -host=gaia -interval=60 -api=127.0.0.1:3334 -logloc=logs/oracle
```
## Appx. 1) Run bot without build (Do not use for prod)
```shell
make run ACTION=oracle FLAGS="-name=<keyname> -host=gaia -interval=900 -api=127.0.0.1:3334 -logloc=logs/oracle -display"
make run ACTION=stake FLAGS="-name=<keyname> -host=gaia -interval=600 -api=127.0.0.1:3335 -logloc=logs/stake -display"
make run ACTION=restake FLAGS="-name=<keyname> -host=gaia -interval=21600 -api=127.0.0.1:3336 -logloc=logs/restake -display"
make run ACTION=withdraw FLAGS="-name=<keyname> -host=gaia -ch=<ch_id> -interval=1814400 -api=127.0.0.1:3337 -logloc=logs/withdraw -display"
```

## Appx. 2) Run as container
```shell
# Set keyring (only need when start first)
docker run --rm -it -v halscrt:/workspace/keyring a41ventures/hal:0.1.1 <action> -display -new

# Run HAL9000
docker run --rm -d -p <host_port>:<bot_monitoring_port> -v halscrt:/workspace/keyring -v hallog:/workspace/logs a41ventures/hal:0.1.1 <action> [flags]
```
