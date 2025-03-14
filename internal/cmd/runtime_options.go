package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/pflag"
	"gopkg.in/guregu/null.v3"

	"go.k6.io/k6/cmd/state"
	"go.k6.io/k6/internal/lib/summary"
	"go.k6.io/k6/lib"
)

// TODO: move this whole file out of the cmd package? maybe when fixing
// https://github.com/k6io/k6/issues/883, since this code is fairly
// self-contained and easily testable now, without any global dependencies...

var userEnvVarName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func runtimeOptionFlagSet(includeSysEnv bool) *pflag.FlagSet {
	flags := pflag.NewFlagSet("", 0)
	flags.SortFlags = false
	flags.Bool("include-system-env-vars", includeSysEnv, "pass the real system environment variables to the runtime")
	flags.String("compatibility-mode", "extended",
		`JavaScript compiler compatibility mode, "extended" or "base"
base: pure Sobek - Golang JS VM supporting ES6+
extended: base + sets "global" as alias for "globalThis"
`)
	flags.StringP("type", "t", "", "override test type, \"js\" or \"archive\"")
	flags.StringArrayP("env", "e", nil, "add/override environment variable with `VAR=value`")
	flags.Bool("no-thresholds", false, "don't run thresholds")
	flags.Bool("no-summary", false, "don't show the summary at the end of the test")
	flags.String("summary-mode", summary.ModeCompact.String(), "determine the summary mode,"+
		" \"compact\", \"full\" or \"legacy\"")
	flags.String(
		"summary-export",
		"",
		"output the end-of-test summary report to JSON file",
	)
	flags.String("traces-output", "none",
		"set the output for k6 traces, possible values are none,otel[=host:port]")
	return flags
}

func getRuntimeOptions(
	flags *pflag.FlagSet,
	environment map[string]string,
) (lib.RuntimeOptions, error) {
	// TODO: refactor with composable helpers as a part of #883, to reduce copy-paste
	// TODO: get these options out of the JSON config file as well?
	opts, err := populateRuntimeOptionsFromEnv(runtimeOptionsFromFlags(flags), environment)
	if err != nil {
		return opts, err
	}

	// Set/overwrite environment variables with custom user-supplied values
	envVars, err := flags.GetStringArray("env")
	if err != nil {
		return opts, err
	}
	for _, kv := range envVars {
		k, v := state.ParseEnvKeyValue(kv)
		// Allow only alphanumeric ASCII variable names for now
		if !userEnvVarName.MatchString(k) {
			return opts, fmt.Errorf("invalid environment variable name '%s'", k)
		}
		opts.Env[k] = v
	}

	return opts, nil
}

func runtimeOptionsFromFlags(flags *pflag.FlagSet) lib.RuntimeOptions {
	opts := lib.RuntimeOptions{
		TestType:             getNullString(flags, "type"),
		IncludeSystemEnvVars: getNullBool(flags, "include-system-env-vars"),
		CompatibilityMode:    getNullString(flags, "compatibility-mode"),
		NoThresholds:         getNullBool(flags, "no-thresholds"),
		NoSummary:            getNullBool(flags, "no-summary"),
		SummaryMode:          getNullString(flags, "summary-mode"),
		SummaryExport:        getNullString(flags, "summary-export"),
		TracesOutput:         getNullString(flags, "traces-output"),
		Env:                  make(map[string]string),
	}
	return opts
}

func populateRuntimeOptionsFromEnv(opts lib.RuntimeOptions, environment map[string]string) (lib.RuntimeOptions, error) {
	// Only override if not explicitly set via the CLI flag

	if envVar, ok := environment["K6_TYPE"]; !opts.TestType.Valid && ok {
		opts.TestType = null.StringFrom(envVar)
	}

	if envVar, ok := environment["K6_COMPATIBILITY_MODE"]; !opts.CompatibilityMode.Valid && ok {
		opts.CompatibilityMode = null.StringFrom(envVar)
	}

	if _, err := lib.ValidateCompatibilityMode(opts.CompatibilityMode.String); err != nil {
		// some early validation
		return opts, err
	}

	if err := saveBoolFromEnv(environment, "K6_INCLUDE_SYSTEM_ENV_VARS", &opts.IncludeSystemEnvVars); err != nil {
		return opts, err
	}

	if err := saveBoolFromEnv(environment, "K6_NO_THRESHOLDS", &opts.NoThresholds); err != nil {
		return opts, err
	}

	if err := saveBoolFromEnv(environment, "K6_NO_SUMMARY", &opts.NoSummary); err != nil {
		return opts, err
	}

	if envVar, ok := environment["K6_SUMMARY_MODE"]; !opts.SummaryMode.Valid && ok {
		opts.SummaryMode = null.StringFrom(envVar)
	}

	if _, err := summary.ValidateMode(opts.SummaryMode.String); err != nil {
		// some early validation
		return opts, err
	}

	if envVar, ok := environment["K6_SUMMARY_EXPORT"]; !opts.SummaryExport.Valid && ok {
		opts.SummaryExport = null.StringFrom(envVar)
	}

	if envVar, ok := environment["SSLKEYLOGFILE"]; !opts.KeyWriter.Valid && ok {
		opts.KeyWriter = null.StringFrom(envVar)
	}

	if envVar, ok := environment["K6_TRACES_OUTPUT"]; !opts.TracesOutput.Valid && ok {
		opts.TracesOutput = null.StringFrom(envVar)
	}

	// If enabled, gather the actual system environment variables
	if opts.IncludeSystemEnvVars.Bool {
		opts.Env = environment
	}

	return opts, nil
}

func saveBoolFromEnv(env map[string]string, varName string, placeholder *null.Bool) error {
	strValue, ok := env[varName]
	if !ok {
		return nil
	}
	val, err := strconv.ParseBool(strValue)
	if err != nil {
		return fmt.Errorf("env var '%s' is not a valid boolean value: %w", varName, err)
	}
	// Only override if not explicitly set via the CLI flag
	if !placeholder.Valid {
		*placeholder = null.BoolFrom(val)
	}
	return nil
}
