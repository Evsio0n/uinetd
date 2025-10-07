#!/bin/sh
# stop and disable service
if command -v systemctl >/dev/null 2>&1; then
  if systemctl is-active uinetd >/dev/null 2>&1; then
    systemctl stop uinetd >/dev/null 2>&1 || true
  fi
  if systemctl is-enabled uinetd >/dev/null 2>&1; then
    systemctl disable uinetd >/dev/null 2>&1 || true
  fi
fi


