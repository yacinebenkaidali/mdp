package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := parseContent(data, "")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Errorf("Result content does not match golden file")
	}
}

func TestRunWithFile(t *testing.T) {
	var mockStdOut bytes.Buffer
	c := config{
		fileName:    inputFile,
		skipPreview: true,
		tfname:      "",
		out:         &mockStdOut,
		in:          nil,
	}
	if err := run(c); err != nil {
		t.Fatal(err)
	}
	fileName := strings.TrimSpace(mockStdOut.String())
	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Errorf("Result content does not match golden file")
	}

	os.Remove(fileName)
}

func TestRunWithStdIn(t *testing.T) {
	var mockStdOut bytes.Buffer
	var mockStdIn bytes.Buffer

	data, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}
	mockStdIn.Write(data)

	c := config{
		fileName:    "",
		skipPreview: true,
		tfname:      "",
		out:         &mockStdOut,
		in:          &mockStdIn,
	}

	if err := run(c); err != nil {
		t.Fatal(err)
	}
	fileName := strings.TrimSpace(mockStdOut.String())
	result, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result, expected) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Errorf("Result content does not match golden file")
	}

	os.Remove(fileName)
}
