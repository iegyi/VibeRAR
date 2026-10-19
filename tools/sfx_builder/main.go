package main

import (
	"bytes"
	"fmt"
	"os"
)

var MagicMarker = []byte("WINRAR_SFX_PAYLOAD_V1\x00\x00")

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: sfx_builder <stub_exe_path> <zip_archive_path> <output_exe_path> [destination_path] [setup_cmd]")
		return
	}

	stubPath := os.Args[1]
	zipPath := os.Args[2]
	outputPath := os.Args[3]

	destPath := "extracted_sfx_output"
	if len(os.Args) >= 5 {
		destPath = os.Args[4]
	}

	setupCmd := ""
	if len(os.Args) >= 6 {
		setupCmd = os.Args[5]
	}

	stubBytes, err := os.ReadFile(stubPath)
	if err != nil {
		fmt.Printf("Error reading stub: %v\n", err)
		return
	}

	zipBytes, err := os.ReadFile(zipPath)
	if err != nil {
		fmt.Printf("Error reading archive: %v\n", err)
		return
	}

	configScript := fmt.Sprintf("\n;--- VibeRAR SFX Script Configuration ---\nPath=%s\nSetup=%s\nSilent=0\nOverwrite=1\nTitle=VibeRAR self-extracting archive\n", destPath, setupCmd)

	var buffer bytes.Buffer
	buffer.Write(stubBytes)
	buffer.Write(MagicMarker)
	buffer.Write(zipBytes)
	buffer.WriteString(configScript)

	if err := os.WriteFile(outputPath, buffer.Bytes(), 0755); err != nil {
		fmt.Printf("Error writing SFX executable: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated standalone SFX executable: %s (%d bytes)\n", outputPath, buffer.Len())
}
