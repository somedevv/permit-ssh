package cmd

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func InteractiveLocal(db *bolt.DB) {
	utils.CallClear()

	answers := struct {
		Key          string
		Ip           string
		Confirmation string
	}{}

	err := survey.Ask(utils.SimpleConnection, &answers)
	if err != nil {
		colors.Red.Println(err)
		os.Exit(1)
	}

	utils.CallClear()

	utils.PrintKeyandIP(answers.Key, answers.Ip)
	err = survey.AskOne(&utils.Prompt_confirmation, &answers.Confirmation, survey.WithValidator(survey.Required))
	if err != nil {
		colors.Red.Println(err)
		os.Exit(1)
	}

	if answers.Confirmation == "Yes" {
		utils.AddKey(answers.Ip, answers.Key)
	} else {
		colors.Red.Println("Key not added")
	}

	defer db.Close()
	os.Exit(0)
}
