package main

import (
	"time"

	"github.com/luabagg/orcgen/internal"
	orcgen "github.com/luabagg/orcgen/pkg"
)

// import "github.com/luabagg/orcgen"

func main() {
	gen := orcgen.New(orcgen.PDF)
	f, _ := gen.ConvertWebpage("https://www.google.com")
	f.Output("google.pdf")

	f, _ = gen.SetExt(orcgen.PNG).ConvertWebpage("https://www.github.com")
	f.Output("test.png")

	go func(d *internal.Director) {
		f, _ = gen.SetExt(orcgen.JPEG).SetFullPage(false).SetPageIdleTime(1500 * time.Millisecond).ConvertWebpage("https://www.twitter.com")
		f.Output("test.jpeg")
	}(gen)
}
