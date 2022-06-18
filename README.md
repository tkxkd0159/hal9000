![cross compile](https://github.com/Carina-labs/HAL9000/actions/workflows/build.yml/badge.svg)
![golangci-lint](https://github.com/Carina-labs/HAL9000/actions/workflows/lint.yml/badge.svg)
![LoC](https://img.shields.io/badge/line%20of%20codes-744-informational)

# HAL9000
The world's most complete oracle feeder

# Cmd
```bash
make build all
make run TARGET=oracle CUSTOM_ORGS="-add=false -name='nova-bot' -interval=2"
GOPRIVATE=github.com/Carina-labs go get github.com/Carina-labs/nova
```
