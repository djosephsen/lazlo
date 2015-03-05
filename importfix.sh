#!/bin/sh
#attempt to fix broken Go import paths as a consequence of forking the repo

#sanity check
[ "${GOPATH}" ] || fail 'Sorry you need to set GOPATH'

#fix sed for mactards (like me)
if [ $(uname) == 'darwin' ]
then	
	s='sed -E'
else
	s='sed'
fi

# setup some functions
function fail {
	echo "$@"
	exit 42
}

function pathEscape {
	echo "${@}" | sed -e 's/\//\\\//g'
}

function ZOMGDOIT {
# cross your fingers
for FILE in ${FLIST}
do
	cat ${FILE} | sed -e "s/$(pathEscape ${OGPATH})/$(pathEscape ${OUR_PATH})/"> ${TMP} && mv ${TMP} ${FILE}
done
rm -f ${TMP}
}

function dryRun {
echo dryRun baby
for FILE in ${FLIST}
do
	LINES=$(grep "${OGPATH}" ${FILE})
	if [ -n "${LINES}" ]
	then
		echo "In file: ${FILE}:" 
		echo "${LINES}"| while read LINE
		do
			echo "I'd replace: ${LINE}"
			echo  "with:        "$(echo ${LINE} | sed -e "s/$(pathEscape ${OGPATH})/$(pathEscape ${OUR_PATH})/")
		done
		echo 
	fi
done
}

#ok lets see here..
TMP='/tmp/SFTEMPFILE'
PACKAGE=$(basename $(pwd))
FLIST=$(find . -name '*.go')
MINUS_THIS="${GOPATH}/src/"
OGPATH=$(echo ${GOPATH}/src/github.com/djosephsen/slacker | sed "s/$(pathEscape ${MINUS_THIS})//")
OUR_PATH=$(pwd | sed -e "s/$(pathEscape ${MINUS_THIS})//")

if [ -z ${1} ]
then
	ZOMGDOIT
else
	dryRun
fi
