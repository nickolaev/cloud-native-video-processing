#!/bin/sh

VIDEO=${VIDEO:-video}
PORT=${PORT:-31380}
RESIZE=$1
TIMESTAMP=$2
LOGO=$3

if [ -z "${RESIZE}" ] || [  "${RESIZE}" != "low" ] && [ "${RESIZE}" != "medium" ] && [ "${RESIZE}" != "high" ] ; then
	RESIZE="source"
fi

if [ -z "${TIMESTAMP}" ] || [ "${TIMESTAMP}" != "enable" ]; then
	TIMESTAMP="disable"
fi

if [ -z "${LOGO}" ] || [ "${LOGO}" != "enable" ]; then
	LOGO="disable"
fi

echo "process-video: ${RESIZE}_${TIMESTAMP}_${LOGO}"

curl -i -XGET -H"process-video: ${RESIZE}_${TIMESTAMP}_${LOGO}" \
	      http://"${VIDEO}":"${PORT}"  \
	| mplayer -really-quiet -cache 8192 - 2>/dev/null

