# superpi-shutdown-libretro

Safe shutdown/reset software for Retroflag SuperPI case running Lakka or LudOS


## SuperPI case usage

Power button will start or shutdown the system.

Reset button will restart RetroArch/Ludo, or reboot the system if pressed for more than 2 seconds.


## Installation

These are the rough steps to install superpi-shutdown-libretro:

- Install Lakka or LudOS as usual.
- Configure network and enable SSH service (Settings -> Services -> SSH).
- Take note of the Raspberry Pi IP address (Main Menu -> Information -> Network Information)
- Copy the `superpi-shutdown-libretro` binary to the Raspberry Pi via SSH.
- Access the Raspberry PI using SSH and run the following commands, assuming that the binary was copied to current directory:
  ```
  # mkdir -p /storage/bin
  # mv superpi-shutdown-libretro /storage/bin/
  # chmod +x /storage/bin/superpi-shutdown-libretro
  # echo "/storage/bin/superpi-shutdown-libretro &" >> /storage/.config/autostart.sh
  # reboot
  ```
