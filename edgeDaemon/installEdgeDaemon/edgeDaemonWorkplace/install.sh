#!/sbin/sh
edgeDaemonWorkplace="/opt/easyfetch/edgeDaemon"
if [ ! -d $edgeDaemonWorkplace ]; then
  mkdir -p edgeDaemonWorkplace
fi
cp * edgeDaemonWorkplace
echo "install finished"
