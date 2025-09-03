FROM debian:bookworm

ARG DEBIAN_FRONTEND=noninteractive

RUN printf "Types: deb\nURIs: https://deb.debian.org/debian\nSuites: bookworm-backports\nComponents: main contrib non-free non-free-firmware\nSigned-By: /usr/share/keyrings/debian-archive-keyring.gpg\n" \
    > /etc/apt/sources.list.d/backports.sources

RUN apt update && apt install -y --no-install-recommends git curl wget unzip build-essential pkg-config ca-certificates \
    libasound2 libasound2-dev libportaudio2 portaudio19-dev libjack-jackd2-0 libjack-jackd2-dev libopenblas0 \
    libgfortran5 libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

RUN apt update && apt -y -t bookworm-backports install rhvoice speech-dispatcher-rhvoice rhvoice-russian \
    && rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:${PATH}"
ENV CGO_ENABLED=1
ENV PKG_CONFIG_PATH="/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/share/pkgconfig:/usr/local/lib/pkgconfig"
ENV VOSK_VERSION=0.3.45
ENV GOLANG_VERSION=1.24.6

RUN wget -q https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz

RUN rm go${GOLANG_VERSION}.linux-amd64.tar.gz

WORKDIR /opt/golosok

COPY . .

RUN mkdir -p build/vosk

RUN wget -q https://github.com/alphacep/vosk-api/releases/download/v${VOSK_VERSION}/vosk-linux-x86_64-${VOSK_VERSION}.zip

RUN unzip -q vosk-linux-x86_64-${VOSK_VERSION}.zip -d /tmp/vosk

RUN cp /tmp/vosk/vosk-linux-x86_64-${VOSK_VERSION}/libvosk.so build/vosk/

RUN cp /tmp/vosk/vosk-linux-x86_64-${VOSK_VERSION}/vosk_api.h build/vosk/ || true

RUN rm -rf /tmp/vosk vosk-linux-x86_64-${VOSK_VERSION}.zip

RUN wget -q https://alphacephei.com/vosk/models/vosk-model-small-ru-0.22.zip

RUN unzip -q vosk-model-small-ru-0.22.zip -d models

RUN rm vosk-model-small-ru-0.22.zip

RUN pkg-config --print-errors --exists alsa portaudio-2.0

RUN pkg-config --cflags --libs alsa portaudio-2.0

ENV LD_LIBRARY_PATH="/opt/golosok/build/vosk"

RUN make build

ENTRYPOINT ["/opt/golosok/build/golosok"]

CMD ["-stt-test", "0"]
