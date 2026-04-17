package main

import (
	"fmt"
	"gitee.com/liujit/shop/server/cmd/migrate"
	"gitee.com/liujit/shop/server/cmd/start"
	"gitee.com/liujit/shop/server/internal/version"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:     version.Name,
	Short:   version.Name + ": An NmsKit microservices.",
	Long:    version.Name + `: An NmsKit microservices.`,
	Version: version.Release,
}

func init() {
	fmt.Println(version.Name, "Server Initial...", "version:", version.Release)
	rootCmd.AddCommand(start.CmdStart)
	rootCmd.AddCommand(migrate.CmdStart)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
