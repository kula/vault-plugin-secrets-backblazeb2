#!/bin/bash

PLUGIN_NAME=vault-plugin-secrets-backblazeb2
PLUGIN_PATH=b2
TOPDIR="$( git rev-parse --show-toplevel )"
PLUGIN="${TOPDIR}/${PLUGIN_NAME}"
WORK="${TOPDIR}/_workspace"
PLUG_DIR="${WORK}/plugins"
TEST="${TOPDIR}/test"
PID_FILE="${WORK}/vault.pid"
LOG_FILE="${WORK}/vault.log"


. ${TEST}/env.sh


die() { 
    echo "### Error: $@" >&2
    exit 1
}

cleanup() {
    if [ -f "$PID_FILE" ]; then
	kill -INT $(cat "$PID_FILE") >/dev/null 2>&1
	rm -f "${PID_FILE}" >/dev/null 2>&1
    fi
}

trap cleanup EXIT


which vault >/dev/null || die "Cannot find vault binary in your path"

[[ -x "$PLUGIN" ]] || die "$PLUGIN_NAME doesn't exist, make it?"
mkdir -p "${PLUG_DIR}" || die "Cannot make work directory"
rm -f "${PLUG_DIR}/${PLUGIN_NAME}" || die "Cannot delete old plugin"

echo "### Starting vault server"

nohup vault server -dev \
    -dev-listen-address="${VAULT_IP}:${VAULT_PORT}" \
    -dev-root-token-id="${VAULT_TOKEN}" \
    -dev-plugin-dir="${PLUG_DIR}" \
    -log-level=debug >> "${LOG_FILE}" 2>&1 &
echo $! > "${PID_FILE}"

ps -p $(cat ${PID_FILE} ) >/dev/null || die "Could not start vault"

echo "### Vault started"
echo "### Log file: ${LOG_FILE}"
echo "### Pid file: ${PID_FILE}"
echo
echo "### Copying and registering plugin"

cp "${TOPDIR}/${PLUGIN_NAME}" "${PLUG_DIR}" || die "Cannot copy ${PLUGIN_NAME} into plugin dir"
SUM=$( sha256sum ${PLUG_DIR}/${PLUGIN_NAME} 2>/dev/null | cut -d " " -f 1 )
[[ -n "$SUM" ]] || die "Could not calculate plugin sha256 sum"

vault plugin register \
    -command="${PLUGIN_NAME}" \
    -sha256="${SUM}" \
    ${PLUGIN_NAME} || die "Could not register plugin in vault"

vault secrets enable \
    -path=${PLUGIN_PATH} \
    -plugin-name=${PLUGIN_NAME} \
    plugin || die "Could not enable plugin"


echo
echo "Plugin enabled at ${PLUGIN_NAME}/"
echo "Starting shell, exit to stop vault server"
echo
echo

PS1="vault-testing: " /bin/bash --noprofile --norc
