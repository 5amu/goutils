package shuttles

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

const Placeholder = "^SHT%d^"

type Shuttle struct {
	ctx     context.Context
	args    []string
	injargs []string
	id      int
	output  *ShuttleOutput
}

func NewShuttle(ctx context.Context, args []string, id int) *Shuttle {
	return &Shuttle{ctx: ctx, args: args, id: id, output: &ShuttleOutput{ID: id, Arguments: args}}
}

func (s *Shuttle) InjectArguments(toinject []string) {
	s.output.Injected = toinject
	res := make([]string, len(s.args))
	copy(res, s.args)
	for i, payload := range toinject {
		placeholder := fmt.Sprintf(Placeholder, i)
		for j := range s.args {
			res[j] = strings.ReplaceAll(res[j], placeholder, payload)
		}
	}
	s.injargs = res
}

func (s *Shuttle) Launch() error {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	cmd := exec.CommandContext(ctx, s.injargs[0], s.injargs[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGABRT}
	var buff bytes.Buffer
	cmd.Stdout = &buff

	res := make(chan error, 1)
	go func() {
		res <- cmd.Run()
	}()

	select {
	case <-s.ctx.Done():
		return cmd.Process.Kill()
	case err := <-res:
		s.output.Output = buff.String()
		return err
	}
}
