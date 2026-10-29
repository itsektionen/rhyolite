package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

func execute(cmdConfig CommandConfig, request ExecuteRequest) ([]byte, error) {
	maxTime := time.Second * 120
	ctx, cancel := context.WithTimeout(context.Background(), maxTime)
	defer cancel()
	command, err := cmdConfig.getBashCommand(request)
	if err != nil {
		err = fmt.Errorf("Cannot run program: %w", err)
		return []byte(err.Error()), err
	}
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	if !cmdConfig.DontOverrideHome {
		// We override home to allow each execution to have its own docker login
		cmd.Env = append(cmd.Environ(), "HOME="+filepath.Join(config.RockDir, request.rock))
	}
	cmd.Dir = filepath.Join(config.RockDir, request.rock)
	output := bytes.Buffer{}
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("Program exited with error:\n%s\n----\n[Internal error: %w]\n", output.String(), err)
		return []byte(err.Error()), err
	}
	return []byte(output.String()), nil
}
