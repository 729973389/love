#!/sbin/sh
flag=1
edgeDaemonWorkplace="/opt/easyfetch/edgeDaemon"
if [ -e "/tmp/installEdgeDaemon.tar.gz" ]; then
	rm /tmp/installEdgeDaemon.tar.gz || echo "WARNING: REMOVE .TAR.GZ ERROR" ||flag=0
fi
if [ -d "/tmp/installEdgeDaemon" ]; then 
	rm -rf /tmp/installEdgeDaemon || echo "WARNING: REMOVE INSTALL FILE ERROR" || flag=0
fi
if [ -e "/etc/systemd/system/edgeDaemon.service" ]; then
	systemctl stop edgeDaemon.service;systemctl disable edgeDaemon.service;systemctl daemon-reload || echo "ERROR: STOP THE SERVICE FAILED" || flag=0
	rm /etc/systemd/system/edgeDaemon.service || echo "ERROR: REMOVE SERVICE ERROR" || flag=0
fi
if [ -d ${edgeDaemonWorkplace} ]; then 
	rm -rf ${edgeDaemonWorkplace} ||echo "ERROR: REMOVE WORKPLACE ERROR" || flag=0
	echo "INFO: REMOVE WORKPLACE SUCCESS"
fi
if [ ${flag} -eq 1 ]; then
       echo "INFO: UNINSTALL SUCCESS"
fi       

