#!/sbin/sh
edgeFlag=0
edgeDaemonWorkplace="/opt/easyfetch/edgeDaemon"
if [ ! -d ${edgeDaemonWorkplace} ]; then
  mkdir -p ${edgeDaemonWorkplace}
else
	  echo "WARNING: ${edgeDaemonWorkplace} exists"
fi
cp ./* ${edgeDaemonWorkplace}
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
