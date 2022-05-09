package main

import "github.com/AlecAivazis/survey/v2"

var SimpleConnection = []*survey.Question{
	{
		Name:      "key",
		Prompt:    &survey.Input{Message: "What is the public key?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:      "ip",
		Prompt:    &survey.Input{Message: "What is the server IP?"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
}

var prompt_confirmation = *&survey.Select{
	Message: "Add RSA key?",
	Options: []string{"Yes", "No"},
}
