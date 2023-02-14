package volume

import (
  "fmt"
  "io/ioutil"
  "math"
  "strconv"
  "strings"
  "time"

  "barista.run/bar"
  "barista.run/base/value"
  "barista.run/outputs"
  "barista.run/timing"
)

SINK_NAME="@DEFAULT_AUDIO_SINK@"
SOURCE_NAME="@DEFAULT_AUDIO_SOURCE@"

type Module struct {
  outputFunc value.Value // of func(Volume) bar.Output
  provider   Provider
}

// Output configures a module to display the output of a user-defined
// function.
func (m *Module) Output(outputFunc func(Volume) bar.Output) *Module {
  m.outputFunc.Set(outputFunc)
  return m
}

// RateLimiter throttles volume updates to once every ~20ms to avoid unexpected behaviour.
var RateLimiter = rate.NewLimiter(rate.Every(20*time.Millisecond), 1)

func (m *Module) Stream(s bar.Sink) {
  var vol value.ErrorValue

  v, err := vol.Get()
  nextV, done := vol.Subscribe()
  defer done()
  go m.provider.Worker(&vol)

  outputFunc := m.outputFunc.Get().(func(Volume) bar.Output)
  nextOutputFunc, done := m.outputFunc.Subscribe()
  defer done()

  for {
    if s.Error(err) {
      return
    }
    if volume, ok := v.(Volume); ok {
      volume.update = func(v Volume) { vol.Set(v) }
      s.Output(outputs.Group(outputFunc(volume)).
        OnClick(defaultClickHandler(volume)))
    }
    select {
    case <-nextV:
      v, err = vol.Get()
    case <-nextOutputFunc:
      outputFunc = m.outputFunc.Get().(func(Volume) bar.Output)
    }
  }
}

// New creates a new module with the given backing implementation.
func New(provider Provider) *Module {
  m := &Module{provider: provider}
  l.Register(m, "outputFunc", "impl")
  // Default output is just the volume %, "MUT" when muted.
  m.Output(func(v Volume) bar.Output {
    if v.Mute {
      return outputs.Text("MUT")
    }
    return outputs.Textf("%d%%", v.Pct())
  })
  return m
}
