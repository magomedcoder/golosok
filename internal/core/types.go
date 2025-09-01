package core

type TTSInitFn func(*Core) error

type TTSSayFn func(*Core, string) error

type TTSToFileFn func(*Core, string, string) error

type NormalizerInitFn func(*Core) error

type NormalizeFn func(*Core, string) string

type PlayWAVInitFn func(*Core) error

type PlayWAVFn func(*Core, string) error
