#!/sbin/sh
flag=2
edgeDaemonWorkplace="/opt/easyfetch/edgeDaemon"
if [ -e "/etc/systemd/system/edgeDaemon.service" ]; then
	flg=0
	systemctl stop edgeDaemon.service;systemctl disable edgeDaemon.service;systemctl daemon-reload || echo "ERROR: STOP THE SERVICE FAILED"
	rm /etc/systemd/system/edgeDaemon.service || echo "ERROR: REMOVE SERVICE ERROR"
	echo "INFO: REMOVE SERVICE SUCCESS"
	flag=2
fi
if [ -d ${edgeDaemonWorkplace} ]; then 
	rm -rf ${edgeDaemonWorkplace} ||echo "ERROR: REMOVE WORKPLACE ERROR"
	echo "INFO: REMOVE WORKPLACE SUCCESS"
	flag=2
fi
if [ ${flag} -eq 2 ]; then
       echo "INFO: UNINSTALL SUCCESS"
fi       

