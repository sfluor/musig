package sound

// Reader interface allows you to read audio files
type Reader interface {
	Read() ([]int, error)
	SampleRate() uint32
}
