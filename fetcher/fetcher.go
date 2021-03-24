package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// FetchProjects retrieves the list of projects in a Rollbar account and
// returns it as a []Project.
//
// Once the list of projects is fetched, the function iterates through the list
// of projects and appends any access tokens associated with each project.
func FetchProjects(accessToken string) (projects []Project) {
	rawProjects := fetch(accessToken, "projects")

	var data projectResponse
	err := json.Unmarshal([]byte(rawProjects), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	projects = data.Result

	for i, _ := range projects {
		fetchProjectAccessTokens(accessToken, &projects[i])
	}

	return projects
}

// FetchTeams retrieves the list of teams in a Rollbar account and returns it
// as a []Team.
//
// Once the initial team list is fetched, this function loops through all of
// the teams and appends any identified projects and users associated with each
// of the teams.
func FetchTeams(accessToken string) (teams []Team) {
	rawTeams := fetch(accessToken, "teams")

	var data teamResponse
	err := json.Unmarshal([]byte(rawTeams), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	teams = data.Result

	for i := range teams {
		fetchTeamProjects(accessToken, &teams[i])
		fetchTeamUsers(accessToken, &teams[i])
	}

	return teams
}

// FetchUsers retrieves the list of users in a Rollbar account and returns it
// as a []User.
//
// Once the initial user metadata is fetched, this function loops through the
// list and captures the projects and teams that each user is associated with
// as []ints.
func FetchUsers(accessToken string) (users []User) {
	rawUsers := fetch(accessToken, "users")

	var data userResponse
	err := json.Unmarshal([]byte(rawUsers), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	users = data.Result.Users
	for i, _ := range users {
		fetchUserTeams(accessToken, &users[i])
	}
	return users
}

// fetchProjectAccessTokens retrieves the access tokens a given project is
// associated with.
//
// It lacks a return value, as it appends the returned access tokens to the
// passed Project struct's AccesTokens property.
func fetchProjectAccessTokens(accessToken string, project *Project) {
	projectID := strconv.Itoa(project.ID)
	endpoint := "project/" + projectID + "/access_tokens"
	rawProjectAccessTokens := fetch(accessToken, endpoint)

	var data accessTokenResponse
	err := json.Unmarshal([]byte(rawProjectAccessTokens), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	accessTokens := data.Result
	for _, accessToken := range accessTokens {
		project.AccessTokens = append(project.AccessTokens, accessToken)
	}
}

// fetchTeamProjects retrieves the projects a given team is associated with.
//
// It lacks a return value as it appends the returned teams to the
// passed Team struct's Projects property.
func fetchTeamProjects(accessToken string, team *Team) {
	teamID := strconv.Itoa(team.ID)
	endpoint := "team/" + teamID + "/projects"
	rawTeamProjects := fetch(accessToken, endpoint)

	var data teamProjectsResponse
	err := json.Unmarshal([]byte(rawTeamProjects), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	for _, project := range data.Result {
		team.Projects = append(team.Projects, project.ProjectID)
	}
}

// fetchTeamUsers retrieves the users a given team is associated with.
//
// It lacks a return value as it appends the returned teams to the
// passed Team struct's Users property.
func fetchTeamUsers(accessToken string, team *Team) {
	teamID := strconv.Itoa(team.ID)
	endpoint := "team/" + teamID + "/users"
	rawTeamUsers := fetch(accessToken, endpoint)

	var data teamUsersResponse
	err := json.Unmarshal([]byte(rawTeamUsers), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	for _, user := range data.Result {
		team.Users = append(team.Users, user.UserID)
	}
}

// fetchUserTeams retrieves the teams a given user is associated with.
//
// It lacks a return value as it appends the returned teams to the
// passed User struct's Teams property.
func fetchUserTeams(accessToken string, user *User) {
	userID := strconv.Itoa(user.ID)
	endpoint := "user/" + userID + "/teams"
	rawUserTeams := fetch(accessToken, endpoint)

	var data userTeamResponse
	err := json.Unmarshal([]byte(rawUserTeams), &data)
	if err != nil {
		log.Fatal("Error parsing JSON response body.", err)
	}
	if data.Err != 0 {
		log.Fatal("API returned an error", data.Result)
	}

	teams := data.Result.Teams

	for _, team := range teams {
		user.Teams = append(user.Teams, team)
	}
}

// fetch will make an HTTP GET request to the requested API endpoint and return
// the response body. If it cannot, the importer will throw a fatal error and
// die, as none of this works unless all the requests are completed
// successfully.
func fetch(accessToken string, endpoint string) (body []byte) {
	baseURL := "https://api.rollbar.com/api/1/"
	apiURL := baseURL + endpoint
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatal("Unable to generate HTTP request.", err)
	}

	req.Header.Set("X-Rollbar-Access-Token", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading HTTP response.", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body.", err)
	}
	return body
}
