# vers
Version file management

To Build cd to this directory and type 

```# go build```


``` ./vers --help
Handle versions....

Usage:
  vers [command]

Available Commands:
  bump        increment either major, minor or patch version number
  delete      delete an entry for version file.
  get         get version info
  help        Help about any command
  init        Make a new version file
  set         Add a new entry to version file

Flags:
      --config string         config file (default is $HOME/.vers.yaml)
  -d, --debug                 Turn on debug messages
  -e, --entry string          Which entry in version file
  -h, --help                  help for vers
  -M, --major int             major number (default: 0)
  -m, --minor int             minor number (default: 0)
  -p, --patch int             patch number  (default 1)
      --prefix string         prefix  (default "v")
      --suffix string         suffix
  -f, --version-file string   version file to use

Use "vers [command] --help" for more information about a command.
```
