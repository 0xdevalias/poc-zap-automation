package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/zaproxy/zap-api-go/zap"
)

// Ref: https://github.com/zaproxy/zap-api-go/blob/master/example/example.go

type ZapClient struct {
	zap.Interface
}

var target string

func init() {
	flag.StringVar(&target, "target", "http://bodgeit:8080/bodgeit/", "target address")
	//flag.StringVar(&target, "target", "http://juiceshop:3000", "target address")
	flag.Parse()
}

func main() {
	cfg := &zap.Config{
		Proxy: "http://127.0.0.1:8080",
	}
	client, err := zap.NewClient(cfg)
	if err != nil {
		log.Fatalf("init error: %s", err)
	}

	zap := ZapClient{client}

	// Configure our spider rules
	client.Spider().SetOptionMaxDepth(5)
	//client.Spider().SetOptionMaxParseSizeBytes(1024*1024*5)
	//client.Spider().SetOptionMaxDuration(5) // In seconds, minutes, ???
	//client.Spider().SetOptionRequestWaitTime() // In seconds, minutes, ???

	client.Spider().SetOptionPostForm(false)
	client.Spider().SetOptionProcessForm(false) // ???

	client.Spider().SetOptionParseComments(true)
	client.Spider().SetOptionParseGit(true)
	client.Spider().SetOptionParseRobotsTxt(true)
	client.Spider().SetOptionParseSVNEntries(true)
	client.Spider().SetOptionParseSitemapXml(true)

	// Start spidering the target
	// Ref: https://github.com/zaproxy/zaproxy/wiki/ApiGen_spider
	log.Printf("Spider: %s", target)
	resp, err := client.Spider().Scan(target, "", "", "", "")
	if err != nil {
		log.Fatalf("spider error: %s", err)
	}

	// The scan now returns a scan id to support concurrent scanning
	scanid := resp["scan"].(string)
	for {
		time.Sleep(1000 * time.Millisecond)
		resp, _ = client.Spider().Status(scanid)
		progress, _ := strconv.Atoi(resp["status"].(string))
		if progress >= 100 {
			break
		}
	}
	log.Println("Spider complete")

	// Give the passive scanner a chance to complete
	time.Sleep(2000 * time.Millisecond)

	//fmt.Println("Active scan : " + target)
	//resp, err = client.Ascan().Scan(target, "True", "False", "", "", "", "")
	//if err != nil {
	//	log.Fatalf("active scan error: %s", err)
	//}
	//
	//// The scan now returns a scan id to support concurrent scanning
	//scanid = resp["scan"].(string)
	//for {
	//	time.Sleep(5000 * time.Millisecond)
	//	resp, _ = client.Ascan().Status(scanid)
	//	progress, _ := strconv.Atoi(resp["status"].(string))
	//	fmt.Printf("Active Scan progress : %d\n", progress)
	//	if progress >= 100 {
	//		break
	//	}
	//}
	//fmt.Println("Active Scan complete")

	//fmt.Println("Alerts:")
	//report, err := client.Core().Xmlreport()
	//if err != nil {
	//	log.Fatalf("xml report error: %s", err)
	//}
	//fmt.Println(string(report))

	//allUrls, err := client.Core().Urls("")
	allUrls, err := zap.Urls("")
	fmt.Printf("%s", strings.Join(allUrls, "\n"))
}

func (zap ZapClient) Urls(baseUrl string) ([]string, error) {
	allUrls, err := zap.Core().Urls(baseUrl)
	if err != nil {
		return nil, err
	}

	var urlSlice []string
	for _, url := range allUrls["urls"].([]interface{}) {
		urlSlice = append(urlSlice, url.(string))
	}

	return urlSlice, nil
}
