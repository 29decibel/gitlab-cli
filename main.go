package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Config is for the information you need to config
type Config struct {
	DefaultProject string
	Host           string
	PrivateToken   string
	ProjectsMap    map[string]string
}

// MergeRequest is the struct for merge reuqest
type MergeRequest struct {
	ID           int    `json: id`
	Title        string `json: title`
	Assignee     Person `json: assignee`
	State        string `json: state`
	SourceBranch string `json: source_branch`
	TargetBranch string `json: target_branch`
	CreatedAt    string `json: created_at`
}

func (mergeRequest MergeRequest) String() string {
	var state string
	switch mergeRequest.State {
	case "closed":
		state = fmt.Sprintf("\033[38;5;160m%s\033[39m", mergeRequest.State)
	case "merged":
		state = fmt.Sprintf("\033[38;5;22m%s\033[39m", mergeRequest.State)
	default:
		state = mergeRequest.State
	}
	return fmt.Sprintf("[%s] %s(%s)", state, mergeRequest.Title, mergeRequest.Assignee.Name)
}

// Person is the general user information
type Person struct {
	ID       int    `json: id`
	UserName string `json: username`
	Email    string `json: email`
	Name     string `json: name`
}

var config *Config
var perPage = flag.Int("per_page", 10, "default counts per page")

func parseConfig() *Config {
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Please make sure your have a config.json exist!")
	}
	json.Unmarshal(configFile, &config)

	return config
}

func main() {
	// parse arguments
	flag.Parse()

	// parse config
	parseConfig()

	for _, mr := range GetMergeRequests(config.ProjectsMap[config.DefaultProject], 1) {
		fmt.Println(mr)
	}
}

// GetMergeRequests is to get all merge request for a project
func GetMergeRequests(projectID string, page int) []MergeRequest {
	url := fmt.Sprintf("%s/projects/%s/merge_requests?private_token=%s&page=%d&per_page=%d", config.Host, projectID, config.PrivateToken, page, *perPage)
	response, err := http.Get(url)
	if err != nil {
		panic("something wrong with the requests")
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Can not parse the response body")
	}

	var mergeRequets []MergeRequest
	json.Unmarshal(contents, &mergeRequets)
	return mergeRequets
}

// GetProjects is to get all projects
func GetProjects() string {
	url := fmt.Sprintf("%s/projects?private_token=%s", config.Host, config.PrivateToken)
	response, err := http.Get(url)
	if err != nil {
		panic("something wrong with the requests")
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Can not parse the response body")
	}
	return string(contents)

}
