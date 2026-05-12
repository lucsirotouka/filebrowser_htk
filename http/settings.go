package fbhttp

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/filebrowser/filebrowser/v2/rules"
	"github.com/filebrowser/filebrowser/v2/settings"
)

type settingsData struct {
	Signup                  bool                  `json:"signup"`
	HideLoginButton         bool                  `json:"hideLoginButton"`
	CreateUserDir           bool                  `json:"createUserDir"`
	MinimumPasswordLength   uint                  `json:"minimumPasswordLength"`
	UserHomeBasePath        string                `json:"userHomeBasePath"`
	Defaults                settings.UserDefaults `json:"defaults"`
	AuthMethod              settings.AuthMethod   `json:"authMethod"`
	Rules                   []rules.Rule          `json:"rules"`
	Branding                settings.Branding     `json:"branding"`
	Tus                     settings.Tus          `json:"tus"`
	Shell                   []string              `json:"shell"`
	Commands                map[string][]string   `json:"commands"`
	AllowedUploadExtensions []string              `json:"allowedUploadExtensions"`
}

var settingsGetHandler = withAdmin(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	data := &settingsData{
		Signup:                  d.settings.Signup,
		HideLoginButton:         d.settings.HideLoginButton,
		CreateUserDir:           d.settings.CreateUserDir,
		MinimumPasswordLength:   d.settings.MinimumPasswordLength,
		UserHomeBasePath:        d.settings.UserHomeBasePath,
		Defaults:                d.settings.Defaults,
		AuthMethod:              d.settings.AuthMethod,
		Rules:                   d.settings.Rules,
		Branding:                d.settings.Branding,
		Tus:                     d.settings.Tus,
		Shell:                   d.settings.Shell,
		Commands:                d.settings.Commands,
		AllowedUploadExtensions: d.settings.AllowedUploadExtensions,
	}

	return renderJSON(w, r, data)
})

// normalizeExtensions converts a list of extension strings to lowercase,
// ensuring each has a leading dot (e.g. "JPG" -> ".jpg", ".PNG" -> ".png").
// Empty strings are dropped.
func normalizeExtensions(exts []string) []string {
	if len(exts) == 0 {
		return exts
	}
	result := make([]string, 0, len(exts))
	for _, e := range exts {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		// Ensure the extension has a leading dot.
		e = strings.ToLower(e)
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		result = append(result, e)
	}
	return result
}

var settingsPutHandler = withAdmin(func(_ http.ResponseWriter, r *http.Request, d *data) (int, error) {
	req := &settingsData{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return http.StatusBadRequest, err
	}

	d.settings.Signup = req.Signup
	d.settings.CreateUserDir = req.CreateUserDir
	d.settings.MinimumPasswordLength = req.MinimumPasswordLength
	d.settings.UserHomeBasePath = req.UserHomeBasePath
	d.settings.Defaults = req.Defaults
	d.settings.Rules = req.Rules
	d.settings.Branding = req.Branding
	d.settings.Tus = req.Tus
	d.settings.Shell = req.Shell
	d.settings.Commands = req.Commands
	d.settings.HideLoginButton = req.HideLoginButton
	d.settings.AllowedUploadExtensions = normalizeExtensions(req.AllowedUploadExtensions)

	err = d.store.Settings.Save(d.settings)
	return errToStatus(err), err
})
