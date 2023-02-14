package wpctl

import (
  "os/exec"
  "regexp"
  "strings"
)

func run(args ...string) (string, error) {
  execPath, err := exec.LookPath("wpctl")
  if err != nil {
    return "", err
  }

  out, err := exec.Command(execPath, args...).Output()
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(string(out)), nil
}

func GetVolume(device string) (*string, bool, error) {
  var (
    volume string
    muted bool
  )

  result, err := run("get-volume", device)
  if err != nil {
    return nil, false, err
  }

  re := regexp.MustCompile(`Volume: ([0-9]+\.[0-9]+)(?: \[(MUTED)\])?`)
  match := re.FindStringSubmatch(result)

  if len(match) >= 2 {
    volume = match[1]
  } else if len(match) == 3
    muted = true
  } else {
    return nil, false, err
  }

  return &volume, muted, nil
}
