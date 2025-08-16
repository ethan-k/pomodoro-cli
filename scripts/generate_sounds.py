#!/usr/bin/env python3
"""
Generate copyright-free audio notification sounds for Pomodoro CLI.

This script creates simple, pleasant notification sounds using basic synthesis.
All generated sounds are copyright-free and safe for distribution.
"""

import numpy as np
import wave
import math
import os

def generate_sine_wave(frequency, duration, sample_rate=44100, amplitude=0.5):
    """Generate a sine wave with the given frequency and duration."""
    frames = int(duration * sample_rate)
    arr = []
    for i in range(frames):
        # Generate sine wave
        value = amplitude * math.sin(2 * math.pi * frequency * i / sample_rate)
        # Apply gentle fade in/out to avoid clicks
        fade_frames = int(0.1 * sample_rate)  # 100ms fade
        if i < fade_frames:
            value *= i / fade_frames
        elif i > frames - fade_frames:
            value *= (frames - i) / fade_frames
        arr.append(value)
    return np.array(arr)

def generate_chord(frequencies, duration, sample_rate=44100, amplitude=0.3):
    """Generate a chord by combining multiple frequencies."""
    waves = []
    for freq in frequencies:
        wave = generate_sine_wave(freq, duration, sample_rate, amplitude)
        waves.append(wave)
    
    # Combine waves
    combined = np.sum(waves, axis=0)
    # Normalize to prevent clipping
    combined = combined / len(frequencies)
    return combined

def save_wav(filename, audio_data, sample_rate=44100):
    """Save audio data as WAV file."""
    # Convert to 16-bit integers
    audio_data = np.int16(audio_data * 32767)
    
    with wave.open(filename, 'w') as wav_file:
        wav_file.setnchannels(1)  # Mono
        wav_file.setsampwidth(2)  # 16-bit
        wav_file.setframerate(sample_rate)
        wav_file.writeframes(audio_data.tobytes())

def create_pomodoro_complete():
    """Create a gentle bell-like sound for pomodoro completion."""
    # Bell-like frequencies (based on C major chord with overtones)
    fundamental = 523.25  # C5
    overtones = [fundamental, fundamental * 2, fundamental * 3, fundamental * 5]
    
    # Create bell sound with decay
    duration = 2.5
    bell_sound = generate_chord([fundamental], duration, amplitude=0.4)
    
    # Add subtle overtones with quick decay
    for i, freq in enumerate(overtones[1:], 1):
        overtone_duration = duration / (i + 1)
        overtone = generate_sine_wave(freq, overtone_duration, amplitude=0.1 / i)
        # Pad with zeros to match main duration
        if len(overtone) < len(bell_sound):
            overtone = np.pad(overtone, (0, len(bell_sound) - len(overtone)))
        bell_sound[:len(overtone)] += overtone[:len(bell_sound)]
    
    return bell_sound

def create_break_complete():
    """Create a soft piano-like chord for break completion."""
    # C major chord (C-E-G) in a comfortable range
    chord_frequencies = [261.63, 329.63, 392.00]  # C4, E4, G4
    
    # Create a soft chord with gentle attack
    chord_sound = generate_chord(chord_frequencies, 2.0, amplitude=0.25)
    
    return chord_sound

def create_session_start():
    """Create a light notification for session start."""
    # Simple two-tone notification
    tone1 = generate_sine_wave(440, 0.3, amplitude=0.3)  # A4
    silence = np.zeros(int(0.1 * 44100))  # 100ms silence
    tone2 = generate_sine_wave(523.25, 0.3, amplitude=0.3)  # C5
    
    notification = np.concatenate([tone1, silence, tone2])
    return notification

def main():
    """Generate all sound files."""
    sounds_dir = os.path.join(os.path.dirname(__file__), '..', 'internal', 'audio', 'sounds')
    os.makedirs(sounds_dir, exist_ok=True)
    
    print("Generating copyright-free notification sounds...")
    
    # Generate pomodoro completion sound
    print("Creating pomodoro completion sound...")
    pomodoro_sound = create_pomodoro_complete()
    save_wav(os.path.join(sounds_dir, 'pomodoro_complete.wav'), pomodoro_sound)
    
    # Generate break completion sound
    print("Creating break completion sound...")
    break_sound = create_break_complete()
    save_wav(os.path.join(sounds_dir, 'break_complete.wav'), break_sound)
    
    # Generate session start sound
    print("Creating session start sound...")
    start_sound = create_session_start()
    save_wav(os.path.join(sounds_dir, 'session_start.wav'), start_sound)
    
    print(f"Sound files generated in: {sounds_dir}")
    print("All sounds are copyright-free and safe for distribution.")

if __name__ == "__main__":
    main()