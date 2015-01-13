package main

import (
	"flag"
	"github.com/jtunison/stopgo"
)

var serverMode bool
var publish bool

func init() {
	flag.BoolVar(&serverMode, "serverMode", false, "turn on server mode")
	flag.BoolVar(&publish, "publish", false, "publish to s3")
}

func main() {

	flag.Parse()

	if serverMode {
		stopgo.Server()
	} else {

		gopath := os.Getenv("GOPATH")
		overlayPath := fmt.Fprintf("%s/src/github.com/jtunison/stopgo/sample/overlay", gopath)
		outputDir := fmt.Fprintf("%s/src/github.com/jtunison/stopgo/sample/public", gopath)
		stopgo.Build(overlayPath, outputDir)

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
