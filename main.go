package main

import (
	"bytes"
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

// ApkInfo ...
type ApkInfo struct {
	AppName           string
	PackageName       string
	VersionCode       string
	VersionName       string
	MinSDKVersion     string
	RawPackageContent string
}

type manifest struct {
	XMLName     xml.Name `xml:"manifest"`
	VersionCode string   `xml:"versionCode,attr"`
	VersionName string   `xml:"versionName,attr"`
	PackageName string   `xml:"package,attr"`
	Application application
	UsesSdk     usesSdk
}

type application struct {
	XMLName     xml.Name `xml:"application"`
	PackageName string   `xml:"name,attr"`
	AppName     string   `xml:"label,attr"`
}

type usesSdk struct {
	XMLName       xml.Name `xml:"uses-sdk"`
	MinSDKVersion string   `xml:"minSdkVersion,attr"`
}

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

	fmt.Printf("Parsing apk file: %s\n", path)

	var manifestContent bytes.Buffer
	enc := xml.NewEncoder(&manifestContent)
	enc.Indent("", "\t")

	// Parse the apk and validate it
	zipErr, resErr, manErr := apkparser.ParseApk(path, enc)
	if zipErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open the APK: %s", zipErr.Error())
		os.Exit(1)
		return
	}

	if resErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse resources: %s", resErr.Error())
		os.Exit(1)
		return
	}

	if manErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse AndroidManifest.xml: %s", manErr.Error())
		os.Exit(1)
		return
	}

	var manifest manifest
	if err := xml.Unmarshal(manifestContent.Bytes(), &manifest); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal AndroidManifest.xml, error: %s", err)
		os.Exit(1)
		return
	}

	var apk_info = ApkInfo{
		AppName:           manifest.Application.AppName,
		PackageName:       manifest.PackageName,
		VersionCode:       manifest.VersionCode,
		VersionName:       manifest.VersionName,
		MinSDKVersion:     manifest.UsesSdk.MinSDKVersion,
		RawPackageContent: string(manifestContent.Bytes()),
	}

	fmt.Fprintf(os.Stdout, "App Name: %s\n", apk_info.AppName)
	fmt.Fprintf(os.Stdout, "Package Name: %s\n", apk_info.PackageName)
	fmt.Fprintf(os.Stdout, "Version Code: %s\n", apk_info.VersionCode)
	fmt.Fprintf(os.Stdout, "App Name: %s\n", apk_info.MinSDKVersion)
	//fmt.Println(apk_info.RawPackageContent)
}
