package main

import (
	"os"

	"github.com/integrii/flaggy"
	"github.com/somedevv/permit-ssh/cmd"
	"github.com/somedevv/permit-ssh/colors"
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
	aws  bool

	// AWS FLAG VARIABLES
	profile string
	region  string

	// SUBCOMMANDS
	remove      *flaggy.Subcommand
	add         *flaggy.Subcommand
	list        *flaggy.Subcommand
	interactive *flaggy.Subcommand
	awsset      *flaggy.Subcommand

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

	//------NESTED SUBCOMMANDS------//
	awsset = flaggy.NewSubcommand("aws")
	awsset.String(&profile, "p", "profile", "AWS Profile")
	awsset.String(&region, "r", "region", "AWS Profile")

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
	list.Bool(&aws, "aws", "aws", "Use AWS CLI")
	list.AttachSubcommand(awsset, 1)
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
		if awsset.Used == true {
			if profile == "" && region == "" {
				colors.Red.Println("Error: Atleast AWS profile or region must be set")
				os.Exit(1)
			}
			cmd.ListAWS(profile, region)

		} else {
			cmd.ListLocal(db)
		}
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
