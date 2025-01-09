#!/bin/bash
set -x
MOTE=""
bash check-mote.sh
if [ $? != 0 ]; then
    echo "error: mote check failed, bailing"
    exit 1
else
    MOTE=`cat mote.txt`
fi
#
# Collect op
#
OP=$1
shift
ARGS=$*
if [ -z "$OP" ]; then
    echo "error: supply op as first arg"
    exit 1
fi
#
# Push first
#
gomote push $MOTE
RC=$?
if [ $RC != 0 ]; then
    echo "error: push failed"
    exit 1
fi
#
# Execute op
#
# FIXME/TODO:
# - handle windows here
# - test with darwin
#
case $OP in
    all) gomote run $MOTE bash -c "cd \$WORKDIR/go/src ; bash all.bash" ;;
    make) gomote run $MOTE bash -c "cd \$WORKDIR/go/src ; bash make.bash" ;;
    dotest) gomote run -path "\$WORKDIR/go/bin:\$PATH" $MOTE bash -c "cd \$WORKDIR/go/src ; go test $ARGS" ;;
    *) echo "unrecognized op: $OP"
       exit 1 ;;
esac
echo "... done."
exit 0

