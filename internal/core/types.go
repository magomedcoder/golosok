package core

type TTSInitFn func(*Core) error

type TTSSayFn func(*Core, string) error

type NormalizerInitFn func(*Core) error

type NormalizeFn func(*Core, string) string
