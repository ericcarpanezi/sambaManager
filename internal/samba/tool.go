package samba

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Tool struct {
	binary string
}

func NewTool() *Tool {
	return &Tool{binary: "samba-tool"}
}

func (t *Tool) Run(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, t.binary, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("samba-tool %s failed: %w (%s)", strings.Join(args, " "), err, string(out))
	}
	return string(out), nil
}
