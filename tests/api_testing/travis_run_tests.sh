if [ "$TRAVIS_EVENT_TYPE" == "cron" ] || [ "$TRAVIS_EVENT_TYPE" == "api" ]
then

    if [ "$environment" == "production" ]
    then
        url='https://itsyou.online/'
        user1_username=$production_user1_username
        user1_password=$production_user1_password
        user1_totp_secret=$production_user1_totp_secret
        user1_applicationid=$production_user1_applicationid
        user1_appsecret=$production_user1_appsecret
        user2_username=$production_user2_username
        user2_password=$production_user2_password
        user2_applicationid=$production_user2_applicationid
        user2_appsecret=$production_user2_appsecret
        
    elif [ "$environment" == "staging" ]
    then
        url='https://staging.itsyou.online/'
        user1_username=$staging_user1_username
        user1_password=$staging_user1_password
        user1_totp_secret=$staging_user1_totp_secret
        user1_applicationid=$staging_user1_applicationid
        user1_appsecret=$staging_user1_appsecret
        user2_username=$staging_user2_username
        user2_password=$staging_user2_password
        user2_applicationid=$staging_user2_applicationid
        user2_appsecret=$staging_user2_appsecret
    fi

    #run testsuite
    echo "Start tests on : ${url}" 
    nosetests-2.7 -v -s testsuite --tc-file=config.ini --tc=main.itsyouonline_url:$url --tc=main.validation_email:$validation_email --tc=main.validation_email_password:$validation_email_password --tc=main.user1_totp_secret:$user1_totp_secret --tc=main.user1_username:$user1_username --tc=main.user1_password:$user1_password --tc=main.user1_applicationid:$user1_applicationid --tc=main.user1_secret:$user1_appsecret --tc=main.user2_username:$user2_username --tc=main.user2_password:$user2_password --tc=main.user2_applicationid:$user2_applicationid --tc=main.user2_secret:$user2_appsecret  
    
else
    echo "Not a cron or api job"
fi