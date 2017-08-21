
## Goal of this repo:
This repo is used to execute our automated tests for the itsyouonline api testsuite.

Tests are performed using the descriptions in the Folders

**testsuite**
In this folder we have automated test testsuite for the itsyouonline api

** Installation **
Steps to install itsyouonline locally:  
1. sudo apt-get update  
2. sudo apt-get -y upgrade  
3. sudo apt-get install docker.io  
4. sudo usermod -G docker -a cloudscalers  
5. sudo docker run -it -p 27017:27017 -name=mongo mongo  
6. sudo curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz  
7. sudo tar -xvf go1.6.linux-amd64.tar.gz  
8. sudo mv go /usr/local  
9. echo "export PATH=\$PATH:/usr/local/go/bin" » .profile  
10. source ~/.profile  
11. mkdir $HOME/gopath  
12. export GOPATH=$HOME/gopath  
13. go get -u github.com/itsyouonline/identityserver  
14. cd gopath/src/github.com/itsyouonline/identityserver  
15. go build  
16. ./identityserver  


# Pre-Requirements:

1. Create two itsyouonline accounts using sign-up.

> During signing up the first account scan the QR using any QR scanner and save the secret value, we will put it in config.ini file later.

2. For each account create **APIKEY**.  [See how to create apikey](create_apikey.md)
3. Fill the data required in ```config.ini``` file according to your environment as follow:

```
itsyouonline_url          : Environment url
validation_email          : Validation Email
validation_email_password : Validation email password
user1_totp_secret         : Authentication secret of the user (1) (from step 1)
user1_username            : User (1) username
user1_password            : User (1) password
user1_applicationid       : User (1) apikey applicationid
user1_secret              : User (1) apikey secret
user2_username            : User (2) username
user2_password            : User (2) password
user2_applicationid       : User (2) apikey applicationid
user2_secret              : User (2) apikey secret
```

<a href="vmail_config.md"><b> See how to configure verification email</b></a>



# Requirements:

If you don't have python 2.7 use this commands to install:
-----------------------------------------------------------
```
sudo add-apt-repository ppa:fkrull/deadsnakes
sudo apt-get update
sudo apt-get install python2.7
sudo apt-get install python-setuptools python-dev build-essential
sudo easy_install pip
```

Install Python Packages:
------------------------
Note That: you may use virtual env for this step
```
pip install virtualenv
virtualenv venv
source venv/bin/activate
pip install -r requirements.txt
```

Run the tests
---------------
change the necessary parameters in config.ini according to your environment
```
(venv)$> nosetests -xv testsuite --tc-file config.ini  2>testresults.log
```

or overwrite it using the following command
```
(venv)$> nosetests -xv testsuite --tc-file config.ini --tc=main.url:https://itsyou.online/  --tc=main.user:alim 2>testresults.log
```

or run **run_api_tests.sh** script as follow:
```
cd api_testing
bash run_api_tests.sh -d testsuite/
```
