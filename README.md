![cross compile](https://github.com/Carina-labs/HAL9000/actions/workflows/build.yml/badge.svg)
![golangci-lint](https://github.com/Carina-labs/HAL9000/actions/workflows/lint.yml/badge.svg)
![LoC](https://img.shields.io/badge/line%20of%20codes-3917-informational)

# HAL9000
The world's most complete oracle feeder

## Bot types
* **oracle** : Update host's base token price every 15 minutes.
* **stake** : Delegate the tokens sent by the user to the host chain via IBC to the a4x validator through the controller account every 10 mintues.
* **restake** : Automatically re-stake the host account's rewards through IBC. The amount to be re-deposited is inquired from the distribution module of the host chain every 6 hours.
* **Withdraw** : Undelegate and withdraw token from host account to nova. The interval depends on the rules of the host chain.

```bash
GOPRIVATE=github.com/Carina-labs go get -u github.com/Carina-labs/nova@v0.5.1
```


# Cmd
```bash
# Set keyring if you need
make run FLAGS="-display -new -name=<keyname>"

# Build all bots
make all [ARCH=<arm64|amd64>] # if you don't set ARCH, it follows GOARCH

# Run bot without build (test)
make run TARGET=oracle FLAGS="-name=<keyname> -host=gaia -interval=5 -api=127.0.0.1:3334 -display"
make run TARGET=stake FLAGS="-name=<keyname> -host=gaia -interval=5 -api=<addr> -display"
make run TARGET=restake FLAGS="-name=<keyname> -host=gaia -interval=5 -api=<addr> -display"
make run TARGET=withdraw FLAGS="-name=<keyname> -host=gaia -ch=channel-45 -interval=5 -api=<addr> -display"

# Run bot (prod)
./out/<bot> [flags]  # use --help to show usages

```
