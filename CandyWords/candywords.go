package main

import (
    "github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gdk"
    "os"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
    // "fmt"
	"unsafe"
)


type App struct {
    window *gtk.Window

    indexLabel *gtk.Label
    searchEntry *gtk.Entry
    searchButton *gtk.ToolButton
    clearButton *gtk.ToolButton

    notebook *gtk.Notebook
    textArea1 *gtk.Label
    textArea2 *gtk.Label
    textArea3 *gtk.Label
}


func (app *App) createToolBar() {
    app.indexLabel = gtk.NewLabel("输入单词✍")
    app.searchEntry = gtk.NewEntry()
    app.searchButton = gtk.NewToolButtonFromStock(gtk.STOCK_FIND)
    app.clearButton = gtk.NewToolButtonFromStock(gtk.STOCK_CLEAR)
}


func (app *App) createWidget() {
    app.notebook = gtk.NewNotebook()

    page1 := gtk.NewFrame("有道词典")
    app.notebook.AppendPage(page1, gtk.NewLabel("有道词典"))
    app.textArea1 = gtk.NewLabel("")
    // app.textArea1.SetSelectable(true)
    textContainer1 := gtk.NewFixed()
    textContainer1.Put(app.textArea1, 20, 20)
    // textContainer.Add(textArea)
    page1.Add(textContainer1)


    page2 := gtk.NewFrame("金山词霸")
    app.notebook.AppendPage(page2, gtk.NewLabel("金山词霸"))
    app.textArea2 = gtk.NewLabel("")
    // app.textArea2.SetSelectable(true)
    textContainer2 := gtk.NewFixed()
    textContainer2.Put(app.textArea2, 20, 20)
    page2.Add(textContainer2)


    page3 := gtk.NewFrame("牛津词典")
    app.notebook.AppendPage(page3, gtk.NewLabel("牛津词典"))
    app.textArea3 = gtk.NewLabel("")
    // app.textArea3.SetSelectable(true)
    textContainer3 := gtk.NewFixed()
    textContainer3.Put(app.textArea3, 20, 20)
    page3.Add(textContainer3)
}


func (app *App) setSize() {
    app.indexLabel.SetSizeRequest(100, 40)
    app.searchEntry.SetSizeRequest(300, 30)
    app.searchButton.SetSizeRequest(40, 40)
    app.clearButton.SetSizeRequest(40, 40)
    app.notebook.SetSizeRequest(600, 360)
}


func (app *App) setLayout(mainLayout *gtk.Fixed) {
    mainLayout.Put(app.indexLabel, 0, 0)
    mainLayout.Put(app.searchEntry, 100, 5)
    mainLayout.Put(app.searchButton, 400, 0)
    mainLayout.Put(app.clearButton, 440, 0)
    mainLayout.Put(app.notebook, 0, 40)
}


func (app *App) windowBindButtonConnection() {
    app.searchButton.Connect("clicked", app.showWordMeans)

    app.clearButton.Connect("clicked", func() {
        app.searchEntry.SetText("")
    })
}

func (app *App) showWordMeans() {
    words := app.searchEntry.GetText()
    if len(words) > 0 {
        means := searchYouDao(words)
        // fmt.Println(means)
        app.textArea1.SetText(means)
        means = searchJinShan(words)
        app.textArea2.SetText(means)
        // means = searchOxford(words)
        app.textArea3.SetText("敬请期待...")
    }
}


func formatStrings(strs []string) string {
    var res []string
    // var rune []rune
    for i, v := range strs {
        if i == 0 {
            res = append(res, v)
        } else {
            rune := []rune(v)
            if len(rune) <= 38 {
                res = append(res, "        " + v)
            } else {
                res = append(res, "        " + string(rune[:38]))
                rune = rune[38:]
                for len(rune) > 38 {
                    res = append(res, string(rune[:38]))
                    rune = rune[38:]
                }
                if len(rune) > 0 {
                    res = append(res, string(rune))
                }
            }            
        }
    }
    return strings.Join(res, "\n")
}


func searchYouDao(words string) string {
    url := "http://www.youdao.com/w/eng/" + words + "/#keyfrom=dict2.index"
    resp, err := http.Get(url)
    if err != nil {
        return "发生了网络错误......"
    }
    content, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    re := regexp.MustCompile(`<li>(.*?)</li>`)
    search_result := re.FindAllSubmatch(content[:len(content)/2], -1)
    if len(search_result) == 0 {
        return "没查到......"
    }
    var means []string = []string{words + ":"}
    // var rune []rune
    for _, e := range search_result {
        // mean = string(e[1])
        means = append(means, string(e[1]))
    }
    return formatStrings(means)
}


func searchJinShan(words string) string {
    url := "http://www.iciba.com/" + words
    resp, err := http.Get(url)
    if err != nil {
        return "发生了网络错误......"
    }
    content, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    re := regexp.MustCompile(`<li class="clearfix">([\s\S]*?)</li>`)
    result := re.FindAll(content, -1)
    if len(result) == 0 {
        return "没查到......"
    }

    re_prop := regexp.MustCompile(`<span class="prop">(.*?)</span>`)
    re_span := regexp.MustCompile(`<span>(.*?)</span>`)
    var means []string = []string{words + ":"}
    var mean string
    for _, i := range result {
        mean = string(re_prop.FindSubmatch(i)[1])
        spans := re_span.FindAllSubmatch(i, -1)
        for _, j := range spans {
            mean += string(j[1])
        }
        means = append(means, mean)
    }
    return formatStrings(means)
}


func searchOxford(words string) string {
    url := "http://www.icooc.com/service/oxford/search/?word=h"
    resp, err := http.Get(url)
    if err != nil {
        return "发生了网络错误......"
    }
    content, _ := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    re := regexp.MustCompile(`<span class="def".*?>(.*?)</span>`)
    result := re.FindAllSubmatch(content, -1)
    var means []string = []string{words + ":"}
    for _, i := range result {
        means = append(means, string(i[1]))
    }
    return strings.Join(means, "\n")
}


func (app *App) initWindow() {
    gtk.Init(&os.Args)
    app.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
    app.window.SetTitle("🎉🍡Candy Words 1.0")
    app.window.SetResizable(false)
    app.window.Connect("destroy", gtk.MainQuit)

    app.createToolBar()
    app.createWidget()
    app.setSize()
    mainLayout := gtk.NewFixed()
    app.setLayout(mainLayout)
    app.windowBindButtonConnection()
	app.windowBindKeyEvent()


    app.window.Add(mainLayout)
    app.window.ShowAll()
    app.window.SetSizeRequest(600, 400)

    gtk.Main()
}


func (app *App) windowBindKeyEvent() {
	app.window.Connect("key-press-event", func(ctx *glib.CallbackContext){
		arg := ctx.Args(0)
		event := *(**gdk.EventKey)(unsafe.Pointer(&arg))
		key := event.Keyval
		if key == 65293 {
			// fmt.Println("yes")
            app.showWordMeans()
		}
        if key == 65365 || key == 65366{
            app.searchEntry.SetText("")
        }
		// fmt.Println(key)
	})
}


func main() {
    app := App{}
    app.initWindow()
}
