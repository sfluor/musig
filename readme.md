# musig :speaker:

[![GoDoc](https://godoc.org/github.com/sfluor/musig?status.svg)](https://godoc.org/github.com/sfluor/musig)
[![CircleCI](https://circleci.com/gh/sfluor/musig/tree/master.svg?style=svg)](https://circleci.com/gh/sfluor/musig/tree/master)

A shazam-like tool that allows you to compute song's fingerprints and reverse lookup song names.

## Installation

You will need to have [go](https://golang.org/doc/install) on your computer (version > 1.11 to be able to use go modules).

To build the binary:

```bash
git clone git@github.com:sfluor/musig.git
cd musig
make
```

You will then be able to run the binary with:

`./bin/musig help`

## Usage

To do some testing you can download `wav` songs and put them in `./assets/dataset/wav/`

Load them with `./bin/musig load "./assets/dataset/wav/*.wav"`

And try to find one of your song name with:

`./bin/musig read "$(ls ./assets/dataset/wav/*.wav | head -n 1)"`

For more details on the usage see the help command:

```
A shazam like CLI tool

Usage:
  musig [command]

Available Commands:
  help        Help about any command
  load        Load loads all the audio files matching the provided glob into the database (TODO: only .wav are supported for now)
  read        Read reads the given audio file trying to find it's song name
  spectrogram spectrogram generate a spectrogram image for the given audio file in png (TODO: only .wav are supported for now)

Flags:
      --database string   database file to use (default "/tmp/musig.bolt")
  -h, --help              help for musig
      --version           version for musig

Use "musig [command] --help" for more information about a command.
```

## Testing

To run the tests you can use `make test` in the root directory.

## TODOs

- [ ] improve the documentation
- [ ] add more audio files to allow to test without having to download them separately
- [ ] support for `mp3` files
- [ ] `listen` to read audio from the mic
