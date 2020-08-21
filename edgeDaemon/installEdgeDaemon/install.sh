#!/sbin/sh
remoteUrl=${1}
installerToken=${3}
id=${2}
flag=1
edgeDaemonWorkplace="/opt/easyfetch/edgeDaemon"
if [ -d ${edgeDaemonWorkplace} ]; then
  rm -rf ${edgeDaemonWorkplace} || echo "ERROR: REMOVE WORKPLACE ERR" || flag=0
else
  mkdir -p ${edgeDaemonWorkplace} || echo "creat working directory error" || flag=0

fi
cp ./* ${edgeDaemonWorkplace} || echo "ERROR:can't copy files" || flag=0
cd
if [ -e "${edgeDaemonWorkplace}/edgeDaemon" ]; then
  if [ -e "${edgeDaemonWorkplace}/edgeDaemon.service" ]; then
    cp ${edgeDaemonWorkplace}/edgeDaemon.service /etc/systemd/system || echo "ERROR: CANT COPY SERVER" || flag=0
    systemctl daemon-reload
    systemctl enable edgeDaemon.service
    systemctl restart edgeDaemon || echo "ERROR: CANT START THE SERVICE" || flag=0
  else
    echo "ERROR: no edgeDaemon.service"
  fi
else
  echo "ERROR: no edgeDaemon file"
fi
if [ -e "${edgeDaemonWorkplace}/jsonCreater" ]; then
  cd ${edgeDaemonWorkplace}
  ./jsonCreater "-url" ${remoteUrl} "-id" ${id} "-token" ${installerToken} || echo "ERROR: CREAT CONFIGFILE ERROR" || flag=0
else
  echo "ERROR: NO jsonCreater"
fi

if [ ${flag} -eq 1 ]; then
  echo "INFO: install success!"
else
  echo "ERROR: install faild"
fi
