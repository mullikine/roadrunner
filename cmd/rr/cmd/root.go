// Copyright (c) 2018 SpiralScout
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spiral/roadrunner/cmd/rr/utils"
	"github.com/spiral/roadrunner/service"
	"os"
)

// Service bus for all the commands.
var (
	cfgFile string
	verbose bool

	// Logger - shared logger.
	Logger = logrus.New()

	// Container - shared service bus.
	Container = service.NewContainer(Logger)

	// CLI is application endpoint.
	CLI = &cobra.Command{
		Use:           "rr",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         utils.Sprintf("<green>RoadRunner, PHP Application Server.</reset>"),
	}
)

// Execute adds all child commands to the CLI command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the CLI.
func Execute() {
	if err := CLI.Execute(); err != nil {
		utils.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
		os.Exit(1)
	}
}

func init() {
	CLI.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	CLI.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .rr.yaml)")

	cobra.OnInitialize(func() {
		if verbose {
			Logger.SetLevel(logrus.DebugLevel)
		}

		if cfg := initConfig(cfgFile, []string{"."}, ".rr"); cfg != nil {
			if err := Container.Configure(cfg); err != nil {
				utils.Printf("<red+hb>Error:</reset> <red>%s</reset>\n", err)
				os.Exit(1)
			}
		}
	})
}

func initConfig(cfgFile string, path []string, name string) service.Config {
	cfg := viper.New()

	if cfgFile != "" {
		// Use cfg file from the flag.
		cfg.SetConfigFile(cfgFile)
	} else {
		// automatic location
		for _, p := range path {
			cfg.AddConfigPath(p)
		}

		cfg.SetConfigName(name)
	}

	// read in environment variables that match
	cfg.AutomaticEnv()

	// If a cfg file is found, read it in.
	if err := cfg.ReadInConfig(); err != nil {
		Logger.Warnf("config: %s", err)
		return nil
	}

	return &utils.ViperWrapper{Viper: cfg}
}
