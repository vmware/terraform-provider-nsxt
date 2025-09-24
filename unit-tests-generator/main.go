package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

func main() {
	rawData, err := ioutil.ReadFile("generator_input.yaml")
	data := make(map[string]interface{})
	if err != nil {
		log.Fatalf("error reading generator input YAML file: %v", err)
	}
	err = yaml.Unmarshal(rawData, &data)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}
	var inputFileName string
	tempSegs := strings.Split(data["InputFileName"].(string), "/")
	if len(tempSegs) > 1 {
		inputFileName = tempSegs[len(tempSegs)-1]
	} else {
		inputFileName = data["InputFileName"].(string)
	}
	inArr := strings.Split(strings.TrimSuffix(inputFileName, ".go"), "_")
	for i, str := range inArr {
		inArr[i] = capitalizeFirst(str)
	}
	data["ResourceName"] = strings.Join(inArr, "")
	inResource := strings.TrimSuffix(inputFileName, ".go")
	path := "../nsxt/"

	tmpl, err := template.ParseFiles("ut_go_template.tmpl")
	if err != nil {
		log.Fatal("parse error:", err)
	}

	outFile, err := os.Create(path + "ut_" + inResource + "_test.go")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, data)
	if err != nil {
		log.Fatal("execute error:", err)
	}

	log.Println("Generated the unit test boiler plate code and written to ", path+"ut_"+inResource+"_test.go")
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
