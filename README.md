# Barista module for volume

## 💡 How to Use

Modify your Barista config to include two modules:

```go
// Speaker/Headphones volume (auto-detects active device)
speaker := volume.New("@DEFAULT_AUDIO_SINK@", false)

// Microphone volume
microphone := volume.New("@DEFAULT_AUDIO_SOURCE@", true)

barista.Add(speaker)
barista.Add(microphone)
```
