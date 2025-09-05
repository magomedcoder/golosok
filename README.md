# Golosok - офлайн-голосовой ассистент на Go

**Golosok** - это лёгкий и расширяемый голосовой ассистент, который работает полностью локально, без интернета.

Возможности:

- Распознавание речи с микрофона через Vosk (STT)
- Озвучивание ответов голосом (TTS через RHVoice или консоль)
- Нормализация текста (преобразование чисел в слова и очистка от лишних символов)
- Управление таймерами

---

## Сборка

### Linux

```bash
docker build -t golosok-linux .

CID=$(docker create golosok-linux)

docker cp "$CID":/opt/golosok/build ./build/golosok-linux

docker rm "$CID"
```

### Windows

```bash
docker build -f Dockerfile-windows -t golosok-windows .

CID=$(docker create golosok-windows)

docker cp "$CID":/opt/golosok/build ./build/golosok-windows

docker rm "$CID"
```

### Тестовый запуск с фейковым STT

```bash
# Linux
cd build/golosok-linux
LD_LIBRARY_PATH=./lib ./golosok -stt-test 1

# Windows
cd build/golosok-windows
golosok.exe -stt-test 1
```

## Запуск

```bash
# Linux
cd build/golosok-linux
LD_LIBRARY_PATH=./lib ./golosok

# Windows
cd build/golosok-windows
golosok.exe
```

---

## Примеры команд

- голосок привет
- голосок дата
- голосок время
- голосок поставь таймер | поставь таймер
- голосок удали таймер | отмени таймер
- голосок удали все таймеры | сбрось все таймеры | отмени все таймеры
- голосок команды

------------------------------------------------------------------------

## Создание своей команды

1. Создайте папку с именем вашей команды в каталоге `internal/commands/`.

   Например: `internal/commands/example/`.


2. В этой папке можно держать одну или несколько команд.

   Пример файла `example.go`:

   ``` go
   package example

   import (
       "github.com/magomedcoder/golosok/internal/core"
   )

   // Register регистрирует команды
   func Register(c *core.Core) {
       c.RegisterCommand("пример", func(c *core.Core, phrase string) {
           c.Say("Я пример")
       })

       c.RegisterCommand("ещёпример", func(c *core.Core, phrase string) {
           c.Say("Я ещё один пример")
       })
   }
   ```

   > Здесь `"голосок пример"` и `"голосок ещёпример"` - ключевые фразы, которые должен произнести пользователь.
   >
   > `phrase` - остаток распознанного текста после ключа.
   > Например, если сказать «голосок пример тест», то в функцию придёт `"тест"`.


3. Подключите новую команду в `cmd/golosok/main.go`.

   В начале файла (в блоке `import`) добавьте строку:

   ``` go
   "github.com/magomedcoder/golosok/internal/commands/example"
   ```

   А в `main()`, сразу после инициализации ядра (`c := core.NewCore()`) вызовите:

   ``` go
   example.Register(c)
   ```

5. Пересоберите проект, чтобы все команды из папки стали доступны.

### Пример использования

После добавления и регистрации команды: `голосок пример`

Ассистент ответит: `Я пример`
