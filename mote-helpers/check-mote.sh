#/bin/bash
#
# Do we have an active mote?
#
MOTE=""
if [ -f mote.txt ]; then
    M=`cat mote.txt`
    if [ -z "$M" ]; then
	echo "... mote.txt empty, new lease needed"
	exit 1
    else
	gomote ping $M
	if [ $? != 0 ]; then
	    echo "... ping of $M failed, new lease needed"
	    exit 2
	else
	    RESULT=`gomote ping $M 2>&1`
	    echo $RESULT | fgrep -q alive
	    RC=$?
	    if [ $RC != 0 ]; then
		echo "... mote $M ping bad result, new lease needed"
		exit 3
	    else
		echo "... existing mote $M looks good."
		MOTE=$M
	    fi
	fi
    fi
fi
exit

