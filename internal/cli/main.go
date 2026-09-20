package cli

import (
	"fmt"

	"adolf-codler/buffaring/internal/transfer"
)

// CLI is the top-level kong CLI structure.
type CLI struct {
	Send SendCmd `cmd:"" help:"Send a buffer to a discovered receiver."`
	Recv RecvCmd `cmd:"" help:"Receive a buffer from a discovered sender."`
}

// SendCmd holds the arguments for the send sub-command.
type SendCmd struct {
	Buff string `arg:"" help:"The buffer string to send."`
}

// Run executes the send sub-command.
func (s *SendCmd) Run() error {
	return transfer.SendBuff(s.Buff)
}

// RecvCmd is the receive sub-command (no extra args).
type RecvCmd struct{}

// Run executes the receive sub-command.
// The received buffer is printed to stdout (no prefix) so it can be piped:
//
//	buffaring recv | pbcopy     (macOS)
//	buffaring recv | xclip      (Linux)
func (r *RecvCmd) Run() error {
	data, err := transfer.RecvBuff()
	if err != nil {
		return err
	}
	// Raw output to stdout — pipeable to pbcopy / xclip / xsel
	fmt.Print(data)
	return nil
}
