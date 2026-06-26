#!/usr/bin/env bash
set -euo pipefail

APP_NAME="ag-directory-manager"
INSTALL_DIR="/opt/${APP_NAME}"
DATA_DIR="/var/lib/${APP_NAME}"
LOG_DIR="/var/log/${APP_NAME}"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

if [[ "${EUID}" -ne 0 ]]; then
  echo "Execute como root." >&2
  exit 1
fi

mkdir -p "${INSTALL_DIR}" "${DATA_DIR}" "${LOG_DIR}"

cat > "${SERVICE_FILE}" <<EOF
[Unit]
Description=AG Directory Manager
After=network.target

[Service]
Type=simple
WorkingDirectory=${INSTALL_DIR}
ExecStart=${INSTALL_DIR}/agdm-server
Restart=always
RestartSec=5
Environment=APP_ENV=production
Environment=APP_SQLITE_PATH=${DATA_DIR}/agdm.db

[Install]
WantedBy=multi-user.target
EOF

chmod 644 "${SERVICE_FILE}"
systemctl daemon-reload
systemctl enable "${APP_NAME}"

echo "Instalação base concluída. Copie o binário para ${INSTALL_DIR}/agdm-server e inicie com: systemctl start ${APP_NAME}"
