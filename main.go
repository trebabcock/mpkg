package main

import (
	"os"
	"fmt"

	"github.com/mingrammer/cfmt"
	"github.com/gosuri/uitable"
)

type Command struct {
	Name	string
	Desc	string
	Func	func()
}

var Commands = make(map[string]Command)

func main() {
	setCommands()

	if val, ok := Commands[os.Args[1]]; ok {
		val.Func()
	} else {
		cfmt.Errorln("Unknown command: " + os.Args[1])
		fmt.Println("Run 'mpkg help' for usage.")
	}
}

func setCommands() {
	help := Command{
		Name: "help",
		Desc: "Show this help message",
		Func: Usage,
	}

	install := Command{
		Name: "install",
		Desc: "Install a package",
		Func: Install,
	}

	remove := Command{
		Name: "remove",
		Desc: "Remove a package",
		Func: Remove,
	}

	search := Command{
		Name: "search",
		Desc: "Search for a package in the repo",
		Func: Search,
	}

	Commands["help"] = help
	Commands["install"] = install
	Commands["remove"] = remove
	Commands["search"] = search
}

func Usage() {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true

	table.AddRow("COMMAND", "DESCRIPTION")
	table.AddRow("", "")
	for _, v := range Commands {
		table.AddRow(v.Name, v.Desc)
	}

	fmt.Println(table)
}

func Search() {
	if len(os.Args) > 2 {
		if packageExists(os.Args[2]) {
			cfmt.Infoln("Found package '"+os.Args[2]+"' in repo")
		} else {
			cfmt.Warningln("Unable to find pacakeg '" + os.Args[2] + "' in repo")
		}
	}
}

func Install() {
	if len(os.Args) > 2 {
		if packageExists(os.Args[2]) && !packageInstalled(os.Args[2]) {
			cfmt.Successln("Successfully installed: " + os.Args[2])
		} else if packageInstalled(os.Args[2]) {
			
		}
	} else {
		installHelp()
	}
}

func installHelp() {
	fmt.Println("Usage: mpkg install <package-name(s)>")
}

func Remove() {
	
}

func removeHelp() {
	fmt.Println("Usage: mpkg remove <package-name(s)>")
}

func packageInstalled(name string) bool {
	return false
}

func packageExists(name string) bool {
	return true
}
