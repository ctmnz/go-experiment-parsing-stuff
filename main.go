package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	log "gopkg.in/op/go-logging.v1"
)

const semVerPattern = `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

var semVerRegex = regexp.MustCompile(semVerPattern)

func isValidSemVer(tag string) bool {
	if !semVerRegex.MatchString(tag) {
		return false
	}
	return true
}

func getTag(imgName string) (string, error) {
	nameSplit := strings.SplitN(imgName, ":", 2)
	fmt.Println(len(nameSplit))
	if len(nameSplit) < 2 {
		return "", errors.New("you don't have tags in the image name")
	}
	return nameSplit[1], nil
}

func analyzeFile(fullFilePath string) {
	log.SetLevel(log.WARNING, "yq-lib")

	myCustPrefs := yqlib.YamlPreferences{
		Indent:                      2,
		ColorsEnabled:               false,
		LeadingContentPreProcessing: true,
		PrintDocSeparators:          false,
		UnwrapScalar:                true,
		EvaluateTogether:            false,
		FixMergeAnchorToSpec:        false,
	}

	data, err := os.ReadFile(fullFilePath)
	if err != nil {
		panic(err)
	}
	// fmt.Print(string(data))
	expression := "select(.kind == \"Deployment\") .spec.template.spec.containers.[].image"
	evaluator := yqlib.NewStringEvaluator()
	myEncoder := yqlib.NewYamlEncoder(myCustPrefs)
	myDecoder := yqlib.NewYamlDecoder(myCustPrefs)
	result, err := evaluator.Evaluate(expression, string(data), myEncoder, myDecoder)
	if err != nil {
		panic(err)
	}
	if len(result) > 0 {
		if len(result) > 1 {
			scanner := bufio.NewScanner(strings.NewReader(result))
			for scanner.Scan() {
				fmt.Println("image:", scanner.Text())
				t, err := getTag(scanner.Text())
				if err != nil {
					panic(err)
				}
				fmt.Printf("Tag: %v\n", t)
				fmt.Printf("Is the Tag semver? : %v\n", isValidSemVer(t))
			}
		} else {
			fmt.Println("image:", result)
		}
	}
}

func getYamlFiles(dirPath string) ([]string, error) {
	yamlFiles, err := filepath.Glob(filepath.Join(dirPath, "*.yaml"))
	if err != nil {
		return nil, err
	}
	return yamlFiles, nil
}

func main() {
	fmt.Println("Script started")

	analyzeFile("./manifests/some-deployment.yaml")

	yfs, err := getYamlFiles("./manifests/")
	if err != nil {
		panic(err)
	}
	for _, yf := range yfs {
		fmt.Println(yf)
		analyzeFile(yf)
	}

	fmt.Println("Script ended")
}
