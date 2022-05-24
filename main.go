package main

import (
	"os"

	"github.com/integrii/flaggy"
	"github.com/somedevv/permit-ssh/cmd"
	"github.com/somedevv/permit-ssh/conf"
	"github.com/somedevv/permit-ssh/utils"
)

// Version of the program, set at buildtime with -ldflags "-X main.version=X"
var version = ""

var (
	// FLAG VARIABLES
	user string
	key  string
	ip   string

	// SUBCOMMANDS
	remove      *flaggy.Subcommand
	add         *flaggy.Subcommand
	list        *flaggy.Subcommand
	interactive *flaggy.Subcommand

	// Configuration
	config conf.Config
)

func init() {
	// Load the config file
	config.GetConf()

	//------META------//
	flaggy.SetName("permit")
	flaggy.SetDescription("Your own SSH key manager and friend, made by somedevv")
	flaggy.SetVersion(version)
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/somedevv/permit-ssh"

	//---------SUBCOMMANDS---------//

	//------DELETE------//
	remove = flaggy.NewSubcommand("remove")
	remove.String(&user, "u", "user", "The user to remove")
	remove.String(&key, "k", "key", "The key to remove")
	remove.String(&ip, "ip", "address", "The IP of the server to remove the user")
	flaggy.AttachSubcommand(remove, 1)

	//------ADD------//
	add = flaggy.NewSubcommand("add")
	add.String(&user, "u", "user", "The user to add")
	add.String(&key, "k", "key", "The key to add")
	add.String(&ip, "ip", "address", "The IP of the server to add the user")
	flaggy.AttachSubcommand(add, 1)

	//------LIST------//
	list = flaggy.NewSubcommand("list")
	flaggy.AttachSubcommand(list, 1)

	//------INTERACTIVE------//
	interactive = flaggy.NewSubcommand("interactive")
	flaggy.AttachSubcommand(interactive, 1)

	//------PARSE------//
	flaggy.Parse()
}

func main() {
	if config.DB == "local" {
		RunWithLocalDB()
	}
}

func RunWithLocalDB() {
	db := utils.SetupLocalDB()

	if interactive.Used {
		cmd.InteractiveLocal(db)
	}

	if list.Used {
		cmd.ListLocal(db)
	}

	if remove.Used {
		cmd.RemoveLocal(db, user, key, ip)
	}

	if add.Used {
		cmd.AddLocal(db, user, key, ip)
	}

	defer db.Close()
	os.Exit(0)
}

func AWS() {

}
