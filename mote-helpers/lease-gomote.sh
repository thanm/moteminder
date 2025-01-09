#!/bin/sh
set -x
#
# Pseudocode:
# - see if mote.txt exists
#   + read and ping
# - if new mote needed, gomote create
#
#-------------
#
# Read mote flavor file.
#
MOTEFLAVOR=`cat moteflavor.txt`
if [ -z "$MOTEFLAVOR" ]; then
    echo "error: missing or empty moteflavor.txt file"
    exit 1
fi
FLAVTEST=`gomote create -list  | egrep "^${MOTEFLAVOR}\$"`
if [ "$FLAVTEST" != "$MOTEFLAVOR" ]; then
    echo "warning: moteflavor $MOTEFLAVOR not recognized by gomote create -list"
    exit 1
fi
HERE=`pwd`
BN=`basename $HERE`
if [ "$FLAVTEST" != "$MOTEFLAVOR" ]; then
    echo "warning: dirname $BN does not match mote flavor $MOTEFLAVOR"
fi
#
# Do we have an active mote?
MOTE=""
bash check-mote.sh
if [ $? != 0 ]; then
    echo "... new lease needed"
else
    MOTE=`cat mote.txt`
fi
if [ -z "$MOTE" ]; then
    # Kick off creation
    echo "... kicking off creation of $MOTEFLAVOR gomote"
    gomote create $MOTEFLAVOR > /tmp/newmote.${MOTEFLAVOR}.txt
    if [ $? != 0 ]; then
	echo "error: gomote create failed"
	exit 1
    fi
    MOTE=`cat /tmp/newmote.${MOTEFLAVOR}.txt`
    echo "... done: new mote is $MOTE"
    cp  /tmp/newmote.${MOTEFLAVOR}.txt mote.txt
fi
