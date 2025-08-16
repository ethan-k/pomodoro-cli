#!/usr/bin/env python3
"""
Download free bell notification sounds from freesound.org
All sounds are Creative Commons licensed and free to use.
"""

import os
import urllib.request
import urllib.parse

def download_sound(url, filename, description=""):
    """Download a sound file from a URL."""
    sounds_dir = os.path.join(os.path.dirname(__file__), '..', 'internal', 'audio', 'sounds')
    os.makedirs(sounds_dir, exist_ok=True)
    
    filepath = os.path.join(sounds_dir, filename)
    
    try:
        print(f"📥 Downloading: {description}")
        print(f"   URL: {url}")
        print(f"   Saving as: {filename}")
        
        # Add headers to simulate a browser request
        req = urllib.request.Request(url, headers={
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
        })
        
        with urllib.request.urlopen(req) as response:
            with open(filepath, 'wb') as f:
                f.write(response.read())
        
        print(f"✅ Downloaded successfully: {filename}")
        return True
        
    except Exception as e:
        print(f"❌ Failed to download {filename}: {e}")
        return False

def main():
    """Download notification bell sounds."""
    print("🔔 Downloading free bell notification sounds from freesound.org")
    print("All sounds are Creative Commons licensed and free to use.\n")
    
    # Note: Direct download URLs from freesound.org require authentication
    # These are example URLs - you would need to get actual download links
    
    sounds = [
        {
            "url": "https://freesound.org/data/previews/571/571512_11450107-lq.mp3",
            "filename": "pomodoro_complete.mp3", 
            "description": "Soft notification bell by LegitCheese (CC0)"
        }
    ]
    
    print("ℹ️  Note: Freesound.org requires user authentication for downloads.")
    print("Manual download instructions:")
    print()
    
    for sound in sounds:
        print(f"1. Visit: https://freesound.org/people/LegitCheese/sounds/571512/")
        print(f"2. Click 'Download' button (you may need to create free account)")
        print(f"3. Save as: {sound['filename']}")
        print(f"4. Description: {sound['description']}")
        print()
    
    print("Alternative sounds to search for on freesound.org:")
    print("- Search: 'bell notification' + Filter: CC0 license")
    print("- Search: 'soft ding' + Filter: CC0 license") 
    print("- Search: 'gentle chime' + Filter: CC0 license")
    print()
    print("🎯 Recommended: 0.5-2 seconds duration, quiet volume, pleasant tone")

if __name__ == "__main__":
    main()