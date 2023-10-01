package token

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"
)

// 这是我们的数据结构
// 一个诗句包含一个标题和一个内容
type Poem struct {
	Title   string
	Author  string
	Content string
}

func (p Poem) Contains(keyword string) bool {
	return strings.Contains(p.Title, keyword) ||
		strings.Contains(p.Author, keyword) ||
		strings.Contains(p.Content, keyword)
}

//go:embed embed/*
var fs embed.FS

// 这是我们的模板文件的路径
const templateFile = "embed/index.html"

// 这是我们的模板对象
var tmpl *template.Template

var poems []Poem

func init() {
	var err error
	tmpl = template.New("index.html").Funcs(template.FuncMap{
		"safehtml": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
		"replace": func(s string, r string, content string) string {
			return strings.ReplaceAll(content, s, r)
		},
	})
	tmpl, err = tmpl.ParseFS(fs, templateFile)
	if err != nil {
		panic(err)
	}
	poems = loadPoems()
}

func loadPoems() []Poem {
	const poemsDir = "embed/poems"
	entries, err := fs.ReadDir(poemsDir)
	if err != nil {
		panic(err)
	}
	var poems []Poem
	for _, entry := range entries {
		poem, err := fs.ReadFile(path.Join(poemsDir, entry.Name()))
		if err != nil {
			panic(err)
		}
		poems = append(poems, parsePoem(string(poem)))
	}
	return poems
}

func parsePoem(poem string) Poem {
	lines := strings.Split(poem, "\n")
	title := lines[0]
	author := lines[2]
	content := strings.Join(lines[4:], "\n")
	return Poem{
		Title:   title,
		Author:  author,
		Content: content,
	}
}

// handler 是我们的处理函数
// 它必须接收两个参数：一个 ResponseWriter 和一个指向 Request 的指针
func Handler(w http.ResponseWriter, r *http.Request) {
	// 我们先从请求中获取用户输入的字或词
	// 我们假设用户输入的参数名是 q
	q := r.FormValue("q")

	// 我们用一个切片来存储搜索结果
	// 切片是一种动态数组，可以根据需要增长或缩小
	// 我们先创建一个空的切片，类型是 Poem
	var results []Poem

	// // 我们用一个 for 循环来遍历我们收集的古诗词
	// // 我们假设我们的古诗词是一个 map，键是标题，值是内容
	// // 这里我们只用了一些示例数据，你可以用你自己的数据替换它
	// poems := map[string]string{
	// 	"春晓":  "春眠不觉晓，处处闻啼鸟。\n夜来风雨声，花落知多少。",
	// 	"静夜思": "床前明月光，疑是地上霜。\n举头望明月，低头思故乡。",
	// 	"鹿柴":  "空山不见人，但闻人语响。\n返景入深林，复照青苔上。",
	// 	"相思":  "红豆生南国，春来发几枝。\n愿君多采撷，此物最相思。",
	// 	"江雪":  "千山鸟飞绝，万径人踪灭。\n孤舟蓑笠翁，独钓寒江雪。",
	// }

	// 我们用 range 来遍历 map，得到每一个键值对
	for _, poem := range poems {
		// 我们用 Contains 方法来判断内容是否包含用户输入的字或词
		// 如果包含，我们就把这个诗句加入到结果切片中
		if poem.Contains(q) {
			results = append(results, poem)
		}
	}

	// 我们用 Execute 方法来执行模板，把结果写入到响应中
	// 我们把用户输入的字或词和搜索结果作为数据传递给模板
	err := tmpl.Execute(w, map[string]interface{}{
		"Query":   q,
		"Results": results,
	})
	if err != nil {
		// 如果执行出错，我们就打印错误信息
		fmt.Println("Error executing template:", err)
	}
}
