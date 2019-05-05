# musig :speaker:

[![GoDoc](https://godoc.org/github.com/sfluor/musig?status.svg)](https://godoc.org/github.com/sfluor/musig)
[![CircleCI](https://circleci.com/gh/sfluor/musig/tree/master.svg?style=svg)](https://circleci.com/gh/sfluor/musig/tree/master)

A shazam-like tool that allows you to compute song's fingerprints and reverse lookup song names.

## Installation

You will need to have [go](https://golang.org/doc/install) on your computer (version > 1.11 to be able to use go modules).

You will also need to have [portaudio](http://www.portaudio.com/) installed (`brew install portaudio`, or `apt install portaudio19-dev`)

To build the binary:

```bash
git clone git@github.com:sfluor/musig.git
cd musig
make
```

You will then be able to run the binary with:

`./bin/musig help`

## Usage

![musig usage](./docs/musig.gif)

To do some testing you can download `wav` songs by doing `make download`.

Load them with `./bin/musig load "./assets/dataset/wav/*.wav"`

And try to find one of your song name with:

`./bin/musig read "$(ls ./assets/dataset/wav/*.wav | head -n 1)"`

You can also try to use it with your microphone using the `listen` command:

`./bin/musig listen`

For more details on the usage see the help command:

```
A shazam like CLI tool

Usage:
  musig [command]

Available Commands:
  help        Help about any command
  listen      listen will record the microphone input and try to find a matching song from the database (Ctrl-C will stop the recording)
  load        Load loads all the audio files matching the provided glob into the database (TODO: only .wav are supported for now)
  read        Read reads the given audio file trying to find it's song name
  spectrogram spectrogram generate a spectrogram image for the given audio file in png (TODO: only .wav are supported for now)

Flags:
      --database string   database file to use (default "/tmp/musig.bolt")
  -h, --help              help for musig

Use "musig [command] --help" for more information about a command.
```

## Testing

To run the tests you can use `make test` in the root directory.

## TODOs

- [ ] improve the documentation
- [ ] support for `mp3` files
- [ ] improve the fingerprinting (using listen gives bad results :sad:)
