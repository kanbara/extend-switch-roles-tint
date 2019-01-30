# extend-switch-roles-tint
Colourify your AWS Extend Switch Roles

`go build tint/main.go`

copy and paste your extended role config to a file, run the command

 `./main from-file config.txt -g "*.dev" -g "*.int" -g "*.prd" -g "*poc*" -g "*.ci" --show`

 copy the config back, and you're good to go!
