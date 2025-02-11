package volume

import (
	"time"

	"github.com/barista-run/barista/bar"
	"github.com/barista-run/barista/base/value"
	"github.com/barista-run/barista/outputs"
	"github.com/barista-run/barista/timing"
)

// Module represents a volume bar module
type Module struct {
	outputFunc value.Value // of func(int, bool, bool) bar.Output
	scheduler  *timing.Scheduler
	device     string // The device to monitor (e.g., @DEFAULT_AUDIO_SINK@ or specific device ID)
	isMic      bool   // True for microphone, false for speaker/headphones
}

// Initializes a new volume module for a specific device
func New(device string, isMic bool) *Module {
	m := &Module{
		scheduler: timing.NewScheduler(),
		device:    device,
		isMic:     isMic,
	}

	m.RefreshInterval(5 * time.Second)

	m.outputFunc.Set(func(vol int, muted bool, isMic bool) bar.Output {
		return outputs.Textf("%d%%", vol)
	})

	return m
}

func (m *Module) RefreshInterval(interval time.Duration) *Module {
	m.scheduler.Every(interval)
	return m
}

// Allows custom output formatting
func (m *Module) Output(outputFunc func(int, bool, bool) bar.Output) *Module {
	m.outputFunc.Set(outputFunc)
	return m
}

// Streams real-time volume updates
func (m *Module) Stream(s bar.Sink) {
	outputFunc := m.outputFunc.Get().(func(int, bool, bool) bar.Output)

	// Detect the active audio device if using default playback
	if !m.isMic {
		activeDevice, err := getActiveAudioSink()
		if err != nil {
			s.Error(err)
			return
		}
		m.device = activeDevice
	}

	// Initial volume fetch
	volume, muted, err := getVolume(m.device)
	if err != nil {
		s.Error(err)
		return
	}
	s.Output(outputFunc(volume, muted, m.isMic))

	for {
		volume, muted, err = getVolume(m.device)
		if err != nil {
			s.Error(err)
			return
		}
		s.Output(outputFunc(volume, muted, m.isMic))

		<-m.scheduler.C
	}
}
