FROM golang:1.19-alpine AS builder

FROM builder AS builder-amd64
ARG arch=x86_64

FROM builder AS builder-arm64
ARG arch=aarch64

FROM builder-$TARGETARCH AS release
RUN set -eux; apk add --no-cache ca-certificates build-base;
RUN apk add git
ARG GHTOKEN
RUN git config --global url."https://$GHTOKEN@github.com/".insteadOf "https://github.com/"
WORKDIR /workspace
COPY . .
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9ecb037336bd56076573dc18c26631a9d2099a7f2b40dc04b6cae31ffb4c8f9a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 6e4de7ba9bad4ae9679c7f9ecf7e283dd0160e71567c6a7be6ae47c81ebe7f32
RUN cp /lib/libwasmvm_muslc.${arch}.a /lib/libwasmvm_muslc.a
RUN LINK_STATICALLY=true make build

FROM alpine:3.16
ARG USER=hal
ENV HOME /home/$USER
RUN apk add --update sudo ca-certificates libstdc++ yq
RUN adduser -D $USER \
        && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
        && chmod 0440 /etc/sudoers.d/$USER
USER $USER
WORKDIR $HOME

ENV TARGET=$USER
ENV PATH="${PATH}:$HOME"
COPY --from=release /workspace/build/$TARGET ./$TARGET
# comment out below if you need config dynamic linking
COPY .chaininfo.yaml .secret.yaml ./config/
RUN mkdir /home/hal/keyring && chown hal:hal /home/hal/keyring
CMD ["hal", "--help"]
