package volume

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Detects the currently active audio device for playback
func getActiveAudioSink() (string, error) {
	cmd := exec.Command("wpctl", "status")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get PipeWire status: %v", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	activeDevice := "@DEFAULT_AUDIO_SINK@" // Default fallback
	deviceRegex := regexp.MustCompile(`\*\s(\d+).*\b(Audio|Headphone|Bluetooth)\b`)

	for scanner.Scan() {
		line := scanner.Text()
		if matches := deviceRegex.FindStringSubmatch(line); len(matches) > 1 {
			activeDevice = matches[1] // Extracts the device ID
			break
		}
	}

	return activeDevice, nil
}

// Retrieves the volume level and mute status for a given device
func getVolume(device string) (int, bool, error) {
	cmd := exec.Command("wpctl", "get-volume", device)
	output, err := cmd.Output()
	if err != nil {
		return -1, false, fmt.Errorf("failed to get volume for %s: %v", device, err)
	}
	result := strings.TrimSpace(string(output))

	re := regexp.MustCompile(`Volume:\s([\d\.]+)(\s\[MUTED\])?`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return -1, false, fmt.Errorf("unexpected output format")
	}

	volumeFloat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return -1, false, fmt.Errorf("error parsing volume value")
	}

	muted := strings.Contains(result, "[MUTED]") // If "[MUTED]" is present, the sound is muted

	return int(volumeFloat * 100), muted, nil
}
