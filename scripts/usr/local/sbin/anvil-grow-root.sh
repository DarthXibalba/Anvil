#!/usr/bin/env bash
set -euo pipefail

echo "== Anvil: Growing root filesystem =="

ROOT_SRC=$(findmnt -n -o SOURCE /)
DISK="/dev/$(lsblk -no PKNAME "$ROOT_SRC" 2>/dev/null || true)"

if [[ -z "$DISK" ]]; then
  echo "ERROR: Unable to determine root disk"
  exit 1
fi

echo "Root device: $ROOT_SRC"
echo "Parent disk: $DISK"

# Case 1: LVM
if [[ "$ROOT_SRC" == /dev/mapper/* ]]; then
  VG=$(lvs --noheadings -o vg_name "$ROOT_SRC" | xargs)
  LV=$(lvs --noheadings -o lv_name "$ROOT_SRC" | xargs)

  echo "Detected LVM root: $VG/$LV"

  sudo pvresize "$DISK"
  sudo lvextend -l +100%FREE "/dev/$VG/$LV"

  FS=$(findmnt -n -o FSTYPE /)
  case "$FS" in
    ext4) resize2fs "/dev/$VG/$LV" ;;
    xfs)  xfs_growfs / ;;
    *)
      echo "Unsupported filesystem: $FS"
      exit 1
      ;;
  esac

  echo "LVM root filesystem grown successfully."
  exit 0
fi

# Case 2: Single-partition root (vda1)
PART="${ROOT_SRC##*/}"
PART_NUM="${PART##*[!0-9]}"

if [[ -z "$PART_NUM" ]]; then
  echo "ERROR: Unable to determine partition number"
  exit 1
fi

echo "Detected partition root: $DISK$PART_NUM"

sudo growpart "$DISK" "$PART_NUM"

FS=$(findmnt -n -o FSTYPE /)
case "$FS" in
  ext4) resize2fs "$ROOT_SRC" ;;
  xfs)  xfs_growfs / ;;
  *)
    echo "Unsupported filesystem: $FS"
    exit 1
    ;;
esac

echo "Partition root filesystem grown successfully."
