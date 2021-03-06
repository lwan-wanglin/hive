package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	genericapiserver "k8s.io/apiserver/pkg/server"
	utilflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"

	"github.com/openshift/hive/pkg/certs"
	hiveapiserver "github.com/openshift/hive/v1alpha1apiserver/pkg/cmd/hive-apiserver"
	"github.com/openshift/hive/pkg/version"
)

func main() {
	stopCh := genericapiserver.SetupSignalHandler()

	rand.Seed(time.Now().UTC().UnixNano())

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	logs.InitLogs()
	defer logs.FlushLogs()

	log.Infof("Version: %s @ %s", version.String, version.Commit)

	log.SetLevel(log.InfoLevel)

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	// TODO: Would be better to get this from the --tls-cert-file argument.
	const certsDir = "/apiserver.local.config/certificates"
	certs.TerminateOnCertChanges(certsDir)

	command := NewHiveAPIServerCommand(stopCh)
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func NewHiveAPIServerCommand(stopCh <-chan struct{}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hive-apiserver",
		Short: "Command for the Hive v1alpha1 API Server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}
	start := hiveapiserver.NewHiveAPIServerCommand("start", os.Stdout, os.Stderr, stopCh)
	cmd.AddCommand(start)

	return cmd
}
