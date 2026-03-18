#!/bin/sh
set -eu

CTF_WEB_PORT="${CTF_WEB_PORT:-7681}"
QUIZ_BASE_URL="${QUIZ_BASE_URL:-http://quiz.ktf.ninja}"
INSTALL_TEMPLATE="/app/install.html.template"
INSTALL_HTML="/usr/share/nginx/html/install.html"
INSTALLER_PATH="/app/QuizLauncherInstaller-macos.zip"

cleanup() {
  if [ -n "${CTF_PID:-}" ] && kill -0 "${CTF_PID}" 2>/dev/null; then
    kill "${CTF_PID}" 2>/dev/null || true
    wait "${CTF_PID}" 2>/dev/null || true
  fi
}

trap cleanup EXIT INT TERM

if [ ! -x /app/ctf-tool ]; then
  echo "ctf-tool binary missing" >&2
  exit 1
fi

# Render install page (if template + installer present)
if [ -f "${INSTALL_TEMPLATE}" ]; then
  mkdir -p /usr/share/nginx/html
  sed "s|__QUIZ_BASE_URL__|${QUIZ_BASE_URL}|g" "${INSTALL_TEMPLATE}" > "${INSTALL_HTML}"
fi

# Start the ctf-tool built-in web server (replaces ttyd)
/app/ctf-tool -web -port "${CTF_WEB_PORT}" &
CTF_PID=$!

# Start nginx in foreground
exec nginx -g "daemon off;"
