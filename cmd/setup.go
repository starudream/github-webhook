package main

import (
	"github.com/starudream/go-lib/cobra/v2"

	"github.com/starudream/github-webhook/tea"
)

var setupCmd = cobra.NewCommand(func(c *cobra.Command) {
	c.Use = "setup"
	c.RunE = func(cmd *cobra.Command, args []string) error {
		m, err := tea.Run(tea.NewEvent())
		if err != nil {
			return tea.Return(err)
		}
		m.View()
		return nil
	}
})

func init() {
	rootCmd.AddCommand(setupCmd)
}
