module github.com/rollbar/rollbar-terraform-importer

go 1.16

require (
	github.com/fatih/color v1.10.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/rollbar/rollbar-terraform-importer/fetcher v0.0.0
	github.com/rollbar/rollbar-terraform-importer/writer v0.0.0
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/rollbar/rollbar-terraform-importer/fetcher => ./fetcher

replace github.com/rollbar/rollbar-terraform-importer/writer => ./writer
