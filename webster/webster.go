package webster

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

// Word ..
type Word struct {
	DictForm      string
	PronunSignal  string
	FormClass     string
	Forms         string
	Definitions   []*Definition
	Phrases       []*Phrase
	Pronunciation string
}

//Definition ...
type Definition struct {
	Senses []*Sense
}

func (d *Definition) String() string {
	return fmt.Sprintf("%v", d.Senses)
}

//Sense ...
type Sense struct {
	Sgram    string
	Note     string
	DefText  string
	Examples []*Example
}

func (s *Sense) String() string {
	return fmt.Sprintf("{Sgram: %s,\nNote: %s,\nDefText: %s,\n Examples: %s} \n", s.Sgram, s.Note, s.DefText, s.Examples)
}

// Phrase ...
type Phrase struct {
	Text        string
	Definitions []*Definition
}

func (p *Phrase) String() string {
	return fmt.Sprintf("{Text: %s,\nDefinitions %s}\n", p.Text, p.Definitions)
}

// Example ...
type Example struct {
	ExtraText string
	Text      string
}

func (e *Example) String() string {
	return fmt.Sprintf("{Text:%s, \n ExtraText: %s} \n", e.Text, e.ExtraText)
}

func parseExamples(s *goquery.Selection) []*Example {
	examples := make([]*Example, 0, 10)
	s.Each(func(i int, sb *goquery.Selection) {
		examples = append(examples, &Example{
			ExtraText: "",
			Text:      cleanTrailingWhiteSpace(sb.Text()),
		})
	})
	return examples
}

func parseSenses(s *goquery.Selection) []*Sense {
	senses := make([]*Sense, 0, 3)
	s.Each(func(i int, sb *goquery.Selection) {
		def := strings.Join(sb.Find(".def_text").Map(func(i int, s *goquery.Selection) string {
			return cleanTrailingWhiteSpace(s.Text())
		}), " : ")
		senses = append(senses, &Sense{
			Sgram:    cleanTrailingWhiteSpace(sb.Find(".sgram").Text()),
			DefText:  def,
			Examples: parseExamples(sb.Find(".vi_content")),
		})
	})
	return senses
}

func parseDefinition(s *goquery.Selection) []*Definition {
	definitions := make([]*Definition, 0, 10)
	s.Each(func(i int, sb *goquery.Selection) {
		definitions = append(definitions, &Definition{
			Senses: parseSenses(sb.Find(".sense")),
		})
	})
	return definitions
}

func parsePhrase(s *goquery.Selection) []*Phrase {
	phrases := make([]*Phrase, 0, 5)
	s.Each(func(i int, sb *goquery.Selection) {
		phrases = append(phrases, &Phrase{
			Text:        sb.Find("h2.dre").Text(),
			Definitions: parseDefinition(sb.Find(".sblock.sblock_dro")),
		})
	})
	return phrases
}

func cleanTrailingWhiteSpace(s string) string {
	var rst []string
	for _, v := range strings.Split(s, "\n") {
		t := strings.TrimSpace(v)
		if t != "" {
			rst = append(rst, t)
		}
	}
	return strings.Join(rst, "\n")
}

const dictPrefix = "http://learnersdictionary.com/definition/"

//ParseWebster ...
func ParseWebster(word string) ([]*Word, error) {

	doc, err := goquery.NewDocument(dictPrefix + word)
	if err != nil {
		return nil, err
	}
	wrds := make([]*Word, 0, 5)
	doc.Find(".entry").Each(func(i int, s *goquery.Selection) {

		d := strings.TrimLeftFunc(cleanTrailingWhiteSpace(s.Find(".hw_txt.gfont").Text()), func(r rune) bool {
			return unicode.IsDigit(r) || unicode.IsSpace(r)
		})
		if d == "" {
			return
		}
		word := &Word{DictForm: d}
		p := cleanTrailingWhiteSpace(s.Find(".hw_d>.hpron_word.ifont").Text())
		// if p == "" {
		// 	p = cleanTrailingWhiteSpace(s.Find(".hpron_word.ifont").Text())
		// }
		if p == "" {
			p = cleanTrailingWhiteSpace(s.Find(".uro_line>.pron_w.ifont").Text())
		}
		if p == "" {
			p = cleanTrailingWhiteSpace(s.Find(".hw_vars_d>.pron_w.ifont").Text())
		}
		word.PronunSignal = p
		word.Pronunciation = s.Find("a.fa").AttrOr("data-file", "")
		word.FormClass = cleanTrailingWhiteSpace(s.Find(".hw_d>.fl").Text())
		word.Forms = cleanTrailingWhiteSpace(s.Find(".hw_infs_d").Text())
		word.Definitions = parseDefinition(s.Find(".sblock.sblock_entry"))
		word.Phrases = parsePhrase(s.Find(".dro"))
		wrds = append(wrds, word)
	})
	return wrds, nil
}
