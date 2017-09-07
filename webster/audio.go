package webster

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const bytesToMegaBytes = 1048576.0

type passThru struct {
	io.Reader
	curr  int64
	total float64
}

func (pt *passThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.curr += int64(n)

	// last read will have EOF err
	// if err == nil || (err == io.EOF && n > 0) {
	// 	printProgress(float64(pt.curr), pt.total)
	// }

	return n, err
}

func wget(url string, out io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http get error %v", err)
		return err
	}

	defer resp.Body.Close()

	src := &passThru{Reader: resp.Body, total: float64(resp.ContentLength)}

	_, err = io.Copy(out, src)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

var audioPrefix = "http://media.merriam-webster.com/audio/prons/en/us/mp3/"

//FetchWordAudio ...
func FetchWordAudio(word *Word, out io.Writer) {
	f := word.Pronunciation
	if f == "" {
		return
	}
	if out == nil {
		f, _ := os.Create(f + ".mp3")
		defer f.Close()
		out = f
	}
	wget(audioPrefix+string(f[0])+"/"+f+".mp3", out)
}
