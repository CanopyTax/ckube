package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
)

type OutputManager struct {
	sync.RWMutex
	output        []string
	HeaderColumns []string
}

func (o *OutputManager) Append(s ...string) {
	o.Lock()
	defer o.Unlock()

	o.output = append(o.output, s...)
}

func (o *OutputManager) GetOutput() []string {
	return o.output
}

func (o *OutputManager) HeaderLine() string {
	var headerLine string
	for _, s := range o.HeaderColumns {
		headerLine += s + "\t"
	}
	return headerLine
}

func (o *OutputManager) Print() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.StripEscape)
	fmt.Fprintln(w, o.HeaderLine())
	for _, line := range o.output {
		fmt.Fprintln(w, o.tabbedString(line))
	}
	w.Flush()
}

func (o *OutputManager) FormattedStringSlice() []string {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		fmt.Printf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	w := tabwriter.NewWriter(tempFile, 0, 0, 1, ' ', tabwriter.StripEscape)
	var headerLine string
	for _, s := range o.HeaderColumns {
		headerLine += s + "\t"
	}
	fmt.Fprintln(w, headerLine)
	for _, line := range o.output {
		fmt.Fprintln(w, o.tabbedString(line))
	}
	w.Flush()

	var tabbedLines []string
	// read the temp file into a []string
	formattedFile, err := os.Open(tempFile.Name())
	defer formattedFile.Close()
	if err != nil {
		fmt.Printf("error opening temp file: %v", err)
	}
	scanner := bufio.NewScanner(formattedFile)
	for scanner.Scan() {
		tabbedLines = append(tabbedLines, scanner.Text())
	}
	return tabbedLines
}

func (o *OutputManager) tabbedString(output string) string {
	splits := strings.Split(output, " ")
	var tabbedString string
	for _, s := range splits {
		if s != "" {
			tabbedString += s + "\t"
		}
	}
	return tabbedString
}
