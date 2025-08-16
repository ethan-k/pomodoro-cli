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
    """Create a simple, quiet bell sound for pomodoro completion."""
    # Simple bell tone - much quieter and simpler
    fundamental = 800  # Higher pitch bell sound
    
    # Create a simple bell with natural decay
    duration = 1.5  # Shorter duration
    sample_rate = 44100
    frames = int(duration * sample_rate)
    
    bell_sound = []
    for i in range(frames):
        # Simple sine wave with exponential decay
        t = i / sample_rate
        decay = math.exp(-t * 2)  # Exponential decay
        amplitude = 0.15 * decay  # Much quieter (was 0.4)
        
        # Simple sine wave
        value = amplitude * math.sin(2 * math.pi * fundamental * t)
        
        # Gentle fade in for first 10ms to avoid click
        fade_frames = int(0.01 * sample_rate)
        if i < fade_frames:
            value *= i / fade_frames
            
        bell_sound.append(value)
    
    return np.array(bell_sound)

def create_break_complete():
    """Create a simple, soft tone for break completion."""
    # Single soft tone instead of chord - much simpler
    frequency = 523.25  # C5
    
    # Create a simple tone with gentle envelope
    duration = 1.0  # Shorter
    tone_sound = generate_sine_wave(frequency, duration, amplitude=0.12)  # Much quieter
    
    return tone_sound

def create_session_start():
    """Create a very quiet notification for session start."""
    # Simple single soft tone
    tone = generate_sine_wave(660, 0.5, amplitude=0.08)  # E5, very quiet
    
    return tone

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