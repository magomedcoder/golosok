```bash
docker build -t golosok .

# Запуск на Linux
docker run --rm -it --device /dev/snd --group-add audio --name golosok golosok

# Запуск на Windows / macOS
docker run --rm -it -e PULSE_SERVER=tcp:host.docker.internal:4713 golosok

# Тестовый запуск с фейковым STT
docker run --rm golosok -stt-test 1
```

```bash
#mkdir -p build
#CID=$(docker create golosok)
#docker cp "$CID":/opt/golosok/build ./build/vosk
#docker rm "$CID"
```