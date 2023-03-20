package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func findDependencies(fileContent string) []string {
	dependencies := []string{}
	dependencyBlockRegex := regexp.MustCompile(`dependencies\s*{([\s\S]*?)}`)
	dependencyLineRegex := 
regexp.MustCompile(`(?:implementation|api|compile)\s*\(?["']([^:]+):([^:]+):([^'"\s$]+)["']`)

	matches := dependencyBlockRegex.FindAllStringSubmatch(fileContent, -1)

	for _, match := range matches {
		deps := dependencyLineRegex.FindAllStringSubmatch(match[1], -1)
		for _, dep := range deps {
			dependencies = append(dependencies, fmt.Sprintf("%s:%s:%s", 
dep[1], dep[2], dep[3]))
		}
	}

	return dependencies
}

func processFile(filePath string, depsMap map[string]bool) {
	contentBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", filePath)
		return
	}

	content := string(contentBytes)
	deps := findDependencies(content)

	for _, dep := range deps {
		depsMap[dep] = true
	}
}

func main() {
	rootDir := "node_modules"
	depsMap := make(map[string]bool)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) 
error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), "build.gradle") {
			processFile(path, depsMap)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		return
	}

	dependencies := make([]string, 0, len(depsMap))
	for dep := range depsMap {
		dependencies = append(dependencies, dep)
	}

	sort.Strings(dependencies)

	for _, dep := range dependencies {
		fmt.Println(dep)
	}
}

