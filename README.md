[![Build Status](https://travis-ci.org/Schtolc/mooncore.svg?branch=master)](https://travis-ci.org/Schtolc/mooncore)[![Go Report Card](https://goreportcard.com/badge/github.com/Schtolc/mooncore)](https://goreportcard.com/report/github.com/Schtolc/mooncore)

# Mooncore
Backend for moon project written in go. Under heavy development atm.

### How to build project
1. `mkdir -p $GOPATH/src/github.com/Schtolc`
1. `git clone git@github.com:Schtolc/mooncore.git`
2. `go build github.com/Schtolc/mooncore`

### How to launch project
1. install & run mysql-server on default port
2. specify mysql credentials in `config.yaml`
3. `echo 'CREATE DATABASE IF NOT EXISTS moondb' | mysql -u root -p`
4. `sudo mkdir /var/log/mooncore`
5. `sudo chmod 777 /var/log/mooncore`
6. `$GOPATH/src/github.com/Schtolc/mooncore/mooncore`

### How to deploy project
1. TODO @Schtolc #4

### How to develop project
1. Submit an issue
2. Submit a pull request for an issue with "[MCORE-{ISSUE NUMBER}]" suffix
3. Make changes until an approval
4. Merge pull request
5. Close issue

### Development rules
If you pushing any changes, make sure that:
1. Travis is green (see next section)
2. All changes are directly related to an issue(!)
3. Your changes response to all (if any) comments in the pull request
4. All commits have "[MCORE-{ISSUE NUMBER}]" suffix

### Tools we use during development and in Travis
1. [golint](https://github.com/golang/lint): `golint ./...`
2. [gocyclo](https://github.com/fzipp/gocyclo): `gocyclo --over 10 .`
3. gofmt: `gofmt -l .`