#!/sbin/sh
remoteUrl=${1}
installerToken=${2}
id=${3}
edgeFlag=0
edgeDaemonWorkplace="/opt/easyfetch/edgeDaemon"
if [ -d ${edgeDaemonWorkplace} ]; then
	rm -rf ${edgeDaemonWorkplace} || echo "ERROR: REMOVE WORKPLACE ERR"
else
       	mkdir -p ${edgeDaemonWorkplace} || echo "creat working directory error"

fi
cp ./* ${edgeDaemonWorkplace} || echo "ERROR:can't copy files"
if [ -e "${edgeDaemonWorkplace}/edgeDaemon" ]; then
	edgeFlag=1
	if [ -e "${edgeDaemonWorkplace}/edgeDaemon.service" ]; then 
		edgeFlag=2
		cp ${edgeDaemonWorkplace}/edgeDaemon.service /etc/systemd/system
		systemctl daemon-reload;systemctl enable edgeDaemon.service;systemctl restart edgeDaemon;
	else
		echo "ERROR: no edgeDaemon.service"
	fi
else
	echo "ERROR: no edgeDaemon file"
fi
if [ ${edgeFlag} -eq 2 ]; then
	echo "INFO: install success!"
else
		echo "ERROR: install faild"
fi
if [ -e "${edgeDaemonWorkplace}/jsonCreater" ]; then 
	cd ${edgeDaemonWorkplace};${edgeDaemonWorkplace}/jsonCreater "-url" ${remoteUrl} "-token" ${installerToken} "-id" ${id}||echo "ERROR: create configfile failed"
fi
