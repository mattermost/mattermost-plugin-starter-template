// main handles deployment of the plugin to a development server using either the Client4 API
// or by copying the plugin bundle into a sibling mattermost-server/plugin directory.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mholt/archiver/v3"
	"github.com/pkg/errors"
)

func main() {
	err := deploy()
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}

// deploy handles deployment of the plugin to a development server.
func deploy() error {
	if len(os.Args) < 3 {
		return errors.New("deploy <plugin id> <bundle path>")
	}

	pluginId := os.Args[1]
	bundlePath := os.Args[2]

	siteURL := os.Getenv("MM_SERVICESETTINGS_SITEURL")
	adminUsername := os.Getenv("MM_ADMIN_USERNAME")
	adminPassword := os.Getenv("MM_ADMIN_PASSWORD")
	copyTargetDirectory, _ := filepath.Abs("../mattermost-server")
	if siteURL != "" && adminUsername != "" && adminPassword != "" {
		return uploadPlugin(pluginId, bundlePath, siteURL, adminUsername, adminPassword)
	} else if _, err := os.Stat(copyTargetDirectory); err == nil {
		log.Printf("Installing plugin to mattermost-server found in %s.", copyTargetDirectory)
		log.Print("Server restart and manual plugin enabling required.")
		return copyPlugin(pluginId, copyTargetDirectory, bundlePath)
	} else {
		p, _ := filepath.Abs("../mattermost-server")
		log.Print("ERROR", err.Error(), p)
		return errors.New("No supported deployment method available. Install plugin manually.")
	}

	return nil
}

// uploadPlugin attempts to upload and enable a plugin via the Client4 API.
// It will fail if plugin uploads are disabled.
func uploadPlugin(pluginId, bundlePath, siteURL, adminUsername, adminPassword string) error {
	client := model.NewAPIv4Client(siteURL)
	log.Printf("Authenticating as %s against %s.", adminUsername, siteURL)
	_, resp := client.Login(adminUsername, adminPassword)
	if resp.Error != nil {
		return fmt.Errorf("Failed to login as %s: %s", adminUsername, resp.Error.Error())
	}

	pluginBundle, err := os.Open(bundlePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", bundlePath)
	}
	defer pluginBundle.Close()

	log.Print("Uploading plugin via API.")
	_, resp = client.UploadPluginForced(pluginBundle)
	if resp.Error != nil {
		return fmt.Errorf("Failed to upload plugin bundle: %s", resp.Error.Error())
	}

	log.Print("Enabling plugin.")
	_, resp = client.EnablePlugin(pluginId)
	if resp.Error != nil {
		return fmt.Errorf("Failed to enable plugin: %s", resp.Error.Error())
	}

	return nil
}

// copyPlugin attempts to install a plugin by copying it to a sibling ../mattermost-server/plugin
// directory. A server restart is required before the plugin will start.
func copyPlugin(pluginId, targetPath, bundlePath string) error {
	targetPath = filepath.Join(targetPath, "plugins")

	err := os.MkdirAll(targetPath, 0777)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", targetPath)
	}

	existingPluginPath := filepath.Join(targetPath, pluginId)
	err = os.RemoveAll(existingPluginPath)
	if err != nil {
		return errors.Wrapf(err, "failed to remove existing existing plugin directory %s", existingPluginPath)
	}

	err = archiver.Unarchive(bundlePath, targetPath)
	if err != nil {
		return errors.Wrapf(err, "failed to unarchive %s into %s", bundlePath, targetPath)
	}

	return nil
}
