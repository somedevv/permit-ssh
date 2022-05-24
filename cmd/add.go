package cmd

import (
	"os"

	"github.com/somedevv/permit-ssh/colors"
	"github.com/somedevv/permit-ssh/utils"
)

func AddWithIP(ip, key string) {
	utils.AddKey(ip, key)
}

func AddWithAWS(profile, region, instance, key string) {
	if profile == "" && region == "" {
		colors.Red.Println("Error: At least AWS profile or region must be set")
		os.Exit(1)
	}
	ip := utils.GetAWSInstance(profile, region, instance, key)
	if ip == "" {
		colors.Red.Println("Error: No instance found")
		os.Exit(1)
	}
	utils.AddKey(ip, key)
}
