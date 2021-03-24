package main

import (
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
	"github.com/rollbar/rollbar-terraform-importer/fetcher"
	"github.com/rollbar/rollbar-terraform-importer/writer"
)

func main() {
	// Just handle the flag parsing and hand it off to generate() to do the
	// heavy lifting.

	var accessToken = flag.String("accessToken", "NO_TOKEN", "Rollbar account access token.")
	var singleFile = flag.Bool("singleFile", false, "Write to a single Terraform file.")
	var outPath = flag.String("out", ".", "Output directory for generated files.")
	flag.Parse()

	errorColor := color.New(color.FgRed).Add(color.Bold)

	// Ensure that an access token was provided as an argument.
	if *accessToken == "NO_TOKEN" {
		errorColor.Fprintln(os.Stderr, "[ERROR] A Rollbar access token must be provided.")
		flag.Usage()
		os.Exit(-1)
	}

	// Validate the access token is actually an access token.
	validate := validator.New()
	vErrs := validate.Var(*accessToken, "required,alphanumunicode")
	if vErrs != nil {
		errorColor.Fprintln(os.Stderr, "[ERROR] Provided access token is not a valid access token.")
		os.Exit(-1)
	}

	// Validate that if any path, except the default was given, that it actually exists.
	if _, err := os.Stat(*outPath); os.IsNotExist(err) {
		errorColor.Fprintln(os.Stderr, "[ERROR] Invalid file path provided for output.")
		os.Exit(-2)
	}

	// Do something based on the user-defined options.
	generate(*singleFile, *accessToken, *outPath)
}

// generate takes the values of the user-defined flags and uses them to define
// how to generate the Terraform files.
func generate(singleFile bool, accessToken string, outPath string) {
	// Make output colorful for visibility.
	stdColor := color.New(color.FgWhite).Add(color.Bold)
	successColor := color.New(color.FgGreen).Add(color.Bold)

	// Fetch the necessary data via the Rollbar API.
	projects := fetcher.FetchProjects(accessToken)
	teams := fetcher.FetchTeams(accessToken)
	users := fetcher.FetchUsers(accessToken)

	if singleFile {
		/*
		 * If the user pases the single file flag, just append constantly to the same file.
		 *
		 * FIXME: This is clunky.
		 */
		writer.WriteProviderBlocks(outPath + "/rollbar_account.tf")
		writer.WriteTeams(teams, outPath+"/rollbar_account.tf")
		writer.WriteProjects(projects, teams, outPath+"/rollbar_account.tf")
		writer.WriteProjectAccessTokens(projects, outPath+"/rollbar_account.tf")
		writer.WriteUsers(users, outPath+"/rollbar_account.tf")
		successColor.Fprintln(os.Stdout, "Rendered All Account Resources to rollbar_account.tf.")

		writer.WriteProjectAccessTokenImportCommands(projects, outPath+"/import")
		writer.WriteProjectImportCommands(projects, outPath+"/import")
		writer.WriteTeamImportCommands(teams, outPath+"/import")
		writer.WriteUserImportCommands(users, outPath+"/import")
		stdColor.Fprintln(os.Stdout, "Rendered Terraform Import Commands to import")
	} else {
		writer.WriteProviderBlocks(outPath + "/main.tf")

		writer.WriteTeams(teams, outPath+"/teams.tf")
		successColor.Fprintln(os.Stdout, "Rendered Team Resources to teams.tf.")

		writer.WriteProjects(projects, teams, outPath+"/projects.tf")
		successColor.Fprintln(os.Stdout, "Rendered Project Resources to projects.tf.")

		writer.WriteProjectAccessTokens(projects, outPath+"/access_tokens.tf")
		successColor.Fprintln(os.Stdout, "Rendered Access Token Resources to access_tokens.tf.")

		writer.WriteUsers(users, outPath+"/users.tf")
		successColor.Fprintln(os.Stdout, "Rendered User Resources to users.tf.")

		writer.WriteProjectAccessTokenImportCommands(projects, outPath+"/import")
		writer.WriteProjectImportCommands(projects, outPath+"/import")
		writer.WriteTeamImportCommands(teams, outPath+"/import")
		writer.WriteUserImportCommands(users, outPath+"/import")
		stdColor.Fprintln(os.Stdout, "Rendered Terraform Import Commands to import")
	}
}
