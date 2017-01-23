# Gin-Logrus

[Gin](https://github.com/gin-gonic/gin) middleware for Logging with
[Logrus](https://github.com/sirupsen/logrus). It is meant as drop in
replacement for the default logger used in Gin.

## Requirements

- github.com/gin-gonic/gin
- github.com/sirupsen/logrus

## Installation


```
go get github.com/rocksolidlabs/gin-logrus
```

### Example

```go
package main
import (
    "flag"
    "time"
    log "github.com/sirupsen/logrus"
    "github.com/rocksolidlabs/gin-logrus"
    "github.com/gin-gonic/gin"
)

func main() {
    flag.Parse()
    router := gin.New()
    router.Use(ginlogrus. Logger("MYAPI", false, true, os.Stdout, logrus.WarnLevel))
    //..
    router.Use(gin.Recovery())
    log.Info("API Running")
    router.Run(":8080")
}
```

## Derived from

https://github.com/zalando/gin-glog

## License

See [LICENSE](LICENSE) file.
