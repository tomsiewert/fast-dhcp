package cmd

import (
	"io"
	"os"

	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tomsiewert/fast-dhcp/pkg/dhcp"
	"github.com/tomsiewert/fast-dhcp/pkg/model"
	"github.com/tomsiewert/fast-dhcp/pkg/prometheus"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration management tools",
	}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate the config file at the specified path",
		Run: func(cmd *cobra.Command, args []string) {
			generateConfig(cfgFile)
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(generateCmd)
}

func generateConfig(path string) error {
	defaultConfig := &model.Config{
		SentryDSN:   "",
		PProfListen: "",
		DHCP: dhcp.Config{
			Server: dhcp.DHCPServer{
				Hostname:  "",
				Interface: "vmbr0",
				Port:      67,
			},
		},
		Prometheus: prometheus.Config{
			ListenAddress: ":9182",
			ReadTimeout:   10,
		},
	}

	file, _ := json.MarshalIndent(defaultConfig, "", "  ")

	return os.WriteFile(path, file, 0644)
}

func readConfig(path string) *model.Config {
	var configuration model.Config
	configFile, err := os.Open(path)
	if err != nil {
		logrus.Fatal(err)
	}
	defer configFile.Close()

	configContent, _ := io.ReadAll(configFile)
	json.Unmarshal(configContent, &configuration)
	return &configuration
}
