#!/bin/sh

CURRENT_DIR=$(cd $(dirname $0) || exit 1;pwd)
LOG_DIR=${CURRENT_DIR}
LOG_FILE=${LOG_DIR}/installer.log
EDGE_DAEMON_DIR=/opt/IoTEdge/edgeDaemon
EDGE_INSTALLER_DIR=/opt/IoTEdge-Installer

DESIRE_CPU_NUM=1
DESIRE_MEM_CAP=$(( 768 * 1024))
DESIRE_GLIBC_VESION=2.17
DESIRE_DOCKER_VERSION=17.06

INPUT_ARGU_NUM=6
PROCEDURE_NUM=5

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

initLog()
{
    if [ ! -e "${LOG_DIR}" ]; then
        mkdir "${LOG_DIR}"
    fi

    if [ -e "${LOG_FILE}" ]; then
        rm "${LOG_FILE}"
    fi

    touch "${LOG_FILE}"
}

checkDocker()
{
    LOG_INFO ""
    LOG_INFO "(1/${PROCEDURE_NUM}) check docker version"
    docker_info=$(docker info 2>&1)

    result=$(docker info 2>&1 | sed -n '1p' | grep Cannot | grep connect | grep Docker)
    if [ -n "${result}" ]; then
        LOG_ERROR "docker is not start, please start first."
        return 1
    fi

    result=$(docker info 2>&1 | sed -n '1p' | grep command | grep not | grep found)
    if [ -n "${result}" ]; then
        LOG_ERROR "docker is not installed, please install first."
        return 1
    fi

    result=$(docker info 2>&1 | sed -n "/Version/p" | grep Server)
    if [ -n "${result}" ]; then
        LOG_INFO "got docker ${result}"
        version=$(echo ${result} | awk '{print $3}')
        big_version=$(echo ${version} | awk -F '[.]' '{print $1}')
        small_version=$(echo ${version} | awk -F '[.]' '{print $2}')
        docker_verison="${big_version}.${small_version}"
        if [ $(expr ${docker_verison} \>= ${DESIRE_DOCKER_VERSION}) -eq 1 ]; then
            LOG_INFO "(1/${PROCEDURE_NUM}) check docker version success"
            return 0
        else
            LOG_ERROR "docker version is too low, mini support version is ${DESIRE_DOCKER_VERSION}."
            return 1
        fi
    fi
    LOG_ERROR "unknown error occurred."
    return 1
}

checkCpu()
{
    LOG_INFO ""
    LOG_INFO "(2/${PROCEDURE_NUM}) check cpu info"
    cpu_num=$(sed -n "/processor/p" /proc/cpuinfo | wc -l)
    LOG_INFO "cpu num is ${cpu_num}"
    if [ ${cpu_num} -ge ${DESIRE_CPU_NUM} ]; then
      LOG_INFO "(2/${PROCEDURE_NUM}) check cpu success "
      return 0
    fi

    LOG_ERROR "require mini cpu num is ${DESIRE_CPU_NUM}."
    return 1
}

checkMemory()
{
    LOG_INFO ""
    LOG_INFO "(3/${PROCEDURE_NUM}) check memory info"
    mem_info=$(sed -n "/MemTotal/p" /proc/meminfo)
    LOG_INFO ${mem_info}
    mem_size=$(echo ${mem_info} | awk '{print $2}')
    if [ ${mem_size} -ge ${DESIRE_MEM_CAP} ];then
      LOG_INFO "(3/${PROCEDURE_NUM}) check memory success"
      return 0
    fi

    LOG_ERROR "MemTotal should be more than ${DESIRE_MEM_CAP} kb."
    return 1
}

checkGlibc()
{
    LOG_INFO ""
    LOG_INFO "(4/${PROCEDURE_NUM}) check glibc info"
    glibc_info=$(getconf GNU_LIBC_VERSION)
    LOG_INFO ${glibc_info}
    glibc_version=$(echo ${glibc_info} | awk '{print $2}')
    if [ $(expr ${glibc_version} \>= ${DESIRE_GLIBC_VESION}) -eq 1 ];then
      LOG_INFO "(4/${PROCEDURE_NUM}) check glibc success"
      return 0
    fi

    LOG_ERROR "mini support glibc version is ${DESIRE_GLIBC_VESION}."
    return 1
}

