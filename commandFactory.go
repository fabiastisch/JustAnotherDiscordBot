package main

type CommandFactory struct {
	commands map[string]Command
}

func newCommandFactory() *CommandFactory {
	commandFactory := &CommandFactory{}
	commandFactory.commands = make(map[string]Command)
	commandFactory.commands["/ping"] = &PingCommand{}
	commandFactory.commands["/krossekrabbe"] = &PatrickCommand{}
	return commandFactory
}

func (commandFactory *CommandFactory) getCommand(commandKey string) Command {
	return commandFactory.commands[commandKey]
}

func (commandFactory *CommandFactory) setCommand(commandKey string, command Command) {
	commandFactory.commands[commandKey] = command
}
