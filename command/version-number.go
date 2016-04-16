package command

type VersionNumberCommand struct {
	Meta

	Version string
}

func (c *VersionNumberCommand) Run(args []string) int {
	c.Ui.Output(c.Version)
	return 0
}

func (c *VersionNumberCommand) Synopsis() string {
	return ""
}

func (c *VersionNumberCommand) Help() string {
	return ""
}
