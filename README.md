## About
`AAC` a.k.a. `Article as Code` helps you collect articles from data sources, such as dev.to, and then stores them as static files. 
 It also helps you sync static files to create articles on your publications.

## Install
```shell
go install github.com/huantt/article-as-code@latest
```

Add `GOPATH/bin` directory to your PATH environment variable, so you can run Go programs anywhere.
```shell
export GOPATH=$HOME/go
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### Collect articles
```shell
Usage:
   collect [flags]

Flags:
  -f, --article-folder string   Article folder (default "data/articles")
  -h, --help                    help for collect-articles
      --rps int                 Limit concurrent requests (default 5)
  -u, --username string         Username
```

**For example**
```shell
aac collect \
--username=jacktt \
--rps=5 \
--article-folder=static
```

### Sync articles
```shell
Usage:
   sync [flags]

Flags:
  -f, --article-folder string   Article folder (default "data/articles")
  -a, --auth-token string       Auth token
  -h, --help                    help for sync-articles
      --rps int                 Limit concurrent requests (default 5)
  -u, --username string         Username
```

**For example**
```shell
aac sync \
--article-folder=data/articles \
--rps=5 \
--username=jacktt \
--auth-token=XXXXXXX
```

## Docker

```shell
docker run --rm huanttok/article-as-code:latest --help
```

## References
- [Hashnode API Docs](https://api.hashnode.com/)
- [Forem API Docs](https://developers.forem.com/api/)