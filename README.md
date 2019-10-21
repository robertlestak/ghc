# ghc

A simple script to create repositories in GitHub. There are many other scripts / programs out there that provide a deeper API, but this script's benefit is its simplicity. No extraneous flags or options - just create your repository and get gitting.

## Build

````
go build .
cp ghc /usr/local/bin/ghc
````

## Configuration

To communicate with the GitHub API, you must add your GitHub access token to your `~/.netrc` file.

### Example .netrc

````
machine api.github.com login [token]
````

See here (https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) for more detailed steps to create an access token.

## Usage

### Flags

````
-c	Copy HTTPS URL to clipboard (bool; default: false)
-d	Description (string; default: nil)
-o	Organization (string; default: nil - must be a member)
-i	Initialize git repository
-p	Private Repository (bool; default: true)
-s	Copy ssh URL to clipboard (bool; default: false)
````

### Example

````
ghc -d "My description" -o "github" -s [name]
````

