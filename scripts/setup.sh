#!/usr/bin/env bash

PWD=`pwd`
TF_CLIENT_PATH="${PWD}/client"
TF_LOG_PATH="${PWD}/client/logs"

# Client
if [ ! -d "${TF_CLIENT_PATH}" ] ; then
	mkdir -p ${TF_CLIENT_PATH} ;
fi
# Logs
if [ ! -d "${TF_LOG_PATH}" ] ; then
	mkdir -p ${TF_LOG_PATH} ;
fi

exit 0