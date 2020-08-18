#!/bin/sh

EDGE_INSTALLER_DIR=/opt/IoTEdge-Installer
EDGE_AGENT_DIR=/opt/IoTEdge/edgeAgent
EDGE_DAEMON_DIR=/opt/IoTEdge/edgeDaemon

LOG_ERROR()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | ERROR | $*"
    echo "${message}"
    # echo "${message}" >> "${LOG_FILE}" 2>&1
}

LOG_INFO()
{
    message="$(date +%Y-%m-%dT%H:%M:%S) | INFO | $*"
    echo "${message}"
    # echo "${message}" >> "${LOG_FILE}" 2>&1
}


main()
{
    LOG_INFO "delete containers."
    export LD_LIBRARY_PATH=${EDGE_INSTALLER_DIR}:$LD_LIBRARY_PATH
    umask 0077
    ${EDGE_INSTALLER_DIR}/edgeInstaller -a uninstall

    LOG_INFO "delete edgeAgent."
    rm -rf ${EDGE_AGENT_DIR}

    LOG_INFO "unregister edgeDaemon service."
    systemctl stop edgedaemon
    sleep 1
    systemctl disable edgedaemon
    sleep 1
    systemctl daemon-reload
    sleep 1
    rm -f /lib/systemd/system/edgedaemon.service
    rm -rf ${EDGE_DAEMON_DIR}

    LOG_INFO "delete edgeRuntime datas."
    rm -rf /var/IoTEdge

    LOG_INFO "delete edgeInstaller"
    rm -rf ${EDGE_INSTALLER_DIR}

    LOG_INFO "uninstall edgeRuntime finish."
}

main