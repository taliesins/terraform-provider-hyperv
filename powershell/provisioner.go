package powershell

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/masterzen/winrm"
	"github.com/segmentio/ksuid"
)

func TimeOrderedUUID() string {
	id := ksuid.New()
	return id.String()
}

func winPath(path string) string {
	if len(path) == 0 {
		return path
	}

	if strings.Contains(path, " ") {
		path = fmt.Sprintf("'%s'", strings.Trim(path, "'\""))
	}

	return strings.Replace(path, "/", "\\", -1)
}

func doCopy(client *winrm.Client, maxChunks int, in io.Reader, toPath string) (remoteAbsolutePath string, err error) {
	tempFile := fmt.Sprintf("terraform-%s", TimeOrderedUUID())
	tempPath := fmt.Sprintf(`%s\%s`, `$env:TEMP`, tempFile)
	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Resolving remote temp path of [%s]", tempPath)
	}
	tempPath, err = ResolvePath(client, tempPath)
	if err != nil {
		return "", err
	}
	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Remote temp path resolved to [%s]", tempPath)
	}

	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Resolving remote to path of [%s]", toPath)
	}
	toPath, err = ResolvePath(client, toPath)
	if err != nil {
		return "", err
	}
	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Remote to path resolved to [%s]", toPath)
	}

	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Uploading file to %s", tempPath)
	}
	err = uploadContent(client, maxChunks, in, tempPath)
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s: %v", tempPath, err)
	}

	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Moving file from %s to %s", tempPath, toPath)
	}
	remoteAbsolutePath, err = restoreContent(client, tempPath, toPath)
	if err != nil {
		return "", fmt.Errorf("error restoring file from %s to %s: %v", tempPath, toPath, err)
	}

	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Removing temporary file %s", tempPath)
	}
	err = cleanupContent(client, tempPath)
	if err != nil {
		return "", fmt.Errorf("error removing temporary file %s: %v", tempPath, err)
	}

	return remoteAbsolutePath, nil
}

