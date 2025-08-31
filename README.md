```bash
sudo apt update && sudo apt install cmake pkg-config build-essential libopenblas-dev zlib1g-dev portaudio19-dev sox
```

```bash
mkdir -p third_party
cd third_party

git clone https://github.com/kaldi-asr/kaldi.git
cd third_party/kaldi/tools
extras/check_dependencies.sh
make -j"$(nproc)"
cd ../src
./configure --shared --use-cuda=no --mathlib=OPENBLAS
make -j"$(nproc)"
cd ../..

git clone https://github.com/alphacep/vosk-api 
cd third_party/vosk-api/src
make -j"$(nproc)"
mkdir -p third_party/vosk/include
mkdir -p third_party/vosk/lib
cp vosk_api.h third_party/vosk/include/
cp libvosk.so third_party/vosk/lib/
```
