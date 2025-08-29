package audio

import "testing"

func TestDefaultConfig(t *testing.T) {
    cfg := DefaultConfig()
    if cfg == nil || !cfg.Enabled {
        t.Fatalf("default audio config should be enabled")
    }
    if cfg.Volume <= 0 || cfg.Volume > 1 {
        t.Fatalf("default volume out of range: %v", cfg.Volume)
    }
    if cfg.Sounds[string(PomodoroComplete)] == "" || cfg.Sounds[string(BreakComplete)] == "" || cfg.Sounds[string(SessionStart)] == "" {
        t.Fatalf("default sounds not populated: %#v", cfg.Sounds)
    }
}

func TestNoOpPlayer(t *testing.T) {
    var p Player = &NoOpPlayer{}
    if p.IsEnabled() {
        t.Fatalf("NoOpPlayer should not be enabled")
    }
    if err := p.Play(PomodoroComplete); err != nil {
        t.Fatalf("NoOpPlayer Play should not error: %v", err)
    }
    if err := p.SetVolume(0.8); err != nil {
        t.Fatalf("NoOpPlayer SetVolume should not error: %v", err)
    }
    if err := p.Close(); err != nil {
        t.Fatalf("NoOpPlayer Close should not error: %v", err)
    }
}

