# Rollbar Terraform Importer
This is an importer designed to help migrate Rollbar account data to be managed
by Terraform.

## How It Works
This utility queries the public Rollbar API to obtain user, team and project
information and produces:

- a file with the Terraform import commands,
- a file (or set of files) containing the Rollbar account information as
Terraform resources.

### Flags
- *-accessToken*: Pass a Rollbar account access token with rights to read.
- *-singleFile*: By default, the Terraform files are produced with a file per
type (*e.g.* user.tf, projects.tf, access_tokens.tf), but this can be disabled
to write them all to a single file.
- *-outPath*: The directory to write the generated files to.

### Examples
- `rollbar-terraform-importer -accessToken 53lkj34802lkj2342341l` will
generate an `import` file contain all import commands, as well as
`access_tokens.tf`, `projects.tf`, `teams.tf` and `users.tf` to the current
working directory.
- `rollbar-terraform-importer -accessToken 53lkj34802lkj2342341l -singleFile`
will do the same thing, except it will write all Terraform resources into a
single file called `rollbar_account.tf`. Terraform import files are still
generated into a file called `import`.

## Caveats
The importer requires some manual review to ensure that all resources and names
are correct. For instance, access tokens are not guaranteed to have unique
names and the importer leaves the decision on naming them to the user.

