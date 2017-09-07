package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/do4way/ivy-dict-note/evernote"
	"github.com/do4way/ivy-dict-note/webster"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "array string flags"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var noteBookName string
var noteName string
var tags arrayFlags
var devToken string

func init() {
	const (
		defaultNoteBook = ""
		usageNoteBook   = "The note book name"
		defaultNoteName = ""
		usageNoteName   = "The note name"
	)
	flag.StringVar(&noteBookName, "note_book", "", "The note book name")
	flag.StringVar(&noteBookName, "nb", "", "The note book name")
	flag.StringVar(&noteName, "note", "", "The note title")
	flag.StringVar(&noteName, "n", "", "The note title")
	flag.StringVar(&devToken, "token", "", "The evernote dev token")
	flag.Var(&tags, "t", "The note tag")
}

func fetchAudio(word *webster.Word, wg *sync.WaitGroup) {
	webster.FetchWordAudio(word, nil)
	wg.Done()
}

func exec(strWrds []string, marshaller webster.Marshaller) {
	var wg sync.WaitGroup
	wg.Add(len(strWrds))
	wordsMap := make(map[string][]*webster.Word)
	for _, strWrd := range strWrds {
		go func(w string) {
			words, err := webster.ParseWebster(w)
			if err != nil {
				log.Print(err)
				wg.Done()
				return
			}
			wordsMap[w] = words
			var awg sync.WaitGroup
			awg.Add(len(words))
			for _, word := range words {
				go fetchAudio(word, &awg)
			}
			awg.Wait()
			wg.Done()
		}(strWrd)
	}
	wg.Wait()
	marshaller.WriteHeader(noteBookName, noteName, tags)
	for _, strWrd := range strWrds {
		for _, word := range wordsMap[strWrd] {
			marshaller.WriteWord(word)
		}
	}
}

func main() {

	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(0)
	}
	flag.Parse()
	if devToken != "" {
		marshaller := evernote.NewMarshaller(devToken)
		exec(flag.Args(), marshaller)
		err := marshaller.Sync()
		if err != nil {
			panic(err)
		}
	} else {
		marshaller := &webster.MdMashaller{
			Out: os.Stdout,
		}
		exec(flag.Args(), marshaller)
	}
}