installDaemon()
{
    LOG_INFO ""
    LOG_INFO "(5/${PROCEDURE_NUM}) start install daemon"
    if [ -e "/lib/systemd/system/edgedaemon.service" ];then
        if [ -e "${EDGE_INSTALLER_DIR}/edgeInstaller" ];then
            LOG_ERROR "IoT Edge Lite have already installed in the host, please uninstall first."
            return 1
        fi
        systemctl stop edgedaemon
        sleep 1
        systemctl disable edgedaemon
        sleep 1
        systemctl daemon-reload
        sleep 1
        rm -f /lib/systemd/system/edgedaemon.service
    fi
    if [ -e "${EDGE_DAEMON_DIR}" ];then
        rm -rf ${EDGE_DAEMON_DIR}
    fi
    mkdir -p ${EDGE_DAEMON_DIR}

    tar zxf edgeRuntime.tar.gz
    cp -r edgeDaemon/* ${EDGE_DAEMON_DIR}

    mkdir -p /lib/systemd/system/
    cp ${EDGE_DAEMON_DIR}/scripts/edgedaemon.service /lib/systemd/system/

    chmod 644 /lib/systemd/system/edgedaemon.service
    chmod 755 ${EDGE_DAEMON_DIR}/bin/edgeDaemon
    chmod 600 ${EDGE_DAEMON_DIR}/conf/rootcert.pem
    LOG_INFO "add daemon to systemd."
    mkdir -p /var/IoTEdge/log/edge_daemon
    systemctl enable edgedaemon
    systemctl start edgedaemon
    LOG_INFO "(5/${PROCEDURE_NUM}) install daemon success"
    return 0
}

checkRun()
{
    if [ 0 -ne $1 ]; then
        LOG_ERROR $2
        exit 1
    fi
}

main()
{
    LOG_INFO "Current Work dir is ${CURRENT_DIR}."
    LOG_INFO "Log written to ${LOG_FILE}."
    LOG_INFO "Start to install edgeDaemon."

    if [ $# -lt ${INPUT_ARGU_NUM} ]; then
        LOG_ERROR "input arguement error."
        exit 1
    fi

    SERVER_ADDRESS=$1
    SERVER_PORT=$2
    NODE_ID=$3
    MODULE_ID=$4
    INSTALLER_TOKEN=$5
    PROJECT_ID=$6

    # initLog
    checkDocker
    checkRun $? "(1/${PROCEDURE_NUM}) check docker version failed"
    sleep 1

    checkCpu
    checkRun $? "(2/${PROCEDURE_NUM}) check cpu failed"
    sleep 1

    checkMemory
    checkRun $? "(3/${PROCEDURE_NUM}) check memory failed"
    sleep 1

    checkGlibc
    checkRun $? "(4/${PROCEDURE_NUM}) check glibc failed"
    sleep 1

    installDaemon
    checkRun $? "(5/${PROCEDURE_NUM}) install failed"
    sleep 5

    mkdir -p ${EDGE_INSTALLER_DIR}
    cp -r edgeInstaller/* ${EDGE_INSTALLER_DIR}
    chmod 600 ${EDGE_INSTALLER_DIR}/install.dat

    LOG_INFO "run edgeInstaller..."
    umask 0077
    export LD_LIBRARY_PATH=${EDGE_INSTALLER_DIR}:${LD_LIBRARY_PATH}
    export ASAN_OPTIONS=halt_on_error=0:detect_leaks=1:malloc_context_size=15:log_path=/var/IoTEdge/log/edge_installer/asan.log
    ${EDGE_INSTALLER_DIR}/edgeInstaller -a install -s ${SERVER_ADDRESS} -p ${SERVER_PORT} -n ${NODE_ID} -m ${MODULE_ID} -r ${INSTALLER_TOKEN} -P ${PROJECT_ID}
}

main "$@"