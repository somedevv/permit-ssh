package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/boltdb/bolt"
	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func Interactive_mode(db *bolt.DB) error {
	utils.CallClear()

	answers := struct {
		Key          string
		Ip           string
		Confirmation string
	}{}

	err := survey.Ask(utils.SimpleConnection, &answers)
	if err != nil {
		return err
	}

	utils.CallClear()

	utils.PrintKeyandIP(answers.Key, answers.Ip)
	survey.AskOne(&utils.Prompt_confirmation, &answers.Confirmation, survey.WithValidator(survey.Required))

	if answers.Confirmation == "Yes" {
		utils.AddKey(answers.Ip, answers.Key)
	} else {
		colors.Red.Println("Key not added")
	}

	defer db.Close()
	return nil
}
