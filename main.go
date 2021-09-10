package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mpkg/pkg"
	"mpkg/utils"
	"os"
	"os/exec"
	"strings"

	"github.com/gosuri/uitable"
	"github.com/kelindar/binary"
	"github.com/mingrammer/cfmt"
)

type Command struct {
	Name string
	Desc string
	Func func()
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

	info := Command{
		Name: "info",
		Desc: "Get information about a package",
		Func: Info,
	}

	Commands["help"] = help
	Commands["install"] = install
	Commands["remove"] = remove
	Commands["search"] = search
	Commands["info"] = info
}

func Info() {
	if len(os.Args) > 2 {
		if strings.HasPrefix(os.Args[2], "./") {
			infoLocal(os.Args)
			return
		}
	}
}

func infoLocal(args []string) {
	file, err := os.Open(args[2])
	if err != nil {
		cfmt.Errorln("Unable to open ", args[2])
		os.Exit(1)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		cfmt.Errorln("Unable to read ", args[2])
		os.Exit(1)
	}

	mpkg := pkg.Package{}

	binary.Unmarshal(fileBytes, &mpkg)

	fmt.Printf("mpkg Version: %s\n", mpkg.Header.MpkgVersion)
	fmt.Printf("Package Name: %s\n", mpkg.Meta.Name)
	fmt.Printf("Package Version: %s\n", mpkg.Meta.Version)
	fmt.Println("Package Dependencies:")
	for _, d := range mpkg.Meta.Dependencies {
		fmt.Printf("\t%s\n", d)
	}
	fmt.Println("Package Install Commands:")
	for _, s := range mpkg.Meta.Scripts {
		fmt.Printf("\t%s\n", s)
	}
	cs := len(mpkg.Content.Files)
	fmt.Printf("Package Content Size: %d bytes\n", cs)
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
			cfmt.Infoln("Found package '" + os.Args[2] + "' in repo")
		} else {
			cfmt.Warningln("Unable to find pacakeg '" + os.Args[2] + "' in repo")
		}
	}
}

func Install() {
	//reader := bufio.NewReader(os.Stdin)
	if len(os.Args) > 2 {
		if strings.HasPrefix(os.Args[2], "./") {
			installLocal(os.Args)
			return
		}
		/*if packageExists(os.Args[2]) && !packageInstalled(os.Args[2]) {
			fmt.Print("Are you sure you want to install ", os.Args[2], "? ", chalk.Magenta.Color("(yes, no) "))
			response, _ := reader.ReadString('\n')
			switch response {
			case "y":
				installPackage(os.Args[2])
			case "yes\n":
				installPackage(os.Args[2])
			case "n":
				os.Exit(0)
			case "no":
				os.Exit(0)
			default:
				os.Exit(0)
			}
			//
		} else if packageInstalled(os.Args[2]) {

		}*/
	} else {
		installHelp()
	}
}

func installLocal(args []string) {
	file, err := os.Open(args[2])
	if err != nil {
		cfmt.Errorln("Unable to open ", args[2])
		os.Exit(1)
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		cfmt.Errorln("Unable to read ", args[2])
		os.Exit(1)
	}

	mpkg := pkg.Package{}

	binary.Unmarshal(fileBytes, &mpkg)

	installPackage(&mpkg)
}

func installPackage(mpkg *pkg.Package) {
	//cfmt.Successln("Successfully installed: " + packageName)
	tarBall, err := os.Create("content.tar")
	if err != nil {
		cfmt.Errorln("Unable to create content.tar: ", err)
		os.Exit(1)
	}

	tarBytes := bytes.NewReader(mpkg.Content.Files)

	if _, err := io.Copy(tarBall, tarBytes); err != nil {
		cfmt.Errorln("Unable to write to content.tar: ", err)
		os.Exit(1)
	}

	utils.UntarFile("content.tar")

	for _, c := range mpkg.Meta.Scripts {
		cmd := exec.Command("bash", "-c", c)
		stdout, err := cmd.Output()
		if err != nil {
			cfmt.Errorln("Unable to run Setup commands: ", err)
			os.Exit(1)
		}
		fmt.Println(string(stdout))
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
