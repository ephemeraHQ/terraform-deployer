package main

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/jessevdk/go-flags"
	"github.com/xmtp-labs/terraform-deployer/pkg/deployer"
	"github.com/xmtp-labs/terraform-deployer/pkg/options"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	var opts options.Options
	if _, err = flags.Parse(&opts); err != nil {
		if err, ok := err.(*flags.Error); !ok || err.Type != flags.ErrHelp {
			log.Fatal("Could not parse options", zap.Error(err))
		}
		return
	}
	err = validateOptions(opts)
	if err != nil {
		log.Fatal("Invalid options", zap.Error(err))
	}

	runTitle, err := getRunTitle(opts)
	if err != nil {
		log.Fatal("Could not get run title", zap.Error(err))
	}

	if opts.DryRun {
		log.Info("Dry run. Not deploying", zap.Any("opts", opts), zap.String("runTitle", runTitle))
		return
	}

	deployer, err := getDeployer(opts.TFToken, opts.Organization, opts.Workspace, log)
	if err != nil {
		log.Fatal("Could not create deployer", zap.Error(err))
	}

	err = deployer.Deploy(opts.VariableName, opts.VariableValue, runTitle)
	if err != nil {
		log.Fatal("Could not deploy", zap.Error(err))
	}
}

func getDeployer(token, organization, workspace string, log *zap.Logger) (*deployer.Deployer, error) {
	tfc, err := tfe.NewClient(&tfe.Config{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	return deployer.NewDeployer(context.Background(), log, tfc, &deployer.Config{
		Organization: organization,
		Workspace:    workspace,
	})
}

func getRunTitle(opts options.Options) (string, error) {
	if opts.RunTitle != "" {
		return opts.RunTitle, nil
	}

	out, err := exec.Command("git", "log", "--oneline", "-n 1").Output()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}

func validateOptions(opts options.Options) error {
	if opts.TFToken == "" {
		return errors.New("Terraform token is required")
	}
	if opts.Workspace == "" {
		return errors.New("Workspace is required")
	}
	if opts.Organization == "" {
		return errors.New("Organization is required")
	}
	if opts.VariableName == "" {
		return errors.New("Variable name is required")
	}
	if opts.VariableValue == "" {
		return errors.New("Variable value is required")
	}
	if opts.VariableValueRequiredPrefix != "" {
		if !strings.HasPrefix(opts.VariableValue, opts.VariableValueRequiredPrefix) {
			return errors.New("Variable value does not start with required prefix")
		}
	}

	return nil
}