func uploadContent(client *winrm.Client, maxChunks int, in io.Reader, toPath string) error {
	var err error
	done := false
	for !done {
		done, err = uploadChunks(client, maxChunks, in, toPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadChunks(client *winrm.Client, maxChunks int, in io.Reader, toPath string) (bool, error) {
	shell, err := client.CreateShell()
	if err != nil {
		return false, fmt.Errorf("couldn't create shell: %v", err)
	}
	defer shell.Close()

	// Upload the file in chunks to get around the Windows command line size limit.
	// Base64 encodes each set of three bytes into four bytes. In addition the output
	// is padded to always be a multiple of four.
	//
	//   ceil(n / 3) * 4 = m1 - m2
	//
	//   where:
	//     n  = bytes
	//     m1 = max (8192 character command limit.)
	//     m2 = len(toPath)

	chunkSize := ((8000 - len(toPath)) / 4) * 3
	chunk := make([]byte, chunkSize)

	if maxChunks == 0 {
		maxChunks = 1
	}

	for i := 0; i < maxChunks; i++ {
		n, err := in.Read(chunk)

		if err != nil && err != io.EOF {
			return false, err
		}
		if n == 0 {
			return true, nil
		}

		content := base64.StdEncoding.EncodeToString(chunk[:n])
		if err = appendContent(shell, toPath, content); err != nil {
			return false, err
		}
	}

	return false, nil
}

func restoreContent(client *winrm.Client, fromPath, toPath string) (string, error) {
	shell, err := client.CreateShell()
	if err != nil {
		return "", err
	}
	defer shell.Close()

	var convertBase64FileToTextFileTemplateRendered bytes.Buffer
	err = convertBase64FileToTextFileTemplate.Execute(&convertBase64FileToTextFileTemplateRendered, convertBase64FileToTextFileTemplateOptions{
		Base64FilePath: fromPath,
		FilePath:       toPath,
	})

	if err != nil {
		return "", err
	}

	script := convertBase64FileToTextFileTemplateRendered.String()

	var executePowershellFromCommandLineTemplateRendered bytes.Buffer
	err = executePowershellFromCommandLineTemplate.Execute(&executePowershellFromCommandLineTemplateRendered, executePowershellFromCommandLineTemplateOptions{
		Powershell: script,
	})

	if err != nil {
		return "", err
	}

	script = executePowershellFromCommandLineTemplateRendered.String()

	commandExitCode, stdOutPut, errorOutPut, err := shellExecute(shell, script)

	if err != nil {
		return "", err
	}

	if commandExitCode != 0 {
		return "", fmt.Errorf("restore operation returned code=%d\nstderr:\n%s\nstdOut:\n%s", commandExitCode, errorOutPut, stdOutPut)
	}

	if len(errorOutPut) > 0 {
		return "", fmt.Errorf("restore operation returned \nstderr:\n%s\nstdOut:\n%s", errorOutPut, stdOutPut)
	}

	return stdOutPut, nil
}

func ResolvePath(client *winrm.Client, filePath string) (string, error) {
	shell, err := client.CreateShell()
	if err != nil {
		return "", err
	}
	defer shell.Close()

	var resolvePathTemplateRendered bytes.Buffer
	err = resolvePathTemplate.Execute(&resolvePathTemplateRendered, resolvePathTemplateOptions{
		FilePath: filePath,
	})

	if err != nil {
		return "", err
	}

	script := resolvePathTemplateRendered.String()

	var executePowershellFromCommandLineTemplateRendered bytes.Buffer
	err = executePowershellFromCommandLineTemplate.Execute(&executePowershellFromCommandLineTemplateRendered, executePowershellFromCommandLineTemplateOptions{
		Powershell: script,
	})

	if err != nil {
		return "", err
	}

	script = executePowershellFromCommandLineTemplateRendered.String()

	commandExitCode, stdOutPut, errorOutPut, err := shellExecute(shell, script)

	if err != nil {
		return "", err
	}

	if commandExitCode != 0 {
		return "", fmt.Errorf("resolve path operation returned code=%d\nstderr:\n%s\nstdOut:\n%s", commandExitCode, errorOutPut, stdOutPut)
	}

	if len(errorOutPut) > 0 {
		return "", fmt.Errorf("resolve path operation returned \nstderr:\n%s\nstdOut:\n%s", errorOutPut, stdOutPut)
	}

	return stdOutPut, nil
}

func cleanupContent(client *winrm.Client, filePath string) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}
	defer shell.Close()

	var deleteFileTemplateRendered bytes.Buffer
	err = deleteFileTemplate.Execute(&deleteFileTemplateRendered, deleteFileTemplateOptions{
		FilePath: filePath,
	})

	if err != nil {
		return err
	}

	script := deleteFileTemplateRendered.String()

	var executePowershellFromCommandLineTemplateRendered bytes.Buffer
	err = executePowershellFromCommandLineTemplate.Execute(&executePowershellFromCommandLineTemplateRendered, executePowershellFromCommandLineTemplateOptions{
		Powershell: script,
	})

	if err != nil {
		return err
	}

	script = executePowershellFromCommandLineTemplateRendered.String()

	commandExitCode, stdOutPut, errorOutPut, err := shellExecute(shell, script)

	if err != nil {
		return err
	}

	if commandExitCode != 0 {
		return fmt.Errorf("cleanup operation returned code=%d\nstderr:\n%s\nstdOut:\n%s", commandExitCode, errorOutPut, stdOutPut)
	}

	if len(errorOutPut) > 0 {
		return fmt.Errorf("cleanup operation returned \nstderr:\n%s\nstdOut:\n%s", errorOutPut, stdOutPut)
	}

	return nil
}

func appendContent(shell *winrm.Shell, filePath, content string) error {
	var appendFileTemplateRendered bytes.Buffer
	err := appendFileTemplate.Execute(&appendFileTemplateRendered, appendFileTemplateOptions{
		FilePath: filePath,
		Content:  content,
	})

	if err != nil {
		return err
	}

	script := appendFileTemplateRendered.String()

	commandExitCode, stdOutPut, errorOutPut, err := shellExecute(shell, script)

	if err != nil {
		return err
	}

	if commandExitCode != 0 {
		return fmt.Errorf("upload operation returned code=%d\nstderr:\n%s\nstdOut:\n%s", commandExitCode, errorOutPut, stdOutPut)
	}

	if len(errorOutPut) > 0 {
		return fmt.Errorf("upload operation returned \nstderr:\n%s\nstdOut:\n%s", errorOutPut, stdOutPut)
	}

	return nil
}

func shellExecute(shell *winrm.Shell, command string, arguments ...string) (int, string, string, error) {
	stdOutBytes := new(bytes.Buffer)
	stdErrBytes := new(bytes.Buffer)

	stdOutFunc := func(bytesStdOutWriter io.Writer, osStdOutWriter io.Writer, commandStdOut io.Reader) {
		stdOutReader := io.TeeReader(commandStdOut, bytesStdOutWriter)
		_, _ = io.Copy(osStdOutWriter, stdOutReader)
	}

	stdErrFunc := func(bytesStdErrWriter io.Writer, osStdErrWriter io.Writer, commandErrOut io.Reader) {
		stdErrReader := io.TeeReader(commandErrOut, bytesStdErrWriter)
		_, _ = io.Copy(osStdErrWriter, stdErrReader)
	}

	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Shell execute: %s %s", command, arguments)
	}

	cmd, err := shell.Execute(command, arguments...)

	if err != nil {
		return 0, "", "", err
	}

	var closed = false

	defer func() {
		if !closed {
			cmd.Close()
		}
	}()

	go stdOutFunc(stdOutBytes, os.Stdout, cmd.Stdout)
	go stdErrFunc(stdErrBytes, os.Stderr, cmd.Stderr)

	cmd.Wait()
	exitCode := cmd.ExitCode()

	err = cmd.Close()
	closed = true
	if err != nil {
		return 0, "", "", err
	}

	err = cmd.Stdout.Close()
	if err != nil {
		return 0, "", "", err
	}
	stdOutString := strings.TrimSpace(stdOutBytes.String())

	err = cmd.Stderr.Close()
	if err != nil {
		return 0, "", "", err
	}
	stdErrString := strings.TrimSpace(stdErrBytes.String())

	if os.Getenv("WINRMCP_DEBUG") != "" {
		log.Printf("[DEBUG] Shell execute result: exitCode=%d stdOut=%s stdErr=%s", exitCode, stdOutString, stdErrString)
	}

	return exitCode, stdOutString, stdErrString, nil
}

