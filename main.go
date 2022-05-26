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

	// AWS FLAG VARIABLES
	profile  string
	region   string
	instance string

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

	//AWS
	awsset = flaggy.NewSubcommand("aws")
	awsset.String(&profile, "p", "profile", "AWS Profile")
	awsset.String(&region, "r", "region", "AWS Profile")
	awsset.String(&instance, "i", "instance", "AWS Instance name")

	//---------SUBCOMMANDS---------//

	//------DELETE------//
	remove = flaggy.NewSubcommand("remove")
	remove.String(&user, "u", "user", "The user to remove")
	remove.String(&key, "k", "key", "The key to remove")
	remove.String(&ip, "ip", "address", "The IP of the server to remove the user")
	remove.AttachSubcommand(awsset, 1)
	flaggy.AttachSubcommand(remove, 1)

	//------ADD------//
	add = flaggy.NewSubcommand("add")
	add.String(&user, "u", "user", "The user to add")
	add.String(&key, "k", "key", "The key to add")
	add.String(&ip, "ip", "address", "The IP of the server to add the user")
	add.AttachSubcommand(awsset, 1)
	flaggy.AttachSubcommand(add, 1)

	//------LIST------//
	list = flaggy.NewSubcommand("list")
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
				colors.Red.Println("Error: At least AWS profile or region must be set")
				os.Exit(1)
			}
			cmd.ListAWS(profile, region)

		} else {
			cmd.ListLocal(db)
		}
	}

	if remove.Used {
		if ip != "" && awsset.Used == false {
			utils.RemoveKeyFromLocalDB(db, user, key)
		}

		if user != "" && key == "" {
			key = utils.SearchUserInLocalDB(db, user)
			if key == "" {
				colors.Red.Printf("User [%s] not found\n", user)
				os.Exit(1)
			}
		}

		if ip != "" && key != "" {
			cmd.DeleteWithIP(ip, key)
		} else if awsset.Used == true && key != "" {
			if profile == "" && region == "" {
				colors.Red.Println("Error: At least AWS profile or region must be set")
				os.Exit(1)
			}
			cmd.DeleteWithAWS(profile, region, instance, key)
		} else {
			colors.Red.Println("Error: At least IP or AWS Instance with user or key must be set")
			os.Exit(1)
		}
	}

	if add.Used {
		if user != "" && key != "" {
			utils.SaveKeyInLocalDB(db, user, key, ip)
		} else if user != "" && key == "" {
			key = utils.SearchUserInLocalDB(db, user)
		}

		if awsset.Used == true && instance != "" && key != "" {
			if profile == "" && region == "" {
				colors.Red.Println("Error: At least AWS profile or region must be set")
				os.Exit(1)
			}
			cmd.AddWithAWS(profile, region, instance, key)

		} else if ip != "" && key != "" {
			cmd.AddWithIP(ip, key)
		} else {
			colors.Red.Println("Error: At least IP or AWS Instance with user or key must be set")
			os.Exit(1)
		}
	}

	defer db.Close()
	os.Exit(0)
}
