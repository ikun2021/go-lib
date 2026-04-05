package opus

import (
	"bytes"
	"errors"
	"github.com/pion/opus"
	"github.com/pion/opus/pkg/oggreader"
	"io"
	"os"
	"testing"
)

func TestOgger(t *testing.T) {

	file, err := os.Open("C:\\Users\\dengyongcai\\Downloads\\seE7JSZkxWrq.opus")
	if err != nil {
		panic(err)
	}

	ogg, _, err := oggreader.NewWith(file)
	if err != nil {
		panic(err)
	}

	out := make([]byte, 1920)
	fd, err := os.Create("C:\\Users\\dengyongcai\\Downloads\\dest.pcm")
	if err != nil {
		panic(err)
	}

	decoder := opus.NewDecoder()
	for {
		segments, _, err := ogg.ParseNextPage()

		if errors.Is(err, io.EOF) {
			break
		} else if bytes.HasPrefix(segments[0], []byte("OpusTags")) {
			continue
		}

		if err != nil {
			panic(err)
		}

		for i := range segments {
			if _, _, err = decoder.Decode(segments[i], out); err != nil {
				panic(err)
			}

			if _, err := fd.Write(out); err != nil {
				panic(err)
			}
		}
	}
}
