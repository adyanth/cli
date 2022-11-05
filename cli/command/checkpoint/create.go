package checkpoint

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/completion"
	"github.com/docker/docker/api/types"
	"github.com/spf13/cobra"
)

type createOptions struct {
	container      string
	checkpoint     string
	checkpointDir  string
	leaveRunning   bool
	tcpEstablished bool
	unixSockets    bool
	terminal       bool
	fileLocks      bool
}

func newCreateCommand(dockerCli command.Cli) *cobra.Command {
	var opts createOptions

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] CONTAINER CHECKPOINT",
		Short: "Create a checkpoint from a running container",
		Args:  cli.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.container = args[0]
			opts.checkpoint = args[1]
			return runCreate(dockerCli, opts)
		},
		ValidArgsFunction: completion.NoComplete,
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.leaveRunning, "leave-running", false, "Leave the container running after checkpoint")
	flags.BoolVar(&opts.tcpEstablished, "tcp-established", false, "Dump open TCP sockets in the checkpoint")
	flags.BoolVar(&opts.unixSockets, "ext-unix-sk", false, "Dump Unix sockets in the checkpoint")
	flags.BoolVar(&opts.terminal, "shell-job", false, "Dump shell jobs in the checkpoint")
	flags.BoolVar(&opts.fileLocks, "file-locks", false, "Dump file locks in the checkpoint")
	flags.StringVarP(&opts.checkpointDir, "checkpoint-dir", "", "", "Use a custom checkpoint storage directory")

	return cmd
}

func runCreate(dockerCli command.Cli, opts createOptions) error {
	client := dockerCli.Client()

	checkpointOpts := types.CheckpointCreateOptions{
		CheckpointID:  opts.checkpoint,
		CheckpointDir: opts.checkpointDir,
		Exit:          !opts.leaveRunning,
		OpenTcp:       opts.tcpEstablished,
		UnixSockets:   opts.unixSockets,
		Terminal:      opts.terminal,
		FileLocks:     opts.fileLocks,
	}

	err := client.CheckpointCreate(context.Background(), opts.container, checkpointOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", opts.checkpoint)
	return nil
}
