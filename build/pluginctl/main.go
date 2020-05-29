// main handles deployment of the plugin to a development server using either the Client4 API
// or by copying the plugin bundle into a sibling mattermost-server/plugin directory.
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mholt/archiver/v3"
)

const helpText = `
Usage:
    pluginctl deploy <plugin id> <bundle path>
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
	case "reset":
		if client == nil {
			return errors.New("In order to reset, please set the following three environment variables:\n\n" +
				"MM_SERVICESETTINGS_SITEURL\nMM_ADMIN_USERNAME\nMM_ADMIN_PASSWORD\n\n" +
				"or, if using a token, set: MM_ADMIN_TOKEN")
		}
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

	if siteURL != "" {
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
	}
	return nil, nil
}

// deploy handles deployment of the plugin to a development server.
func deploy(client *model.Client4, pluginID, bundlePath string) error {
	if client != nil {
		return uploadPlugin(client, pluginID, bundlePath)
	}
	return copyPlugin(pluginID, bundlePath)
}

// uploadPlugin attempts to upload and enable a plugin via the Client4 API.
// It will fail if plugin uploads are disabled.
func uploadPlugin(client *model.Client4, pluginID, bundlePath string) error {
	pluginBundle, err := os.Open(bundlePath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", bundlePath, err)
	}
	defer pluginBundle.Close()

	log.Print("Uploading plugin via API.")
	_, resp := client.UploadPluginForced(pluginBundle)
	if resp.Error != nil {
		return fmt.Errorf("Failed to upload plugin bundle: %s", resp.Error.Error())
	}

	log.Print("Enabling plugin.")
	_, resp = client.EnablePlugin(pluginID)
	if resp.Error != nil {
		return fmt.Errorf("Failed to enable plugin: %s", resp.Error.Error())
	}

	return nil
}

// copyPlugin attempts to install a plugin by copying it to a sibling ../mattermost-server/plugin
// directory. A server restart is required before the plugin will start.
func copyPlugin(pluginID, bundlePath string) error {
	targetPath, _ := filepath.Abs("../mattermost-server")
	_, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		return errors.New("no supported deployment method available, please install plugin manually")
	} else if err != nil {
		return fmt.Errorf("failed to stat %s: %w", targetPath, err)
	}

	log.Printf("Installing plugin to mattermost-server found in %s.", targetPath)
	log.Print("Server restart required to load updated plugin.")

	targetPath = filepath.Join(targetPath, "plugins")

	err = os.MkdirAll(targetPath, 0777)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", targetPath, err)
	}

	existingPluginPath := filepath.Join(targetPath, pluginID)
	err = os.RemoveAll(existingPluginPath)
	if err != nil {
		return fmt.Errorf("failed to remove existing existing plugin directory %s: %w", existingPluginPath, err)
	}

	err = archiver.Unarchive(bundlePath, targetPath)
	if err != nil {
		return fmt.Errorf("failed to unarchive %s into %s: %w", bundlePath, targetPath, err)
	}

	return nil
}

// resetPlugin attempts to reset the plugin via the Client4 API.
func resetPlugin(client *model.Client4, pluginID string) error {
	log.Print("Disabling plugin.")
	_, resp := client.DisablePlugin(pluginID)
	if resp.Error != nil {
		return fmt.Errorf("failed to disable plugin: %s", resp.Error.Error())
	}

	log.Print("Enabling plugin.")
	_, resp = client.EnablePlugin(pluginID)
	if resp.Error != nil {
		return fmt.Errorf("failed to enable plugin: %s", resp.Error.Error())
	}

	return nil
}
