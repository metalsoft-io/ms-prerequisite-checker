package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

var version string
var major, minor int

type application struct {
	wg sync.WaitGroup
}

type runtimeConfiguration struct {
	logLevel string
	cmd      string
	args     map[string]string
	handler  func(context.Context, chan<- string, *application, map[string]string)
}

func main() {
	config := processArguments()

	major, minor, _, _ = parseVersionString(version)

	slog.SetDefault(slog.New(NewLogHandler(os.Stdout, config.logLevel)))

	app := &application{}

	ctx, cancel := context.WithCancel(context.Background())

	// Channel for handler to send end message
	endCh := make(chan string, 3)
	go func() {
		endMessage := <-endCh
		slog.Info(endMessage)
		cancel()
	}()

	// Call handler for command
	config.handler(ctx, endCh, app, config.args)

	// Listen for interruptCh signal
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal or end message from handler
	select {
	case <-ctx.Done():
	case <-interruptCh:
		slog.Info("Stopping all services...")
		cancel()
	}

	// Wait for all goroutines to finish
	app.wg.Wait()
}

func processArguments() (conf runtimeConfiguration) {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [flags] command [arg-name=arg-value...]\n", os.Args[0])
		fmt.Print("\nCommands:")

		for _, cmdDetails := range commands {
			arguments := ""
			argumentsHelp := ""
			for _, argDetails := range cmdDetails.arguments {
				arguments += fmt.Sprintf(" %s=<%s>", argDetails.key, argDetails.key)
				argumentsHelp += fmt.Sprintf("        %s: ", argDetails.key)
				if argDetails.required {
					argumentsHelp += "(required) "
				}
				argumentsHelp += argDetails.description
				if argDetails.defaultValue != "" {
					argumentsHelp += fmt.Sprintf(" [default: %s]", argDetails.defaultValue)
				}
				argumentsHelp += "\n"
			}
			fmt.Printf("\n  %s%s\n      %s\n%s", cmdDetails.key, arguments, cmdDetails.description, argumentsHelp)
		}

		fmt.Print("\nFlags:\n")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("  -%s\n      %v\n", f.Name, f.Usage)
		})
	}

	flag.StringVar(&conf.logLevel, "log-level", "INFO", "Sets log level [default 'INFO']")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	conf.cmd = strings.ToLower(flag.Arg(0))
	cmdIndex := slices.IndexFunc(commands, func(cmd commandDetails) bool { return cmd.key == conf.cmd })
	if cmdIndex == -1 {
		flag.Usage()
		os.Exit(1)
	}

	conf.args = make(map[string]string)
	if len(conf.args) > len(commands[cmdIndex].arguments) {
		fmt.Print("Too many arguments\n\n")
		flag.Usage()
		os.Exit(1)
	}

	for _, cmdArg := range flag.Args()[1:] {
		if !strings.Contains(cmdArg, "=") {
			flag.Usage()
			os.Exit(1)
		}
		argParts := strings.Split(cmdArg, "=")
		argName := strings.ToLower(argParts[0])
		argValue := argParts[1]

		argIndex := slices.IndexFunc(commands[cmdIndex].arguments, func(arg argumentDetails) bool { return arg.key == argName })
		if argIndex == -1 {
			fmt.Printf("Unknown argument: %s\n\n", argName)
			flag.Usage()
			os.Exit(1)
		}

		conf.args[argName] = argValue
	}

	for _, argDetails := range commands[cmdIndex].arguments {
		_, ok := conf.args[argDetails.key]
		if !ok && argDetails.required {
			fmt.Printf("Missing required argument: %s\n\n", argDetails.key)
			flag.Usage()
			os.Exit(1)
		}
		if !ok && argDetails.defaultValue != "" {
			conf.args[argDetails.key] = argDetails.defaultValue
		}
	}

	conf.handler = commands[cmdIndex].handler

	return conf
}

func parseVersionString(version string) (int, int, int, error) {
	// Regular expression for version parsing
	re := regexp.MustCompile(`^[A-z]*(\d+)\.(\d+)\.?(\d+)?`)
	if !re.MatchString(version) {
		return 0, 0, 0, errors.New("invalid version string")
	}

	// Split the version string into parts
	versionParts := re.FindStringSubmatch(version)

	major, _ := strconv.Atoi(versionParts[1])

	minor := 0
	if len(versionParts) >= 3 {
		minor, _ = strconv.Atoi(versionParts[2])
	}

	patch := 0
	if len(versionParts) >= 4 {
		patch, _ = strconv.Atoi(versionParts[3])
	}

	return major, minor, patch, nil
}
