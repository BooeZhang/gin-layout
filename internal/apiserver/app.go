package apiserver

import (
	"github.com/BooeZhang/gin-layout/internal/pkg/app"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/pkg/log"
)

// NewApp 创建应用程序
func NewApp() *app.App {
	opts := options.DefaultOption()
	application := app.NewApp(
		app.WithOptions(opts),
		app.WithRunFunc(run(opts)),
	)

	return application
}

// run 创建应用程序的回调函数
func run(opts *options.Options) app.RunFunc {
	return func() error {
		log.Init(opts.Log)
		defer log.Flush()

		fn := func(opts *options.Options) error {
			server, err := createAPIServer(opts)
			if err != nil {
				return err
			}

			return server.Run()
		}

		return fn(opts)
	}
}
