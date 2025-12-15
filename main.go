package main

import (
	"fmt"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	log "gopkg.in/op/go-logging.v1"
)

func main() {
	fmt.Println("Script started")

	log.SetLevel(log.WARNING, "yq-lib")

	yamlMultiContent := `
# Document 0: Service
---
apiVersion: v1
kind: Service
metadata:
  name: my-service
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
`

	myCustPrefs := yqlib.YamlPreferences{
		Indent:                      2,
		ColorsEnabled:               false,
		LeadingContentPreProcessing: true,
		PrintDocSeparators:          false,
		UnwrapScalar:                true,
		EvaluateTogether:            false,
		FixMergeAnchorToSpec:        false,
	}

	//  myPrefs := yqlib.NewDefaultYamlPreferences()
	expression := "select(.kind == \"Deployment\") .spec.template.spec.containers.[].image"
	evaluator := yqlib.NewStringEvaluator()
	myEncoder := yqlib.NewYamlEncoder(myCustPrefs)
	myDecoder := yqlib.NewYamlDecoder(myCustPrefs)
	result, err := evaluator.Evaluate(expression, yamlMultiContent, myEncoder, myDecoder)
	if err != nil {
		panic(err)
	}
	if len(result) > 0 {
		fmt.Println("image:", result)
	}
	fmt.Println("Script ended")
}
