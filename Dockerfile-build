FROM debian:bookworm AS build

ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y --no-install-recommends build-essential git cmake automake autoconf libtool \
    pkg-config zlib1g-dev ca-certificates python3 python3-dev unzip wget curl sox subversion gfortran \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt

RUN git clone -b vosk --single-branch --depth=1 https://github.com/alphacep/kaldi

RUN git clone --depth=1 https://github.com/alphacep/vosk-api

# Kaldi

WORKDIR /opt/kaldi/tools

RUN make -j"$(nproc)" openfst cub

RUN ./extras/install_openblas_clapack.sh

WORKDIR /opt/kaldi/src

RUN ./configure --mathlib=OPENBLAS_CLAPACK --shared

RUN make -j"$(nproc)" online2 lm rnnlm

# Vosk api

WORKDIR /opt/vosk-api/src

RUN KALDI_ROOT=/opt/kaldi make -j"$(nproc)"

RUN install -Dm755 libvosk.so /opt/out-vosk/lib/libvosk.so \
    && install -Dm644 vosk_api.h /opt/out-vosk/include/vosk_api.h

# Golosok

FROM debian:bookworm AS app

RUN apt update && apt install -y --no-install-recommends git cmake wget unzip curl ca-certificates scons build-essential \
    pkg-config libopenblas0 libgfortran5 libstdc++6 libasound2 libasound2-dev libportaudio2 portaudio19-dev \
    libjack-jackd2-0 libjack-jackd2-dev \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --recursive https://github.com/RHVoice/RHVoice.git

WORKDIR /opt/RHVoice

RUN scons

RUN scons install

WORKDIR /opt

ENV GOLANG_VERSION=1.24.6

RUN wget -q https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz

RUN rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

ENV CGO_ENABLED=1

ENV PKG_CONFIG_PATH="/usr/lib/x86_64-linux-gnu/pkgconfig:/usr/share/pkgconfig:/usr/local/lib/pkgconfig"

ENV LD_LIBRARY_PATH="/opt/golosok/build/vosk"

WORKDIR /opt/golosok

COPY . .

RUN wget -q https://alphacephei.com/vosk/models/vosk-model-small-ru-0.22.zip

RUN unzip vosk-model-small-ru-0.22.zip -d models

RUN rm vosk-model-small-ru-0.22.zip

RUN mkdir -p build/vosk

COPY --from=build /opt/out-vosk /opt/golosok/build/vosk

RUN pkg-config --print-errors --exists alsa portaudio-2.0

RUN pkg-config --cflags --libs alsa portaudio-2.0

RUN make build

ENTRYPOINT ["/opt/golosok/build/golosok"]

CMD ["-stt-test", "0"]