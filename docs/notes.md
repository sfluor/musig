# Analog -> Digital

- Remove frequencies over 20khz before sampling to avoid aliasing (otherwise those frequencies will end up between 0Hz and 20khz)

# Mp3 -> headphones

- Uses PCM (Pulse coded modulation), a PCM stream is a stream of bits that can be composed of multiple channels, for instance stereo has 2 channels.

- In a stream the amplitude of the signal is divided into samples (a 44,1khz sampled music will have 44100 samples per second), each sample gives the (quantized) amplitude of the sound for the corresponding time interval.

# How to get frequencies ?

- Apply an FFT on small subsets of the original signal (0.1 s intervals for instance). This will give us the amplitude and frequency over a time period

# Resampling

- Resampling a song allows to process the song faster, if the song was sampled at a 44,1khz rate (and we were doing a 4096-sample window) and we resample it with a 11.khz rate, we will have frequencies from 0 to 5khz only but this would allow us to have the same frequency resolution when doing a 1024-sample window FFT.

- Having frequencies only from 0 to 5kz is not really an issue with songs since the most important frequencies for musics are within this range.

- An easy way of doing resampling is by doing an average of the signal, (for example 4 samples averaged)

- However we also have to be careful about aliasing, we have to apply a low pass filter to avoid it since we now have frequencies from 0kz to 5khz.

# Stero to mono

We also have to transform stereo to mono

# Some sump up

We can have a spectrogram:
- From 0Hz to 5kHz
- With a bin size of 10.7 Hz

[comment]: # (TODO figure out this)
- 512 possible frequencies (here a "frequency" is in fact a bin, this comes from the DFT, try to understand this)
- And a unit of time of 0.1 second

# Filtering

- Since we want to be noise tolerant we have to keep only the loudest notes, but we can't just do that on every 0.1 second bin:
- - Since human ears have more difficulties to hear some low freqs, their amplitude is sometimes increased , if you only take the most powerful frequencies, you will end up with only the low ones
- - Spectral leakage over a powerful frequency will create other powerful frequencies

The solution:

- For each FFT result we put the 512 frequency bins inside 6 logarithmic bands (we could use another number of bands here):
- - [0, 10]
- - [10, 20]
- - [20, 40]
- - [40, 80]
- - [80, 160]
- - [160, 511]

- For each band we keep the strongest bin of frequencies
- Compute the average value of these 6 maximums
- Keep only the bins above this mean (multiplied by a coefficient)
- This can have some caveats though (search `this algorithm has a limitation` [here](http://coding-geek.com/how-shazam-works/)), also search `a simple algorithm to filter the spectrogram` for another way of doing

- We then end up with a spectrogram but without the amplitude, only points (time, frequency)

