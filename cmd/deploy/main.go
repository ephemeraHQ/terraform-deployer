package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

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

	runTitle, err := getRunTitle(opts, log)
	if err != nil {
		log.Fatal("Could not get run title", zap.Error(err))
	}

	if opts.DryRun {
		log.Info("Dry run. Not deploying", zap.Any("opts", opts), zap.String("runTitle", runTitle))
		return
	}

	vars := make(map[string]string)
	for i := range opts.VariableName {
		vars[opts.VariableName[i]] = opts.VariableValue[i]
	}

	deployer, err := getDeployer(opts.TFToken, opts.Organization, opts.Workspace, log, opts.Timeout)
	if err != nil {
		log.Fatal("Could not create deployer", zap.Error(err))
	}

	err = deployer.Deploy(vars, runTitle)
	if err != nil {
		log.Fatal("Could not deploy", zap.Error(err))
	}
}

func getDeployer(
	token, organization, workspace string,
	log *zap.Logger,
	timeout time.Duration,
) (*deployer.Deployer, error) {
	tfc, err := tfe.NewClient(&tfe.Config{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return deployer.NewDeployer(context.Background(), log, tfc, &deployer.Config{
		Organization: organization,
		Workspace:    workspace,
		WaitTimeout:  timeout,
	})
}

func getRunTitle(opts options.Options, log *zap.Logger) (string, error) {
	if opts.RunTitle != "" {
		return opts.RunTitle, nil
	}

	out, err := exec.Command("git", "log", "--oneline", "-n 1").Output()
	if err != nil {
		log.Error("git log failed", zap.String("output", string(out)))
		return "", err
	}

	return strings.Trim(string(out), "\n"), nil
}

func validateOptions(opts options.Options) error {
	if opts.TFToken == "" {
		return errors.New("terraform token is required")
	}
	if opts.Workspace == "" {
		return errors.New("workspace is required")
	}
	if opts.Organization == "" {
		return errors.New("organization is required")
	}

	if len(opts.VariableName) == 0 {
		return errors.New("at least one variable name is required")
	}
	if len(opts.VariableValue) == 0 {
		return errors.New("at least one variable value is required")
	}

	if len(opts.VariableValue) != len(opts.VariableName) {
		return errors.New("variable name and value must be the same length")
	}

	if opts.VariableValueRequiredPrefix != "" {
		for i, val := range opts.VariableValue {
			if !strings.HasPrefix(val, opts.VariableValueRequiredPrefix) {
				return fmt.Errorf(
					"variable %s:%s does not start with required prefix",
					opts.VariableName[i],
					val,
				)
			}
		}
	}

	return nil
}
