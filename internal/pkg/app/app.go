package app

import (
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
)

var (
	progressMessage = color.GreenString("==>")
)

// App 应用
type App struct {
	basename string
	options  *options.Options
	runFunc  RunFunc
	silence  bool
	noConfig bool
	cmd      *cobra.Command
}

// Option 使用配置项
type Option func(*App)

// WithOptions 用于初始化应用程序配置项。
func WithOptions(opt *options.Options) Option {
	return func(a *App) {
		a.options = opt
	}
}

// RunFunc 服务启动后的回调函数
type RunFunc func() error

// WithRunFunc 用于设置应用程序启动回调函数选项。
func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

// WithName 设置服务名字
func WithName(name string) Option {
	return func(app *App) {
		app.basename = name
	}
}

// NewApp 根据给定的应用程序名称和其他选项创建一个新的应用程序实例。
func NewApp(opts ...Option) *App {
	a := &App{}

	for _, o := range opts {
		o(a)
	}
	if len(a.basename) == 0 && a.basename == "" {
		a.basename = "api-server"
	}
	a.buildCommand()
	return a
}

// buildCommand 应用程序命令行程序
func (a *App) buildCommand() {
	cmd := cobra.Command{
		Use:           FormatBaseName("api-server"),
		Short:         "API服务",
		Long:          "API服务",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	if !a.noConfig {
		options.AddConfigFlag(a.basename, cmd.Flags())
	}
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	a.cmd = &cmd
}

// runCommand 应用程序命令行程序
func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, "api-server")
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}
	// run application
	if a.runFunc != nil {
		return a.runFunc()
	}

	return nil
}

// Run 启动应用程序。
func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// FormatBaseName 生成不同操作系统下的可执行文件名。
func FormatBaseName(basename string) string {
	// Make case-insensitive and strip executable suffix if present
	if runtime.GOOS == "windows" {
		basename = strings.ToLower(basename)
		basename = strings.TrimSuffix(basename, ".exe")
	}

	return basename
}

// printWorkingDir 打印工作目录
func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}
