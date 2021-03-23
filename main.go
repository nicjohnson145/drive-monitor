package main

import (
	"bufio"
	"os"
	"regexp"
)

type Disk struct {
	ID               string
	Status           string
	State            string
	FailurePredicted bool
}

func main() {
	// First element is garbage
	sections := splitToSections(os.Stdin)[1:]
	disks := sectionsToStruct(sections)
}

func splitToSections(input *os.File) [][]string {
	scanner := bufio.NewScanner(input)

	outer := [][]string{}
	var currentList []string

	for scanner.Scan() {
		if currentList == nil {
			currentList = []string{}
		}
		line := scanner.Text()
		if line != "" {
			currentList = append(currentList, line)
		} else {
			outer = append(outer, currentList)
			currentList = nil
		}
	}

	return outer
}

func sectionsToStruct(sections [][]string) []Disk {
	disks := []Disk{}

	id := regexp.MustCompile(`^ID:\s+: (\d:\d:\d)$`)
	status := regexp.MustCompile(`^Status\s+: (.*)$`)
	state := regexp.MustCompile(`^State\s+: (.*)$`)

	return disks
}