func uploadScript(client *winrm.Client, fileName string, command string) (remoteAbsolutePath string, err error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), fileName)
	writer := bufio.NewWriter(tmpFile)
	if _, err := writer.WriteString(command); err != nil {
		return "", fmt.Errorf("error preparing shell script: %s", err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("error preparing shell script: %s", err)
	}
	tmpFile.Close()
	f, err := os.Open(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("error opening temporary shell script: %s", err)
	}
	defer f.Close()
	defer os.Remove(tmpFile.Name())

	remotePath := fmt.Sprintf(`%s\%s`, `$env:TEMP`, fileName)

	log.Printf("[DEBUG] Uploading shell wrapper for command from [%s] to [%s] ", tmpFile.Name(), remotePath)

	remoteAbsolutePath, err = doCopy(client, 15, f, winPath(remotePath))
	if err != nil {
		return "", fmt.Errorf("error uploading shell script: %s", err)
	}

	return remoteAbsolutePath, nil
}

func createCommand(vars string, remotePath string) (commandText string, err error) {
	var executeCommandTemplateRendered bytes.Buffer

	err = executeCommandTemplate.Execute(&executeCommandTemplateRendered, executeCommandTemplateOptions{
		Vars: vars,
		Path: remotePath,
	})

	if err != nil {
		fmt.Printf("Error creating command template: %s", err)
		return "", err
	}

	commandText = executeCommandTemplateRendered.String()

	return commandText, err
}

