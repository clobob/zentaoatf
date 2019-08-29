package widget

import (
	"fmt"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
)

var (
	HelpGlobal = ""
)

type HelpWidget struct {
	name string
	x, y int
	w, h int
	body string
}

func NewHelpWidget() {
	initContent()

	maxX, _ := vari.Cui.Size()

	lines := strings.Split(HelpGlobal, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines)
	w = w + 2

	help := HelpWidget{name: "help", x: maxX - w - 3, y: 1, w: w, h: h, body: HelpGlobal}
	help.Layout()
}

func (w *HelpWidget) Layout() error {
	v, err := vari.Cui.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}

func ShowHelpFromView(g *gocui.Gui, v *gocui.View) error {
	return ShowHelp()
}
func ShowHelp() error {
	help, _ := vari.Cui.View("help")

	if help != nil {
		HideHelp()
	} else {
		NewHelpWidget()
	}

	return nil
}

func HideHelp() error {
	help, _ := vari.Cui.View("help")

	if help != nil {
		if err := vari.Cui.DeleteView("help"); err != nil {
			return err
		}
	}

	return nil
}

func initContent() {
	HelpGlobal = fmt.Sprintf("%s \n %s \n",
		i118Utils.I118Prt.Sprintf("help_key_bind"),
		i118Utils.I118Prt.Sprintf("help_exit"))
}
