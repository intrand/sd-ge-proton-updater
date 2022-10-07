# About

This will automatically install the latest stable version of GloriousEggroll's proton-ge-custom fork. GE-Proton includes many improvements to run games more betterer, including many of the games marked as "unsupported" by Valve.

# How it works

Steam Deck GE-Proton Updater will run at every boot. In order, it will...
1. check for the latest version of itself and attempt to update itself if a newer version is found. You shouldn't need to worry about getting a newer version manually, but you still can try if you want to by clicking the Update button.
2. ensure that itself is configured correctly. This generally means that the binary is installed to the correct location (`~/.sd-ge-proton-updater`), that the binary is executable and that the systemd unit is installed and set to run on boot.
3. attempt to get information about the latest GE-Proton release,
4. install the latest release if you don't already have it.

# Caveats

Unfortunately, the Steam UI, even in Gaming Mode, needs to be restarted for new versions of Proton to show up. So, even if you get an update, you may need to reboot for it to become available in the Steam interface.

# Installation & Use

Follow the very simple steps in the latest release.
