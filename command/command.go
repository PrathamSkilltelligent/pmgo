package command

// import "github.com/samber/mo"

type CmdContext interface {
	IsCmdContext()
}

type CmdOutput interface {
	IsCmdOutput()
}

// type Command[T CmdOutput] interface {
// 	Execute(CmdContext) mo.Result[T]
// }

// func ExecuteCommand[T CmdOutput](ctx CmdContext, cmd Command[T]) mo.Result[T] {
// 	return cmd.Execute(ctx)
// }
