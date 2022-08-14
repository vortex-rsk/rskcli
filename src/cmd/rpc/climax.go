package rpc

import (
	"fmt"
	"rskcli/src/utils"
	"rskcli/src/utils/color"
)

type Command struct {
	Names       []string
	Usage       string
	Params      []DefParam
	Description string
	Run         func(*Command, *Context)
}

type Handler struct {
	Methods []Method
	Run     func(*Method, *Context)
}

type Context struct {
	Data        map[string]string
	CommandArgs []string
	Flags       map[string]bool
}

var commands []*Command

var handlers []*Handler

var context map[string]string = make(map[string]string)
var commandArgs []string
var flags map[string]bool

func Handle(args []string) {

	cmd := args[0]

	//cmd := GetCommand(args[0])

	ctx := &Context{context, commandArgs, flags}

	var found bool = false

	// search in the generic handlers for RPC methods
	for _, handler := range handlers {
		for _, method := range handler.Methods {
			if utils.IndexOf(method.Names, cmd) >= 0 {
				found = true
				handler.Run(&method, ctx)
				return
			}
		}
	}

	if !found {
		fmt.Printf("Command %s does not exist. You need help.\n", color.Red(cmd))
	}

}

func AddToContext(key string, val string) {
	context[key] = val
}
func AddFlags(data *map[string]bool) {
	flags = *data
}

func AddArgsToContext(args []string) {
	commandArgs = args
}

func AddCommand(cmd *Command) {
	commands = append(commands, cmd)
}

func GetCommand(name string) *Command {
	for _, cmd := range commands {
		if utils.IndexOf(cmd.Names, name) >= 0 {
			return cmd
		}
	}
	return nil
}

func (ctx Context) Get(key string) string {
	return ctx.Data[key]
}

func AddHandler(handler *Handler) {
	handlers = append(handlers, handler)
}
