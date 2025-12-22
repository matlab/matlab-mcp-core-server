// Copyright 2025 The MathWorks, Inc.

package parser

import "github.com/spf13/pflag"

func GenerateUsageText(flagSet *pflag.FlagSet) string {
	return generateUsageText(flagSet)
}
