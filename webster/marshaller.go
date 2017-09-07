package webster

import (
	"fmt"
	"io"
)

//Marshaller ...
type Marshaller interface {
	WriteHeader(notebook string, noteTitle string, tags []string) error
	WriteWord(word *Word) error
}

//MdMashaller ...
type MdMashaller struct {
	Out io.Writer
}

//WriteHeader ...
func (m *MdMashaller) WriteHeader(notebook string, noteTitle string, tags []string) error {
	if noteTitle != "" {
		fmt.Fprintf(m.Out, "#%s\n", noteTitle)
	}
	var tagStr = "@(" + notebook + ")["
	for _, tag := range tags {
		tagStr += tag + " | "
	}
	tagStr += "]"
	fmt.Fprintln(m.Out, tagStr)
	return nil
}

//WriteWord ...
func (m *MdMashaller) WriteWord(word *Word) error {
	fmt.Fprintf(m.Out, "# %s *[%s]*\n", word.DictForm, word.FormClass)
	fmt.Fprintf(m.Out, "#### *%s* ", word.PronunSignal)
	if word.Pronunciation != "" {
		fmt.Fprintf(m.Out, "*%s.mp3*\n", word.Pronunciation)
	} else {
		fmt.Print("\n")
	}

	if len(word.Definitions) > 0 && len(word.Definitions[0].Senses) > 0 {
		s := word.Definitions[0].Senses[0]
		fmt.Fprintf(m.Out, "- __%s__\n", s.DefText)
	}

	if len(word.Definitions) > 1 && len(word.Definitions[1].Senses) > 0 {
		s := word.Definitions[1].Senses[0]
		fmt.Fprintf(m.Out, "- __%s__\n", s.DefText)
	}

	// for i, d := range word.Definitions {
	// 	fmt.Printf("%d. ", i)
	// 	for _, s := range d.Senses {
	// 		fmt.Printf("- *%s*: %s\n", s.Sgram, s.DefText)
	// 	}
	// }
	fmt.Fprintln(m.Out, "***")
	// fmt.Printf("%s\n", word.Definitions)
	return nil
}
