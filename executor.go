package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

func execute(config CommandConfig, request ExecuteRequest) ([]byte, error) {
	maxTime := time.Second * 120
	ctx, cancel := context.WithTimeout(context.Background(), maxTime)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", config.Command)
	cmd.Dir = filepath.Join(RhyoliteDir, request.directory)
	output := bytes.Buffer{}
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("Program exited with error:\n%s\n----\n[Internal error: %w]\n", output.String(), err)
		return []byte(err.Error()), err
	}
	return []byte(output.String()), nil
}
