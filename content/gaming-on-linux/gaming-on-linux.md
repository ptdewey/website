---
title: "Gaming on Linux (Wayland)"
authors: ["Patrick Dewey"]
date: 2025-08-11
categories: ["Linux"]
tags: ["linux", "linux-gaming", "proton", "gamescope"]
type: "Blog"
bluesky_link: "at://did:plc:hm5f3dnm6jdhrc55qp2npdja/app.bsky.feed.post/3lw3fp652gc2l"
---

## Introduction

Gaming on Linux has come a long way in recent years, but it still presents unique challenges compared to Windows gaming. While tools like Proton and the advent of the Steam Deck have drastically improved the out-of-the-box Linux gaming experience, there are still many issues. Many of these issues are related to hardware/software compatibility, and can often be found when using modern Wayland window managers that lack full X11 support.

This will not be a comprehensive guide, but it includes some tips and tools I've found to be helpful with getting Linux gaming working.

## My Setup

My gaming system runs NixOS Unstable (currently version 25.11), with the Niri compositor. I use a recent AMD GPU since they have significantly fewer issues on Wayland compared to Nvidia cards (I won't be talking about Nvidia-specific issues here).

One of the primary challenges I've run into with Linux gaming is incomplete (or lack of) X11 support. This is a problem as many Steam games only support X11, creating friction when running on pure Wayland environments. In my experience with Niri, I've encountered two particularly frustrating issues: controllers failing to be detected by games, and X11-only Steam games crashing immediately on startup. Another problem is that there are relatively few people using these newer tiling window managers, and as such, there is a lack of good documentation around using them for gaming.

## Solutions

### Proton

Proton[^proton] has been an immensely important tool for Linux gaming, as It allows Windows-exclusive games to be run on Linux using Wine. Proton can be enabled in your Steam settings under the "compatibility" tab. For most Linux users, Proton is likely the only thing from this guide you will need (unless you require additional setup-specific configuration).

### Steam Launch Flags

For window-managers that may not work well out of the box with just Proton, Steam has a couple of useful launch flags that can solve some issues.

Running steam with the `-steamos3` launch flag tells steam to run with some more controller friendly options. I added custom desktop entry files to `~/.local/share/applications/` to ensure steam always launches with the flag enabled.

```desktop
# ~/.local/share/applications/steam.desktop
[Desktop Entry]
Name=SteamUI
Comment=Run steam in steamos3 mode
Exec=steam -steamos3
Icon=steam
Terminal=false
Type=Application
Categories=Game;
```

Another useful flag is `-gamepadui`. It instructs steam to start in big picture mode for a console-like (controller-driven) UI.

```desktop
# ~/.local/share/applications/steam-ui.desktop
[Desktop Entry]
Name=SteamUI (Big Picture)
Comment=Run steam in gamepad mode
Exec=steam -steamos3 -gamepadui
Icon=steam
Terminal=false
Type=Application
Categories=Game;
```

I did run into an issue where my controller can't send inputs to a game using `-gamepadui` and `gamescope` at the same time. Note that this seems to be a more setup-specific issue (with Niri), so your mileage may vary.

### Gamescope

Gamescope[^gamescope] is another useful tool for running games in a Wayland environment. It runs X11 applications with XWayland, allowing Steam games (and other X-only applications) to be run in Wayland environments that lack support for the X server (e.g. Niri). It is packaged for most Linux distros or can be built from source[^gamescope-packages].

To launch a game with Gamescope, add this to the Steam launch settings for each desired game with modifications for your screen width/height (`-W, -w, -H, -h` flags), and desired frame-rate (`-r` flag). Leave `%command%` as is, as this is automatically populated with the path your game.
```sh
gamescope -W 2560 -H 1440 -w 2560 -h 1440 -r 60 -f -- %command%
```

See the [Gamescope repo](https://github.com/ValveSoftware/gamescope) or run `gamescope --help` for more usage details.

---

[^proton]: [Proton](https://github.com/ValveSoftware/Proton)
[^gamescope]: [Gamescope](https://github.com/ValveSoftware/gamescope)

[^gamescope-packages]: [Gamescope package statuses](https://github.com/ValveSoftware/gamescope?tab=readme-ov-file#status-of-gamescope-packages)

