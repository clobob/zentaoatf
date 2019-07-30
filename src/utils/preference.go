package utils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/fatih/color"
	"github.com/jroimartin/gocui"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var Prefer model.Preference

func InitPreference() {
	// preference from yaml
	Prefer = getInst()

	// screen size
	InitScreenSize()

	// internationalization
	InitI118(Prefer.Language)

	if strings.Index(os.Args[0], "atf") > -1 && (len(os.Args) > 1 && os.Args[1] != "set") {
		PrintPreference()
	}
}

func SetPreference(param string, val string, dumb bool) {
	buf, _ := ioutil.ReadFile(PreferenceFile)
	yaml.Unmarshal(buf, &Prefer)

	if param == "lang" {
		Prefer.Language = val
		if !dumb {
			color.Blue(I118Prt.Sprintf("set_preference", I118Prt.Sprintf("lang"), I118Prt.Sprintf(Prefer.Language)))
		}
	} else if param == "workDir" {
		val = convertWorkDir(val)

		Prefer.WorkDir = val
		updateWorkDirHistory()
		if !dumb {
			color.Blue(I118Prt.Sprintf("set_preference", I118Prt.Sprintf("workDir"), Prefer.WorkDir))
		}
	}
	data, _ := yaml.Marshal(&Prefer)
	ioutil.WriteFile(PreferenceFile, data, 0666)
}

func SaveProjectHistory(workDir string) {
	buf, _ := ioutil.ReadFile(workDir + ConfigFile)
	yaml.Unmarshal(buf, &Prefer)

	data, _ := yaml.Marshal(&Prefer)
	ioutil.WriteFile(PreferenceFile, data, 0666)
}

func getInst() model.Preference {
	var once sync.Once
	once.Do(func() {
		Prefer = model.Preference{}
		if FileExist(PreferenceFile) {
			buf, _ := ioutil.ReadFile(PreferenceFile)
			yaml.Unmarshal(buf, &Prefer)
		} else { // init
			Prefer.Language = "en"
			Prefer.WorkDir = convertWorkDir("./")
			Prefer.WorkHistories = []string{Prefer.WorkDir}

			data, _ := yaml.Marshal(&Prefer)
			ioutil.WriteFile(PreferenceFile, data, 0666)
		}
	})
	return Prefer
}

func PrintPreference() {
	color.Blue(I118Prt.Sprintf("current_preference", ""))

	val := reflect.ValueOf(Prefer)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Prefer).NumField(); i++ {
		val := val.Field(i)
		fmt.Printf("  %s: %v \n", typeOfS.Field(i).Name, val.Interface())
	}
}

func PrintPreferenceToView(v *gocui.View) {
	fmt.Fprintln(v, color.BlueString(I118Prt.Sprintf("current_preference", "")))

	val := reflect.ValueOf(Prefer)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(Prefer).NumField(); i++ {
		val := val.Field(i)
		fmt.Fprintln(v, fmt.Sprintf("  %s: %v", typeOfS.Field(i).Name, val.Interface()))
	}
}

func convertWorkDir(path string) string {
	if path == "./" || path == "." {
		path, _ = filepath.Abs(`.`)
		path = path + string(os.PathSeparator)
	} else {
		if strings.LastIndex(path, "/") != len(path)-1 {
			path = path + string(os.PathSeparator)
		}
	}

	return path
}

func updateWorkDirHistory() {
	curr := Prefer.WorkDir
	histories := Prefer.WorkHistories

	// 已经是第一个，不做操作
	if histories[0] == curr {
		return
	}

	// 移除元素
	idx := -1
	for i, item := range histories {
		if item == curr {
			idx = i
		}
	}
	if idx > -1 {
		histories = append(histories[:idx], histories[idx+1:]...)
	}

	// 头部插入元素
	histories = append([]string{curr}, histories...)

	// 保存最后10个
	if len(histories) > 10 {
		histories = histories[:10]
	}

	Prefer.WorkHistories = histories
}