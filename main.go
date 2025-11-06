package main

import (
	"vault-cli/cmd"
	"vault-cli/internal/tui"
)

func main() {
	tui.ShowVaultBanner()
	cmd.Execute()
}
