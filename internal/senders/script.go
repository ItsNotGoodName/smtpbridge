package senders

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

type Script struct {
	scriptPath string
}

func NewScript(scriptPath string) Script {
	return Script{
		scriptPath: scriptPath,
	}
}

type scriptPayload struct {
	Title       string                    `json:"title"`
	Body        string                    `json:"body"`
	Attachments []scriptPayloadAttachment `json:"attachments"`
}

type scriptPayloadAttachment struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (s Script) Send(ctx context.Context, env models.Envelope, tr Transformer) error {
	payload := scriptPayload{
		Title:       "",
		Body:        "",
		Attachments: []scriptPayloadAttachment{},
	}

	var err error

	// Title
	payload.Title, err = tr.Title(ctx, env)
	if err != nil {
		return err
	}

	// Body
	payload.Body, err = tr.Body(ctx, env)
	if err != nil {
		return err
	}

	// Files
	files, err := tr.Files(ctx, env)
	if err != nil {
		return err
	}
	for i := range files {
		path, err := tr.Path(ctx, files[i])
		if err != nil {
			return err
		}
		payload.Attachments = append(payload.Attachments, scriptPayloadAttachment{
			Path: path,
			Name: files[i].Name,
		})
	}

	rd, wt := io.Pipe()
	defer rd.Close()

	go func() {
		json.NewEncoder(wt).Encode(payload)
		wt.Close()
	}()

	cmd := exec.CommandContext(ctx, s.scriptPath)
	cmd.Stdin = rd
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		stdErr := errBuf.String()
		return fmt.Errorf("%w: %s", err, stdErr)
	}

	return nil
}
