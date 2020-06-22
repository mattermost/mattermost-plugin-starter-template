// main handles deployment of the plugin to a development server using the Client4 API.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/mattermost/mattermost-server/v5/model"
)

const helpText = `
Usage:
    pluginctl deploy <plugin id> <bundle path>
    pluginctl disable <plugin id>
    pluginctl enable <plugin id>
    pluginctl reset <plugin id>
`

func main() {
	err := pluginctl()
	if err != nil {
		fmt.Printf("Failed: %s\n", err.Error())
		fmt.Print(helpText)
		os.Exit(1)
	}
}

func pluginctl() error {
	if len(os.Args) < 3 {
		return errors.New("invalid number of arguments")
	}

	client, err := getClient()
	if err != nil {
		return err
	}

	switch os.Args[1] {
	case "deploy":
		if len(os.Args) < 4 {
			return errors.New("invalid number of arguments")
		}
		return deploy(client, os.Args[2], os.Args[3])
	case "disable":
		return disablePlugin(client, os.Args[2])
	case "enable":
		return enablePlugin(client, os.Args[2])
	case "reset":
		return resetPlugin(client, os.Args[2])
	default:
		return errors.New("invalid second argument")
	}
}

func getClient() (*model.Client4, error) {
	siteURL := os.Getenv("MM_SERVICESETTINGS_SITEURL")
	adminToken := os.Getenv("MM_ADMIN_TOKEN")
	adminUsername := os.Getenv("MM_ADMIN_USERNAME")
	adminPassword := os.Getenv("MM_ADMIN_PASSWORD")

	if siteURL == "" {
		return nil, errors.New("MM_SERVICESETTINGS_SITEURL is not set")
	}

	client := model.NewAPIv4Client(siteURL)

	if adminToken != "" {
		log.Printf("Authenticating using token against %s.", siteURL)
		client.SetToken(adminToken)
		return client, nil
	}

	if adminUsername != "" && adminPassword != "" {
		client := model.NewAPIv4Client(siteURL)
		log.Printf("Authenticating as %s against %s.", adminUsername, siteURL)
		_, resp := client.Login(adminUsername, adminPassword)
		if resp.Error != nil {
			return nil, fmt.Errorf("failed to login as %s: %w", adminUsername, resp.Error)
		}
		return client, nil
	}

	return nil, errors.New("one of MM_ADMIN_TOKEN or MM_ADMIN_USERNAME/MM_ADMIN_PASSWORD must be defined")
}

// deploy attempts to upload and enable a plugin via the Client4 API.
// It will fail if plugin uploads are disabled.
func deploy(client *model.Client4, pluginID, bundlePath string) error {
	pluginBundle, err := os.Open(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", bundlePath, err)
	}
	defer pluginBundle.Close()

	log.Print("Uploading plugin via API.")
	_, resp := client.UploadPluginForced(pluginBundle)
	if resp.Error != nil {
		return fmt.Errorf("failed to upload plugin bundle: %s", resp.Error.Error())
	}

	log.Print("Enabling plugin.")
	_, resp = client.EnablePlugin(pluginID)
	if resp.Error != nil {
		return fmt.Errorf("failed to enable plugin: %s", resp.Error.Error())
	}

	return nil
}

// disablePlugin attempts to disable the plugin via the Client4 API.
func disablePlugin(client *model.Client4, pluginID string) error {
	log.Print("Disabling plugin.")
	_, resp := client.DisablePlugin(pluginID)
	if resp.Error != nil {
		return fmt.Errorf("failed to disable plugin: %w", resp.Error)
	}

	return nil
}

// enablePlugin attempts to enable the plugin via the Client4 API.
func enablePlugin(client *model.Client4, pluginID string) error {
	log.Print("Enabling plugin.")
	_, resp := client.EnablePlugin(pluginID)
	if resp.Error != nil {
		return fmt.Errorf("failed to enable plugin: %w", resp.Error)
	}

	return nil
}

// resetPlugin attempts to reset the plugin via the Client4 API.
func resetPlugin(client *model.Client4, pluginID string) error {
	err := disablePlugin(client, pluginID)
	if err != nil {
		return err
	}

	err = enablePlugin(client, pluginID)
	if err != nil {
		return err
	}

	return nil
}
