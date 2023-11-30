package main

import (
	"github.com/starudream/go-lib/cobra/v2"

	"github.com/starudream/github-webhook/github"
)

var serveCmd = cobra.NewCommand(func(c *cobra.Command) {
	c.Use = "serve"
	c.RunE = func(cmd *cobra.Command, args []string) error {
		return github.Serve()
	}
})

func init() {
	rootCmd.AddCommand(serveCmd)
}
