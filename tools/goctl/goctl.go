package main

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/cmd"
)

func main() {
	fmt.Println("[Local Run]")

	logx.Disable()
	load.Disable()
	cmd.Execute()
}
