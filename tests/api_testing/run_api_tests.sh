#!/bin/bash
GREEN='\033[0;32m' # Green color
NC='\033[0m'       # No color

usage()
{
	echo "This script to run itsyou.online api test suite on local environment"
	echo -e "\nUsage:\n$0 [options] \n"
	echo "Options:"
	echo " -d : directory of the testsuite"
}

if [[($1 == "--help") || ($1 == "-h")]]; then
	usage
	exit 0
elif [[($1 == "-d")]] ; then
	directory=$2
else
	echo "Invalid option $1"
	echo "for help use -h or --help"
	exit 0
fi

echo -e "${GREEN}** Checking python-pip ...${NC}";
which pip2 || apt-get install -y python-pip
echo -e "${GREEN}** Activating JumpScale virtual env ...${NC}"
pip2 install virtualenv
virtualenv venv
source venv/bin/activate
echo -e "${GREEN}** Installing G8_testing requirements ...${NC}"
apt-get install build-essential libssl-dev libffi-dev python-dev
pip2 install -r requirements.txt
echo -e "${GREEN}** Running tests ...${NC}"
nosetests -v --logging-level=WARNING $directory --tc-file config.ini --with-xunit --xunit-file='testresults.xml' --with-progressive
# Collect result
echo -e "${GREEN}** DONE ** ...${NC}"
