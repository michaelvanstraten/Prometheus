package prometheus

import (
	"encoding/json"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"github.com/michaelvanstraten/prometheus/rendering"
)

type AppConfig struct {
}

type App struct {
	Router httprouter.Router
	Renderer rendering.HtmlRenderer
}

func New() *App {
	return &App{}
}

func NewFromConfig(Path string) *App {
	if file, err := ioutil.ReadFile(Path); err == nil {
		var NewApp = &App{}
		if err := json.Unmarshal([]byte(file), NewApp); err != nil {
			println("Error Parsing the App Configuration file: " + Path)
			return nil
		} else {
			return NewApp
		}
	} else {
		println("Could not Read the Configuration file: " + Path)
		return nil
	}
}