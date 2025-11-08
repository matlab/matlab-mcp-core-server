// Copyright 2025 The MathWorks, Inc.

package processdetails

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	// The environment variable MW_CONTEXT_TAGS allows MathWorks to understand how this MathWorks product (MATLAB MCP Core Server)
	// is used with MATLAB. This information helps us improve our products.
	// Your content, and information about the content within your files, is not shared with MathWorks.
	// For more information see https://mathworks.com/support/faq/user_experience_information_faq.html
	MWContextTagsEnvVar = "MW_CONTEXT_TAGS"
	MWContextTagsValue  = "MATLAB:MATLAB_MCP_CORE_SERVER:V1"
)

type OSLayer interface {
	Environ() []string
}

type ProcessDetails struct {
	osLayer OSLayer
}

func New(
	osLayer OSLayer,
) *ProcessDetails {
	return &ProcessDetails{
		osLayer: osLayer,
	}
}

func (*ProcessDetails) NewAPIKey() string {
	return uuid.NewString()
}

func (g *ProcessDetails) EnvironmentVariables(sessionDirPath string, apiKey string, certificateFile string, certificateKey string) []string {
	processEnvVars := append(
		g.osLayer.Environ(),
		"MATLAB_LOG_DIR="+sessionDirPath,
		"MW_MCP_SESSION_DIR="+sessionDirPath,
		`MW_DIAGNOSTIC_DEST="filedir=`+sessionDirPath+`"`,
		"MWAPIKEY="+apiKey,
		"MW_CERTFILE="+certificateFile,
		"MW_PKEYFILE="+certificateKey,
	)

	// Add or update the MW_CONTEXT_TAGS environment variable
	// If MW_CONTEXT_TAGS already exists, append a comma and the new value.
	processEnvVars = func(processEnvVars []string) []string {
		envVarPrefix := MWContextTagsEnvVar + "="

		for i, envVar := range processEnvVars {
			if strings.HasPrefix(envVar, envVarPrefix) {
				// Found an existing MW_CONTEXT_TAGS environment variable
				envVarValue := strings.TrimPrefix(envVar, envVarPrefix)
				if envVarValue == "" {
					// Only the prefix was present, weird, but OK
					processEnvVars[i] = envVarPrefix + MWContextTagsValue
				} else {
					// It already has a value, comma separate the values
					processEnvVars[i] = fmt.Sprintf("%s,%s", envVar, MWContextTagsValue)
				}
				// Exit early, we assume only one entry per environment variable
				return processEnvVars
			}
		}

		// We didn't find an existing MW_CONTEXT_TAGS environment variable
		return append(processEnvVars, envVarPrefix+MWContextTagsValue)
	}(processEnvVars)

	return processEnvVars
}

func (*ProcessDetails) StartupFlag(os string, showMATLAB bool, startupCode string) []string {
	startupFlags := []string{}
	if showMATLAB {
		startupFlags = append(startupFlags,
			"-desktop",
		)
	} else {
		startupFlags = append(startupFlags,
			"-nosplash",
			"-softwareopengl",
			"-nodesktop",
		)
		if os == "windows" {
			startupFlags = append(startupFlags,
				"-noDisplayDesktop",
				"-wait",
				"-log",
				"/minimize",
			)
		} else {
			// Unix platforms (Linux/macOS)
			startupFlags = append(startupFlags,
				"-minimize",
			)
		}
	}
	startupFlags = append(startupFlags,
		"-r",
		startupCode,
	)
	return startupFlags
}
