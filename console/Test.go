package console

import (
	"fmt"
	"github.com/ctfang/command"
	"os"
	"os/signal"
	"syscall"
)

type Test struct {
}

func (Test) Configure() command.Configure {
	return command.Configure{
		Name:        "test",
		Description: "测试命令",
		Input:       command.Argument{},
	}
}

func (Test) Execute(input command.Input) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("进程信号是")
		fmt.Println(sig)
	}()

	fmt.Println("awaiting signal")
}
