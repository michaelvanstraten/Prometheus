package rendering

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"html/template"
	"sync"
)

type database interface {
	Get(Key string) []byte
	Set(Key string, Data []byte,)
}

type HtmlRenderer struct {
	DatabaseConnections []database
	Main template.Template
	mux sync.Mutex
}

func (rd *HtmlRenderer) Render(TemplateName string, Data interface{}) []byte {
	var c = make(chan []byte, 2)
	go func() {
		var renderbuffer *bytes.Buffer = bytes.NewBuffer([]byte{})
		rd.Main.ExecuteTemplate(renderbuffer, TemplateName, Data)
		go func (){
			for index := range rd.DatabaseConnections {
				rd.DatabaseConnections[index].Set(HashData(&Data), renderbuffer.Bytes())
			}
		}()
		c <- renderbuffer.Bytes()
	}()
	for index := range rd.DatabaseConnections {
		go func(db database) {
			if data := db.Get(HashData(&Data)); data != nil {
				c <- data
			}
		}(rd.DatabaseConnections[index])
	}
	return <- c
}

func (rd *HtmlRenderer) AddTemplate(Path ...string) {
	defer rd.mux.Unlock()
	rd.mux.Lock()
	for i := 0; i < len(Path); i++ {
		if _, err := rd.Main.ParseGlob(Path[i]); err != nil{
			println(err.Error())
		}
	}
}

func HashData(Data *interface{}) string {
	if databytes, err := json.Marshal(Data); err == nil {
		var hash = sha256.Sum256(databytes)
		return string(hash[:])
	}
	return ""
}

func New() *HtmlRenderer {
	var rd = HtmlRenderer{}
	rd.Main = *template.New("Main")
	rd.mux = sync.Mutex{}
	return &rd
}