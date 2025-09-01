```bash
docker build -f Dockerfile-dependencies -t golosok-dependencies-build .

mkdir -p build/vosk

CID=$(docker create golosok-dependencies-build)
docker cp "$CID":/opt/out-vosk ./build/vosk
docker rm "$CID"
```

```bash
sudo apt install rhvoice rhvoice-russian
```