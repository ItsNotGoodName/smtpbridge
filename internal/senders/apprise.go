package senders

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

//go:embed apprise_script.py
var appriseScript []byte

func AppriseWriteScript(scriptPath string) error {
	return os.WriteFile(scriptPath, appriseScript, 0755)
}

type Apprise struct {
	pythonExecutable string
	scriptPath       string
	urls             []string
}

func NewApprise(pythonExecutable, scriptPath string, urls []string) Apprise {
	return Apprise{
		pythonExecutable: pythonExecutable,
		scriptPath:       scriptPath,
		urls:             urls,
	}
}

type apprisePayload struct {
	URLs        []string                   `json:"urls"`
	Title       string                     `json:"title"`
	Body        string                     `json:"body"`
	Attachments []apprisePayloadAttachment `json:"attachments"`
}

type apprisePayloadAttachment struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

func (a Apprise) Send(ctx context.Context, env models.Envelope, tr Transformer) error {
	payload := apprisePayload{
		URLs:        []string{},
		Title:       "",
		Body:        "",
		Attachments: []apprisePayloadAttachment{},
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

	// Apprise's body cannot be empty
	if payload.Body == "" && payload.Title != "" {
		payload.Body = payload.Title
		payload.Title = ""
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
		payload.Attachments = append(payload.Attachments, apprisePayloadAttachment{
			Path: path,
			Name: files[i].Name,
		})
	}

	// URLs
	payload.URLs = a.urls

	rd, wt := io.Pipe()
	defer rd.Close()

	go func() {
		json.NewEncoder(wt).Encode(payload)
		wt.Close()
	}()

	cmd := exec.CommandContext(ctx, a.pythonExecutable, a.scriptPath)
	cmd.Stdin = rd

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		stdErr := errBuf.String()
		return fmt.Errorf("%w: %s", err, stdErr)
	}

	return nil
}
