package main

import (
	"os"
	"path/filepath"
	"runtime"

	orcgen "github.com/luabagg/orcgen/pkg"
)

// import "github.com/luabagg/orcgen"

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	file := filepath.Join(basepath, "test.html")
	html, _ := os.ReadFile(file)

	gen := orcgen.New(orcgen.PDF)

	f, _ := gen.ConvertHTML(html)
	f.Output("test.pdf")

	f, _ = gen.SetExt(orcgen.PNG).SetFullPage(true).ConvertHTML(html)
	f.Output("test.png")

	f, _ = gen.SetExt(orcgen.JPEG).ConvertHTML(html)
	f.Output("test.jpeg")
}
