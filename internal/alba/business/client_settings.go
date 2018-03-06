package business

import "github.com/spf13/viper"

/*
This package exposes the data and operations regarding settings available from the client.

All app-wide settings are stored in the alba.yml file.
 */

type Settings struct {
	LibraryPath string
	CoversPreferredSource string

	DisableLibraryConfiguration bool
}

type SettingsInteractor struct {}

func (si *SettingsInteractor) GetSettings() Settings {
	var settings Settings

	settings.DisableLibraryConfiguration = viper.GetBool("ClientSettings.DisableLibraryConfiguration")
	settings.LibraryPath = viper.GetString("Library.Path")
	settings.CoversPreferredSource = viper.GetString("Covers.PreferredSource")

	return settings
}

