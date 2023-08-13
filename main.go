package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// const (
//
//		header = `<!DOCTYPE html>
//	  <html>
//		<head>
//		  <meta http-equiv="content-type" content="text/html; charset=utf-8">
//		  <title>Markdown Preview Tool</title>
//		</head>
//		<body>
//	  `
//		footer = `
//		</body>
//	  </html>
//	  `
//
// )
const (
	defaultTemplate = `<!DOCTYPE html>
  <html>
	<head>
	  <meta http-equiv="content-type" content="text/html; charset=utf-8">
	  <title>{{ .Title }}</title>
	</head>
	<body>
  {{ .Body }}
	</body>
  </html>
  `
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	preview := flag.Bool("skip", false, "skip auto-preview")
	templateName := flag.String("t", "", "Alternate template name")
	flag.Parse()
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, !(*preview), *templateName, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
func run(fileName string, skipPreview bool, tfname string, out io.Writer) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	htmlContent, err := parseContent(data, tfname)
	if err != nil {
		return err
	}
	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)
	if err := saveHtml(htmlContent, outName); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}
	defer os.Remove(outName)
	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	var buffer bytes.Buffer

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}
	unsafe := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	c := content{
		Title: "Mardown preview",
		Body:  template.HTML(body),
	}
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func saveHtml(data []byte, outputName string) error {
	return os.WriteFile(outputName, data, 0644)
}
func preview(fname string) error {
	cName := ""
	cParams := []string{}
	switch runtime.GOOS {
	case "linux":
		{
			cName = "xdg-open"
		}
	case "windows":
		{
			cName = "cmd.exe"
			cParams = []string{"/C", "start"}
		}
	case "darwin":
		{
			cName = "open"
		}
	default:
		{
			return fmt.Errorf("OS not supported")
		}
	}
	cParams = append(cParams, fname)
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}
	err = exec.Command(cPath, cParams...).Run()
	time.Sleep(time.Second * 2)
	return err
}
