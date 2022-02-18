## Resumify
Generater tool in Go to build your awesome CV :page_facing_up:

## Why
Reason is to generate simple and clean site based on most popular configuration extension YML, where people can easily adds new experiences and do it via VCS.

## Usage
Make your own repository on github eg. https://github.com/{username}/cv <br />
Get file in raw https://raw.githubusercontent.com/{username}/cv/{branch}/cv.yml (this will be request for Go app) <br />
Put this file into .env under the key "src"  <br />
Run resumify app and check your resume on endpoint "/preview" <br />
If u like that, you can easily generate index.html on endpoint "/generate" <br />
