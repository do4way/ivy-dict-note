package evernote

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"

	eClient "github.com/do4way/evernote-sdk-golang/client"
	"github.com/do4way/evernote-sdk-golang/edam"
	"github.com/do4way/ivy-dict-note/webster"
)

//Marshaller ...
type Marshaller struct {
	title     string
	noteBook  *string
	tags      []string
	words     []string
	resources []*edam.Resource
	noteStore *edam.NoteStoreClient
	authToken string
}

//NewMarshaller ...
func NewMarshaller(authToken string) *Marshaller {
	mar := &Marshaller{
		words:     make([]string, 0, 10),
		resources: make([]*edam.Resource, 0, 10),
	}
	eclient := eClient.NewClient("", "", eClient.PRODUCTION)
	eNoteStore, err := eclient.GetNoteStore(authToken)
	if err != nil {
		panic(err)
	}
	mar.noteStore = eNoteStore
	mar.authToken = authToken
	return mar
}

//WriteHeader ...
func (m *Marshaller) WriteHeader(noteBookName string, noteTitle string, tags []string) error {
	m.title = noteTitle
	m.tags = tags
	// nbs, err := m.noteStore.ListNotebooks(m.authToken)
	// fmt.Printf("--------->%v,,%v\n", nbs, err)
	// if err != nil {
	// 	return nil
	// }
	// for _, nb := range nbs {
	// 	fmt.Printf("--------->%s\n", *nb.Name)
	// 	if *nb.Name == noteBook {
	// 		guid := string(*nb.GUID)
	// 		m.noteBook = &guid
	// 		return nil
	// 	}
	// }
	return nil
}

var audioMime = "audio/mpeg"

const (
	wordHeaderTemplate = `<h1>%s&nbsp;&nbsp;<i>[%s]</i></h1>
		<h3>%s</h3>
	`
	audioTemplate = `<en-media type="audio/mpeg" hash="%s"/>`

	difTemplate = `<li style="margin-left:-20px"><span style="font-size:16pt;">%s</span></li>`

	hr = "<hr/>"
)

//WriteWord ...
func (m *Marshaller) WriteWord(word *webster.Word) error {
	wordENML := fmt.Sprintf(wordHeaderTemplate, word.DictForm, word.FormClass, word.PronunSignal)
	if word.Pronunciation != "" {
		resource, err := createResourceFromFile(word.Pronunciation + ".mp3")
		if err == nil {
			m.resources = append(m.resources, resource)
			wordENML += fmt.Sprintf(audioTemplate, fmt.Sprintf("%x", resource.Data.BodyHash))
		}
	}
	wordENML += "<ol>"
	if len(word.Definitions) > 0 && len(word.Definitions[0].Senses) > 0 {
		wordENML += fmt.Sprintf(difTemplate, word.Definitions[0].Senses[0].DefText)
	}

	if len(word.Definitions) > 1 && len(word.Definitions[1].Senses) > 0 {
		wordENML += fmt.Sprintf(difTemplate, word.Definitions[1].Senses[0].DefText)
	}

	wordENML += "</ol>"
	wordENML += hr

	m.words = append(m.words, wordENML)
	return nil
}

//Sync ...
func (m *Marshaller) Sync() error {

	nBody := `<?xml version="1.0" encoding="UTF-8" ?>
		<!DOCTYPE en-note SYSTEM "http://xml.evernote.com/pub/enml2.dtd">
		<en-note>
	`
	for _, word := range m.words {
		nBody += word
	}

	nBody += "</en-note>"
	ourNote := edam.NewNote()
	ourNote.Title = &m.title
	ourNote.Content = &nBody
	ourNote.TagNames = m.tags
	ourNote.Resources = m.resources
	if m.noteBook != nil {
		ourNote.NotebookGuid = m.noteBook
	}
	// fmt.Println(nBody)
	// return nil
	// println(m.authToken)
	// notebook, err := m.noteStore.GetDefaultNotebook(m.authToken)
	// if err != nil {
	// 	panic(err)
	// } else if notebook == nil {
	// 	panic("Invalid Notebook")
	// }
	// println("Default notebook: ", notebook.GetName())
	// return err
	_, err := m.noteStore.CreateNote(m.authToken, ourNote)

	return err
}

func createResourceFromFile(file string) (*edam.Resource, error) {

	fmt.Println(file)
	fdata, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	resource := edam.NewResource()
	attrs := edam.NewResourceAttributes()
	attrs.FileName = &file
	resource.Attributes = attrs
	resource.Data = edam.NewData()
	resource.Mime = &audioMime
	resource.Data.Body = fdata
	sum := md5.Sum(fdata)
	resource.Data.BodyHash = sum[:]
	return resource, nil
}
