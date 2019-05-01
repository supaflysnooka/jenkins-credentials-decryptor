package main

import (
	"github.com/beevik/etree"
	"io/ioutil"
	"regexp"
)

const (
	credentialsXpath = "//java.util.concurrent.CopyOnWriteArrayList/*"
)

type Credential struct {
	tags map[string]string
}

func ReadCredentials(path string) *[]Credential {
	credentials := make([]Credential, 0)
	for _, credentialNode := range readCredentialsXml(path).FindElements(credentialsXpath) {
		credential := &Credential{
			tags: map[string]string{},
		}
		for _, field := range credentialNode.ChildElements() {
			credential.tags[field.Tag] = field.Text()
		}
		credentials = append(credentials, *credential)
	}
	return &credentials
}

/*
 HACK ALERT:
 Stripping xml version line as current native and third party xml decoders
 refuses to parse xml version 1.0+
 Jenkins uses xml version 1.1+ so this may blow up.
*/
func readCredentialsXml(path string) *etree.Document {
	credentials, err := ioutil.ReadFile(path)
	check(err)
	sanitizedCredentials := regexp.
		MustCompile("(?m)^.*<?xml.*$").
		ReplaceAllString(string(credentials), "")
	document := etree.NewDocument()
	err = document.ReadFromString(sanitizedCredentials)
	check(err)
	return document
}
