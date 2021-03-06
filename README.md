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