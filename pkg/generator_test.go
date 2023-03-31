package pkg

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	PageIdleTime = 700 * time.Millisecond

	gen := New(JPEG)
	defer gen.Close()

	err := gen.ConvertWebpage("https://www.google.com")
	if err != nil {
		log.Fatalf("err converting file: %v", err)
	}

	err = gen.GetOutput("../tmp/test1.jpeg")
	if err != nil {
		log.Fatalf("err creating file: %v", err)
	}

	gen = New(PDF)

	filepath := "../test/test.html"
	b, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal("err reading file")
	}
	err = gen.ConvertHTML(b)
	if err != nil {
		log.Fatalf("err converting file: %v", err)
	}

	err = gen.GetOutput("../tmp/test1.pdf")
	if err != nil {
		log.Fatalf("err creating file: %v", err)
	}

	log.Print("success")
}
