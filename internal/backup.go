package internal

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tokuhirom/blog4/internal/sobs"
)

func StartBackup(config *Config, s3client *sobs.SobsClient) {
	slog.Info("Start taking backup")
	time.AfterFunc(1*time.Second, func() {
		takeBackup(config, s3client)
	})
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		takeBackup(config, s3client)
	}
}

func takeBackup(config *Config, s3client *sobs.SobsClient) {
	slog.Info("takeBackup")

	// Generate dump file path
	date := time.Now()
	dumpFileName := fmt.Sprintf("/tmp/blog3-backup-%s.sql", date.Format("2006-01-02T15-04-05"))
	encryptedFileName := dumpFileName + ".enc"
	slog.Info("mariadb-dump file name", slog.String("filename", dumpFileName))

	// Build mariadb-dump command with conditional --skip-ssl flag
	sslFlag := ""
	if config.LocalDev {
		sslFlag = "--skip-ssl "
	}

	// Execute mariadb-dump command
	err := execCommand(fmt.Sprintf(
		"mariadb-dump %s--host=%s --port=%d --user=%s --password=%s %s > %s",
		sslFlag,
		config.DBHostname,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		dumpFileName,
	))
	if err != nil {
		slog.Error("Error taking database dump", slog.Any("error", err))
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			slog.Error("Error removing dump file", slog.String("file", name), slog.Any("error", err))
		}
	}(dumpFileName)

	// Encrypt the dump file
	err = execCommand(fmt.Sprintf(
		"openssl enc -aes-256-cbc -salt -in %s -out %s -pass pass:%s",
		dumpFileName,
		encryptedFileName,
		config.BackupEncryptionKey,
	))
	if err != nil {
		slog.Error("Error encrypting dump file", slog.Any("error", err))
		return
	}
	slog.Info("Encrypted file created", slog.String("file", encryptedFileName))
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			slog.Error("Error removing encrypted file", slog.String("file", name), slog.Any("error", err))
		}
	}(encryptedFileName)

	// Upload file to S3
	slog.Info("Uploading file to S3", slog.String("file", encryptedFileName))
	fileContent, err := os.ReadFile(encryptedFileName)
	if err != nil {
		slog.Error("Error reading encrypted file", slog.Any("error", err))
		return
	}

	err = s3client.PutObjectToBackupBucket(
		context.Background(),
		filepath.Base(encryptedFileName),
		"application/octet-stream",
		int64(len(fileContent)),
		bytes.NewReader(fileContent))
	if err != nil {
		slog.Error("Error uploading file to S3", slog.Any("error", err))
		return
	}
	slog.Info("File uploaded to S3", slog.String("file", encryptedFileName))

	// Delete old backup files (keep only last 7 days)
	err = s3client.DeleteOldBackups(context.Background(), 7)
	if err != nil {
		slog.Error("Error deleting old backups", slog.Any("error", err))
		// Don't return - this is not a critical error
	}

	// Clean up temporary files
	slog.Info("Removing temporary files", slog.String("dump_file", dumpFileName), slog.String("encrypted_file", encryptedFileName))
}

func execCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command execution failed: cmd=%s, output=%s, error=%v",
			command, string(output), err)
	}
	return nil
}
