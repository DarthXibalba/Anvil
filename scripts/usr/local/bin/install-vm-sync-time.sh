#!/usr/bin/env bash
sudo apt update
sudo apt install qemu-guest-agent chrony
sudo systemctl enable --now qemu-guest-agent chrony
