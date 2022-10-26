# Go Template

This repository acts as a template for golang repositories.

## After creating repo from template

The Template Cleanup workflow will perform some cleanup from the template and publish in the form of a pull request.  Please make sure to review this first pull requet and merge in before attempting to add features as it is trying to set up the repo for proper usage.

The following items should be verified/changed for usage as a child repo:

### The Github settings file

1. Set the repository name in the `.github/settings.yml` (Handled by Template Cleanup)
1. Set the repository description in the `.github/settings.yml` (Handled by Template Cleanup)
1. Add the owner team as admins in `.github/settings.yml` under the teams section.

### The BD configuration file

1. Set the name of the application (A default assumption handled by Template Cleanup)
  - This should be the friendly name to reference the application and should match the key for the SSM param store path.  `/{env}/{name}/{var}`
1. Set the application type (A default assumption handled by Template Config)
  - This type is to reflect the type of repo the repository represents such as: `application`, `library`, `support`
1. Set the team.  The name of this team is not the same as the github team, but a friendly name for the team such as: `web`, `data`, `datascience`, `mobile`, `techops`, `it`
1. Set any additional tags that are desired for filtering the repository

### The CODEOWNERS file

1. If you would like the owner team to be required on all pull requests, add an entry to the CODEOWNERS file.  This team name should match the github team.  See the commented section for guidance.

### The Pull Request template

1. If you would like the team to be added as metadata to all pull requests thus allowing them to search for pull requests under their team.  See the commented section for guidance.

### This README

1. This README will be destroyed by the Template Cleanup and you will be responsible for adding information about the repository.
