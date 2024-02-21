package options

import "time"

type Options struct {
	TFToken                     string        `short:"t" long:"terraform-token" env:"TERRAFORM_TOKEN" description:"Terraform API token"`
	Workspace                   string        `short:"w" long:"workspace" env:"TERRAFORM_WORKSPACE" description:"Terraform Workspace"`
	Organization                string        `short:"o" long:"org" env:"TERRAFORM_ORGANIZATION" description:"Terraform Organization"`
	VariableName                string        `short:"n" long:"variable-name" env:"TERRAFORM_VARIABLE_NAME" description:"Terraform Variable Name"`
	VariableValue               string        `short:"v" long:"variable-value" env:"TERRAFORM_VARIABLE_VALUE" description:"Value to set for the Terraform Variable named in VariableName"`
	RunTitle                    string        `long:"run-title" env:"TERRAFORM_RUN_TITLE" description:"Title for the Terraform Run. Defaults to latest commit message if unset."`
	DryRun                      bool          `long:"dry-run" description:"Do not actually run the Terraform Run. Useful for testing."`
	VariableValueRequiredPrefix string        `long:"variable-value-required-prefix" env:"VARIABLE_VALUE_REQUIRED_PREFIX" description:"If set, the VariableValue must start with this prefix"`
	Timeout                     time.Duration `long:"timeout" description:"Run timeout"`
}