func createElevatedCommand(client *winrm.Client, elevatedUser string, elevatedPassword string, vars string, remotePath string) (commandText string, elevatedRemotePath string, err error) {
	elevatedRemotePath, err = generateElevatedRunner(client, elevatedUser, elevatedPassword, remotePath)
	if err != nil {
		return "", "", fmt.Errorf("error generating elevated runner: %s", err)
	}

	commandText, err = createCommand(vars, elevatedRemotePath)

	return commandText, elevatedRemotePath, err
}

func generateElevatedRunner(client *winrm.Client, elevatedUser string, elevatedPassword string, remotePath string) (elevatedRemotePath string, err error) {
	log.Printf("[DEBUG] Building elevated command wrapper for: %s", remotePath)

	name := fmt.Sprintf("terraform-%s", TimeOrderedUUID())
	fileName := fmt.Sprintf(`elevated-shell-%s.ps1`, name)

	var elevatedCommandTemplateRendered bytes.Buffer
	err = elevatedCommandTemplate.Execute(&elevatedCommandTemplateRendered, elevatedCommandTemplateOptions{
		User:                   elevatedUser,
		Password:               elevatedPassword,
		TaskDescription:        "Terraform elevated task",
		TaskName:               name,
		TaskExecutionTimeLimit: "PT2H",
		ScriptPath:             remotePath,
	})

	if err != nil {
		fmt.Printf("Error creating elevated command template: %s", err)
		return "", err
	}

	elevatedCommand := elevatedCommandTemplateRendered.String()

	elevatedRemotePath, err = uploadScript(client, fileName, elevatedCommand)
	if err != nil {
		return "", err
	}

	return elevatedRemotePath, nil
}

// Run powershell
func RunPowershell(client *winrm.Client, elevatedUser string, elevatedPassword string, vars string, commandText string) (exitStatus int, stdout string, stderr string, err error) {
	name := fmt.Sprintf("terraform-%s", TimeOrderedUUID())
	fileName := fmt.Sprintf(`shell-%s.ps1`, name)

	path, err := uploadScript(client, fileName, commandText)
	if err != nil {
		return 0, "", "", err
	}

	var command string

	if elevatedUser == "" {
		command, err = createCommand(vars, path)
	} else {
		command, path, err = createElevatedCommand(client, elevatedUser, elevatedPassword, vars, path)
	}

	if err != nil {
		return 0, "", "", err
	}

	var executePowershellFromCommandLineTemplateRendered bytes.Buffer
	err = executePowershellFromCommandLineTemplate.Execute(&executePowershellFromCommandLineTemplateRendered, executePowershellFromCommandLineTemplateOptions{
		Powershell: command,
	})

	if err != nil {
		return 0, "", "", err
	}

	command = executePowershellFromCommandLineTemplateRendered.String()

	shell, err := client.CreateShell()
	if err != nil {
		return 0, "", "", err
	}
	defer shell.Close()

	commandExitCode, stdOutPut, errorOutPut, err := shellExecute(shell, command)

	if err != nil {
		return 0, "", "", err
	}

	if commandExitCode != 0 {
		return 0, "", "", fmt.Errorf("run command operation returned code=%d\nstderr:\n%s\nstdOut:\n%s", commandExitCode, errorOutPut, stdOutPut)
	}

	if len(errorOutPut) > 0 {
		return 0, "", "", fmt.Errorf("run command operation returned \nstderr:\n%s\nstdOut:\n%s", errorOutPut, stdOutPut)
	}

	err = cleanupContent(client, path)
	if err != nil {
		return 0, "", "", fmt.Errorf("error removing temporary file %s: %v", path, err)
	}

	return commandExitCode, stdOutPut, errorOutPut, nil
}
