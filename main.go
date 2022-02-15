package main

import (
	"bufio"
	"fmt"
	"github.com/levigross/grequests"
	"os"
	"regexp"
)

type Disk struct {
	ID               string
	Status           string
	State            string
	FailurePredicted bool
}

type Config struct {
	UserToken string
	AppToken  string
}

func main() {
	conf := parseConfigOrDie()
	// First element is garbage
	sections := splitToSections(os.Stdin)[1:]
	disks := sectionsToStruct(sections)
	failingDisks := findFailingDrives(disks)
	alertFailingDisks(failingDisks, conf)
}

func parseConfigOrDie() Config {
	c := Config{}

	if val, ok := os.LookupEnv("APP_TOKEN"); !ok {
		fmt.Println("APP_TOKEN not set")
		os.Exit(1)
	} else {
		c.AppToken = val
	}

	if val, ok := os.LookupEnv("USER_TOKEN"); !ok {
		fmt.Println("USER_TOKEN not set")
		os.Exit(1)
	} else {
		c.UserToken = val
	}

	return c
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

	id := regexp.MustCompile(`^ID\s+: (?P<ID>\d+:\d+:\d+)$`)
	status := regexp.MustCompile(`^Status\s+: (?P<Status>.*)$`)
	state := regexp.MustCompile(`^State\s+: (?P<State>.*)$`)
	failurePredicted := regexp.MustCompile(`^Failure Predicted\s+: (?P<FailurePredicted>.*)$`)

	for _, lines := range sections {
		d := Disk{}
		for _, line := range lines {
			if match := id.FindStringSubmatch(line); match != nil {
				d.ID = match[1]
				continue
			}
			if match := status.FindStringSubmatch(line); match != nil {
				d.Status = match[1]
				continue
			}
			if match := state.FindStringSubmatch(line); match != nil {
				d.State = match[1]
				continue
			}
			if match := failurePredicted.FindStringSubmatch(line); match != nil {
				d.FailurePredicted = match[1] != "No"
				continue
			}
		}
		disks = append(disks, d)
	}

	return disks
}

func findFailingDrives(disks []Disk) []Disk {
	failing := []Disk{}
	for i, d := range disks {
		if d.Status != "Ok" || d.FailurePredicted == true {
			failing = append(failing, disks[i])
		}
	}

	return failing
}

func alertFailingDisks(disks []Disk, config Config) {
	if len(disks) == 0 {
		return
	}

	resp, err := grequests.Post(
		"https://api.pushover.net/1/messages.json",
		&grequests.RequestOptions{
			Data: map[string]string{
				"token": config.AppToken,
				"user": config.UserToken,
				"message": fmt.Sprintf("%v disks reported unhealthy", len(disks)),
			},
		},
	)
	if err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Printf("Bad response from pushover API\n%+v\n", resp)
		os.Exit(1)
	}
}

