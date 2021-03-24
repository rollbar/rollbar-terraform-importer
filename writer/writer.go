package writer

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/rollbar/rollbar-terraform-importer/fetcher"
)

// WriteProviderBlocks writes, to a user-defined file, the boilerplate
// necessary for Terraform to pull the Rollbar provider and make the output
// a functioning Terraform project.
func WriteProviderBlocks(filename string) {
	outputFile := writeFile(filename)

	resource := `terraform {
  required_providers {
    rollbar = {
      source = "rollbar/rollbar"
      version = "1.0.6"
    }
  }
}

provider "rollbar" {
}` + "\n\n"
	outputFile.WriteString(resource)
	outputFile.Sync()
	outputFile.Close()
}

// WriteProjectAccessTokens writes, to a user-defined file, the project access
// tokens as Terraform resources.
//
// The name of the resource is just the name of the team, made to conform to
// the limitations of a Terraform resource identifier via sanitizeIdentifier().
func WriteProjectAccessTokens(projects []fetcher.Project, filename string) {
	outputFile := writeFile(filename)

	for _, project := range projects {
		for _, accessToken := range project.AccessTokens {
			resource := `resource "rollbar_project_access_token" "` +
				sanitizeIdentifier(project.Name) + "_" +
				sanitizeIdentifier(accessToken.Name) + `" {` + "\n" +
				`  name = "` + accessToken.Name + `"` + "\n" +
				`  project_id = rollbar_project.` + sanitizeIdentifier(project.Name) + ".id" + "\n" +
				`  depends_on = [rollbar_project.` + sanitizeIdentifier(project.Name) + `]` + "\n" +
				`  rate_limit_window_size = ` + strconv.Itoa(accessToken.RateLimitWindowSize) + "\n" +
				`  rate_limit_window_count = ` + strconv.Itoa(accessToken.RateLimitWindowCount) + "\n" +
				`}` + "\n\n"
			outputFile.WriteString(resource)
		}
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteProjects writes Rollbar projects as Terraform resources to the
// user-defined file.
//
// The name of the resource is just the name of the team, made to conform to
// the limitations of a Terraform resource identifier via sanitizeIdentifier().
func WriteProjects(projects []fetcher.Project, teams []fetcher.Team, filename string) {
	outputFile := writeFile(filename)

	for _, project := range projects {
		projectTeamIDs := "["
		for _, team := range teams {
			for _, teamProject := range team.Projects {
				if teamProject == project.ID {
					projectTeamIDs += "rollbar_team." +
						sanitizeIdentifier(team.Name) + ".id, "
				}
			}
		}
		projectTeamIDs += "]"

		// projectTeams = strings.ReplaceAll(projectTeams, ", ]", "]")

		var resource string
		if projectTeamIDs != "[]" {
			resource = `resource "rollbar_project" "` + sanitizeIdentifier(project.Name) + `" {` + "\n" +
				`  name = "` + project.Name + `"` + "\n" +
				`  team_ids = ` + projectTeamIDs + "\n" +
				`  depends_on = ` + strings.ReplaceAll(projectTeamIDs, ".id", "") + "\n" +
				`}` + "\n\n"
		} else {
			resource = `resource "rollbar_project" "` + sanitizeIdentifier(project.Name) + `" {` + "\n" +
				`  name = "` + project.Name + `"` + "\n" +
				`}` + "\n\n"
		}
		outputFile.WriteString(resource)
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteTeams writes Rollbar teams as Terraform resources to the user-defined
// file.
//
// The name of the resource is just the name of the team, made to conform to
// the limitations of a Terraform resource identifier via sanitizeIdentifier().
func WriteTeams(teams []fetcher.Team, filename string) {
	outputFile := writeFile(filename)

	for _, team := range teams {
		resource := `resource "rollbar_team" "` + sanitizeIdentifier(team.Name) + `" {` + "\n" +
			`  name = "` + team.Name + `"` + "\n" +
			`}` + "\n\n"
		outputFile.WriteString(resource)
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteUsers writes all the Rollbar users as Terraform resources to the
// user-defined file.
//
// The name of the resource is just the name of the team, made to conform to
// the limitations of a Terraform resource identifier via sanitizeIdentifier().
func WriteUsers(users []fetcher.User, filename string) {
	outputFile := writeFile(filename)

	for _, user := range users {
		teams := "["
		for _, team := range user.Teams {
			teams += "rollbar_team." + sanitizeIdentifier(team.Name) + ".id, "
		}
		teams += "]"

		resource := `resource "rollbar_user" "` + sanitizeIdentifier(user.Username) + `" {` + "\n"
		if user.Email != "" {
			resource = resource + `  email = "` + user.Email + `"` + "\n"
		}
		if len(user.Teams) > 0 {
			resource = resource + `  team_ids = "` + teams + `"` + "\n"
		}
		resource = resource + `}` + "\n\n"
		_, err := outputFile.WriteString(resource)
		if err != nil {
			log.Fatal("Failed to write to file.", err)
		}
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteProjectAccessTokenImportCommands extracts the project name and the
// access token value for each access token to generate a Terraform import
// for every access token in a given project. The resource names for the
// projects and access tokens the same way they are for the resources
// themselves, via sanitizeIdentifier().
func WriteProjectAccessTokenImportCommands(projects []fetcher.Project, filename string) {
	outputFile := writeFile(filename)
	for _, project := range projects {
		for _, accessToken := range project.AccessTokens {
			outputFile.WriteString("terraform import rollbar_project_access_token." +
				sanitizeIdentifier(project.Name) + "_" +
				sanitizeIdentifier(accessToken.Name) + " " + strconv.Itoa(project.ID) +
				"/" + accessToken.AccessToken + "\n")
		}
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteProjectImportCommands iterates through an array of Project structs and
// extracts the project names and IDs from them to generate the Terraform
// import command for each project. The resource names for the
// projects are generated via sanitizeIdentifier() like the are done for the
// resources themselves.
func WriteProjectImportCommands(projects []fetcher.Project, filename string) {
	outputFile := writeFile(filename)
	for _, project := range projects {
		outputFile.WriteString("terraform import rollbar_project." +
			sanitizeIdentifier(project.Name) + " " + strconv.Itoa(project.ID) + "\n")
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteTeamImportCommands iterates through an array of Team structs and
// extracts the team names and IDs from them to generate the Terraform import
// command for each project. The resource names for the teams are generated via
// sanitizeIdentifier() like the are done for the resources themselves.
func WriteTeamImportCommands(teams []fetcher.Team, filename string) {
	outputFile := writeFile(filename)
	for _, team := range teams {
		outputFile.WriteString("terraform import rollbar_team." +
			sanitizeIdentifier(team.Name) + " " + strconv.Itoa(team.ID) + "\n")
	}
	outputFile.Sync()
	outputFile.Close()
}

// WriteUserImportCommands iterates through an array of User structs and
// extracts the usernames and IDs from them to generate the Terraform import
// command for each user. The resource names for the users are generated via
// sanitizeIdentifier() like the are done for the resources themselves.
func WriteUserImportCommands(users []fetcher.User, filename string) {
	outputFile := writeFile(filename)
	for _, user := range users {
		outputFile.WriteString("terraform import rollbar_user." +
			sanitizeIdentifier(user.Username) + " " + strconv.Itoa(user.ID) + "\n")
	}
	outputFile.Sync()
	outputFile.Close()
}

// writeFile accepts a filename and returns a file descriptor. This is intended
// for writing the Terraform files and import commands to disk and to avoid
// having to write this for every single function above.
func writeFile(filename string) (outputFile *os.File) {
	outputFile, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return outputFile
}

// sanitizeIdentifier takes a given string and tries to make it conform to the
// limitations of valid Terraform identifiers.
func sanitizeIdentifier(token string) (identfier string) {
	// First, replace all the spaces and dots with underscores.
	forbiddenUnderscore := regexp.MustCompile("[\\./, ]")
	token = forbiddenUnderscore.ReplaceAllString(token, "_")

	// Then, remove known invalid characters.
	forbiddenStrip := regexp.MustCompile(`[\\(\\)\\?]`)
	token = forbiddenStrip.ReplaceAllString(token, "")

	return token
}
