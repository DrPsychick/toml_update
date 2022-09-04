[![CircleCI](https://img.shields.io/circleci/build/github/DrPsychick/toml_update)](https://app.circleci.com/pipelines/github/DrPsychick/toml_update)
[![Coverage Status](https://coveralls.io/repos/github/DrPsychick/toml_update/badge.svg?branch=main)](https://coveralls.io/github/DrPsychick/toml_update?branch=master)
[![license](https://img.shields.io/github/license/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update)
[![Contributors](https://img.shields.io/github/contributors/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update/graphs/contributors)
[![Paypal](https://img.shields.io/badge/donate-paypal-00457c.svg?logo=paypal)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=FTXDN7LCDWUEA&source=url)
[![GitHub Sponsor](https://img.shields.io/badge/github-sponsor-blue?logo=github)](https://github.com/sponsors/DrPsychick)

[![GitHub issues](https://img.shields.io/github/issues/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update/issues?q=is%3Aissue+is%3Aclosed)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update/pulls)
[![GitHub closed pull requests](https://img.shields.io/github/issues-pr-closed/drpsychick/toml_update.svg)](https://github.com/drpsychick/toml_update/pulls?q=is%3Apr+is%3Aclosed)

# toml_update
Update a toml configuration file from ENV variables

## Example
This example uses a `telegraf` configuration file, same works for other projects using the same configuration format
```shell
docker run --rm telegraf telegraf config > test.conf
export CONF_UPDATE=test.conf
export CONF_PREFIX=PFX

export PFX_myvar1=inputs.cpu.totalcpu=false
export PFX_whatever2=agent.interval="20s"
export PFX_ping_urls=inputs.ping.urls=["google.com","amazon.com"]
export PFX_emptysection=inputs.netstat.enabled=#no configuration

grep -e 'totalcpu' -e '^[^#]*interval = ' test.conf
# if you've cloned the repo, simply run `go build toml_update.go` to get the binary
./toml_update
grep -e 'totalcpu' -e '^[^#]*interval = ' test.conf
```

## Built and released with [goreleaser](https://goreleaser.com)
### Setup
```shell
docker run --rm --privileged \
  -v $PWD:/go/src/github.com/drpsychick/toml_update \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /go/src/github.com/drpsychick/toml_update \
  goreleaser/goreleaser init
  
go mod init github.com/drpsychick/toml_update
```

### Test
```shell
docker run --rm --privileged \
  -v $PWD:/go/src/github.com/drpsychick/toml_update \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /go/src/github.com/drpsychick/toml_update \
  goreleaser/goreleaser --snapshot --skip-publish --rm-dist
```

### Release
```shell
git tag -a v0.0.8 -m "Reviewed"
git push origin v0.0.8

docker run --rm --privileged \
  -v $PWD:/go/src/github.com/drpsychick/toml_update \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /go/src/github.com/drpsychick/toml_update \
  -e GITHUB_TOKEN=XXX \
  goreleaser/goreleaser release --rm-dist
```