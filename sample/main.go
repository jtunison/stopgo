package main

import (
	"flag"
	"github.com/jtunison/stopgo"
	"os"
)

var serverMode bool
var publish bool
var outputDir string

func init() {
	flag.BoolVar(&serverMode, "serverMode", false, "turn on server mode")
	flag.BoolVar(&publish, "publish", false, "publish to s3")
	flag.StringVar(&outputDir, "outputDir", "public", "name of output directory")
}

func main() {

	flag.Parse()

	if serverMode {
		stopgo.Server()
	} else {

		stopgo.Build()

		// publish!
		myS3bucket := "jameslee.com"
		awsAccessKey := "ASDFASDFASDF"
		awsSecretKey := "SDFSDFSDF"
		if publish {
			err := stopgo.Publish("public", myS3bucket, awsAccessKey, awsSecretKey)
			if err != nil {
				panic(err)
			}
		}
	}

}
