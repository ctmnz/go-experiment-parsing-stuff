package main

import (
	"fmt"
	"log"

	// The official Kubernetes YAML library
	yamlutil "sigs.k8s.io/yaml"
)

// Define a struct that mirrors the structure you expect in a YAML file.
// Notice the use of `json:"..."` tags. sigs.k8s.io/yaml respects these tags
// (which is crucial for Kubernetes compatibility).
type Config struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Replicas int    `json:"replicas"`
		Image    string `json:"image"`
	} `json:"spec"`
}

func main() {
	// --- PART 1: UNMARSHAL (YAML to Go Struct) ---

	// A sample YAML string representing a simple Kubernetes Deployment-like object
	yamlInputOld := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: default
spec:
  replicas: 3
  image: nginx:latest
`
  yamlInput := `
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



	fmt.Println("--- 1. Original YAML Input ---")
	fmt.Println(yamlInputOld)
	fmt.Println(yamlInput)

	// 1. Create an empty struct to hold the unmarshaled data
	var config Config

	// 2. Unmarshal the YAML byte slice into the Config struct
	if err := yamlutil.Unmarshal([]byte(yamlInput), &config); err != nil {
		log.Fatalf("error unmarshaling YAML: %v", err)
	}

	fmt.Println("\n--- 2. Go Struct Data (after Unmarshal) ---")
	fmt.Printf("API Version: %s\n", config.APIVersion)
	fmt.Printf("Kind: %s\n", config.Kind)
	fmt.Printf("Metadata Name: %s\n", config.Metadata.Name)
	fmt.Printf("Replicas: %d\n", config.Spec.Replicas)
	fmt.Printf("Image: %s\n", config.Spec.Image)

	// --- PART 2: MARSHAL (Go Struct to YAML) ---

	// 1. Modify the struct data (e.g., scale up)
	config.Spec.Replicas = 5
	config.Metadata.Name = "my-app-scaled"

	// 2. Marshal the Go struct back into a YAML byte slice
	yamlOutput, err := yamlutil.Marshal(config)
	if err != nil {
		log.Fatalf("error marshaling Go struct to YAML: %v", err)
	}

	fmt.Println("\n--- 3. YAML Output (after Marshal and modification) ---")
	// Print the resulting YAML (converted to a string)
	fmt.Println(string(yamlOutput))
}
