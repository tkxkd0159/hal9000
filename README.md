![cross compile](https://github.com/Carina-labs/HAL9000/actions/workflows/build.yml/badge.svg)
![golangci-lint](https://github.com/Carina-labs/HAL9000/actions/workflows/lint.yml/badge.svg)
![LoC](https://img.shields.io/badge/line%20of%20codes-1430-informational)

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
make build all

# Set keyring
make run TARGET=oracle FLAGS="-display -add -name=nova_bot"

# Run bot
make run TARGET=oracle FLAGS="-name=nova_bot -host=gaia -interval=5 -test=false -display"
make run TARGET=stake FLAGS="-interval=5 -display"
make run TARGET=restake FLAGS="-interval=5 -display"
make run TARGET=withdraw FLAGS="-interval=5 -display"
```
