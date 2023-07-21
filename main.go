package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/avast/apkparser"
	flag "github.com/ogier/pflag"
)

// flags
var (
	path string
)

func init() {
	flag.StringVarP(&path, "filepath", "p", "", "File path for an apk")
}

func main() {
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Searching user(s): %s\n", path)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	zipErr, resErr, manErr := apkparser.ParseApk(path, enc)
	if zipErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open the APK: %s", zipErr.Error())
		os.Exit(1)
		return
	}

	if resErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse resources: %s", resErr.Error())
	}
	if manErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse AndroidManifest.xml: %s", manErr.Error())
		os.Exit(1)
		return
	}
	fmt.Println()
}
