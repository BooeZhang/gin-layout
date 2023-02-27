package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/BooeZhang/gin-layout/internal/apiserver"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}

func main() {
	apiserver.NewApp("api-server").Run()
}
