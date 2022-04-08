#!/usr/bin/env bash

PWD=`pwd`
TF_CLIENT_PATH="${PWD}/client"
HOSTNAME="lumen.com"

# Clean-up
if [ -f "${TF_CLIENT_PATH}/.terraform.lock.hcl" ] ; then
	rm ${TF_CLIENT_PATH}/.terraform.lock.hcl ;
fi
if [ -d "${TF_CLIENT_PATH}/.terraform" ] ; then
	rm -rf ${TF_CLIENT_PATH}/.terraform ;
fi

rm -rf ~/.terraform.d/plugins/${HOSTNAME}/* ;

exit 0
