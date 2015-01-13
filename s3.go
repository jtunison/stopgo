package stopgo

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

func RetrieveS3Map(b *s3.Bucket) map[string]string {
	var m = make(map[string]string)

	res, err := b.List("", "", "", 1000)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range res.Contents {
		m[v.Key] = strings.Replace(v.ETag, `"`, "", -1)
	}

	return m
}

func Publish(localpath, bucketName, awsAccessKey, awsSecretKey string) error {

	auth := aws.Auth{
		AccessKey: awsAccessKey,
		SecretKey: awsSecretKey,
	}

	connection := s3.New(auth, aws.USEast)
	b := connection.Bucket(bucketName)
	s3map := RetrieveS3Map(b)

	// step 1:  copy all files that aren't already there
	err := filepath.Walk("public", func(path string, fi os.FileInfo, _ error) error {

		name := filepath.Base(path)
		hidden := strings.HasPrefix(name, ".")
		if !fi.IsDir() && !hidden {

			// has it changed?

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			// have to omit the leading 'public/'
			ext := filepath.Ext(name)
			contentType := mime.TypeByExtension(ext)
			key := path[7:]
			md5bytes := md5.Sum(content)
			md5string := hex.EncodeToString(md5bytes[:])

			if s3map[key] == "" {
				log.Printf("Adding '%s'\n", key)
				err = b.Put(key, content, contentType, s3.PublicRead)
				if err != nil {
					return err
				}

			} else if s3map[key] != md5string {
				log.Printf("Updating '%s'\n", key)
				// b.Del(key)
				// if err != nil {
				// 	return err
				// }
				err = b.Put(key, content, contentType, s3.PublicRead)
				if err != nil {
					return err
				}

			}
			// else {
			// 	log.Printf("Skipping '%s' (unchanged)\n", key)
			// }

			delete(s3map, key)

		}
		return nil
	})
	if err != nil {
		return err
	}

	// step 2: remove any s3 files that are no longer relevant
	for key, _ := range s3map {
		log.Printf("Deleting '%s'\n", key)
		err = b.Del(key)
		if err != nil {
			return err
		}
	}

	return nil

}
