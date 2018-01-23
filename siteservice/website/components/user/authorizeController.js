(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("AuthorizeController", AuthorizeController);


    AuthorizeController.$inject = ['$scope', '$location', '$window', '$q', '$translate',
        'UserService', 'UserDialogService', 'NotificationService'];

    function AuthorizeController($scope, $location, $window, $q, $translate,
                                 UserService, UserDialogService, NotificationService) {
        var vm = this;
        vm.isAuthorizeController = true;

        var queryParams = $location.search();
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requestedorganizations = [];
        vm.username = UserService.getUsername();

        vm.user = {};
        vm.loaded = {};
        vm.pendingNotifications = [];
        vm.pendingOrganizationInvites = {};

        UserDialogService.init(vm);
        vm.showAvatarDialog = addAvatar;
        vm.showEmailDialog = addEmail;
        vm.showPhonenumberDialog = addPhone;      
        vm.showAddressDialog = addAddress;
        vm.showBankAccountDialog = bank;
        vm.showPublicKeyDialog = addPublicKey;
        vm.verifyEmail = verifyEmail;
        vm.submit = submit;
        vm.showDigitalWalletAddressDialog = digitalWalletAddress;
        vm.isEmailVerified = isEmailVerified;
        vm.isPhoneVerified = isPhoneVerified;
        vm.verifyEmailByLabel = verifyEmailByLabel;
        vm.verifyPhoneByLabel = verifyPhoneByLabel;
        var properties = ['avatars', 'addresses', 'emailaddresses', 'phonenumbers', 'bankaccounts', 'digitalwallet', 'publicKeys', 'validatedemailaddresses', 'validatedphonenumbers'];
        $scope.requested = {
            organizations: {}
        };
        $scope.authorizations = {
            ownerof: {
                emailaddresses: []
            }
        };
        angular.forEach(properties, function (prop) {
            $scope.authorizations[prop] = [];
        });

        $scope.update = update;
        $scope.isNew = true;

        activate();

        function activate() {
            fetch();
        }

        function fetch() {

            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.user = data;
                        parseScopes();
                        loadVerifiedPhones();
                        loadVerifiedEmails();
                        getInitialData();
                    });
            NotificationService
                .get(vm.username)
                .then(
                    function (data) {
                        vm.pendingNotifications = data.invitations.filter(function (invitation) {
                            return invitation.status === 'pending';
                        });
                        angular.forEach(vm.pendingNotifications, function (invite) {
                            vm.pendingOrganizationInvites[invite.organization] = true;
                        });
                    }
                );

        }

        function parseScopes() {
            if (vm.requestedScopes) {
                var listAuthorizations = {
                    'address': 'addresses',
                    'email': 'emailaddresses',
                    'phone': 'phonenumbers',
                    'bankaccount': 'bankaccounts',
                    'publickey': 'publicKeys',
                    'avatar': 'avatars'
                };
                var scopes = vm.requestedScopes.split(',');
                // Filter duplicated scopes
                scopes = scopes.filter(function (item, pos, self) {
                    return self.indexOf(item) === pos;
                });
                angular.forEach(scopes, function (scope) {
                    var splitPermission = scope.split(':');
                    if (!splitPermission.length > 1) {
                        return;
                    }
                    // Empty label -> 'main'
                    var userScope = splitPermission[1];
                    var permissionLabel = splitPermission.length > 2 && splitPermission[2] ? splitPermission[2] : 'main';
                    // last part is always the read or write permission
                    var readWrite = splitPermission[splitPermission.length - 1];
                    if (!['read', 'write'].includes(readWrite)) {
                        readWrite = null;
                    }
                    var auth = {
                        requestedlabel: permissionLabel,
                        reallabel: '',
                        scope: readWrite
                    };
                    var listScope = listAuthorizations[userScope];
                    if (listScope) {
                        auth.reallabel = vm.user[listScope].length ? vm.user[listScope][0].label : null;
                        $scope.authorizations[listScope].push(auth);
                    }
                    else if (scope === 'user:name') {
                        $scope.authorizations.name = true;
                    }
                    else if (scope.startsWith('user:memberof:')) {
                        $scope.requested.organizations[permissionLabel] = true;
                    }
                    else if (scope.startsWith('user:digitalwalletaddress:')) {
                        auth.reallabel = vm.user.digitalwallet.length ? vm.user.digitalwallet[0].label : null;
                        auth.currency = splitPermission.length === 4 ? splitPermission[3] : null;
                        $scope.authorizations.digitalwallet.push(auth);
                    }
                    else if (scope === 'user:github') {
                        $scope.authorizations.github = true;
                    }
                    else if (scope === 'user:facebook') {
                        $scope.authorizations.facebook = true;
                    }
                    else if (scope === 'user:keystore') {
                        $scope.authorizations.keystore = true;
                    }
                    else if (scope === 'user:see') {
                        $scope.authorizations.see = true;
                    }
                    else if (scope.startsWith('user:validated:')){
                        permissionLabel = splitPermission.length > 3 && splitPermission[3] ? splitPermission[3] : 'main';
                        auth.requestedlabel = permissionLabel;
                        switch (splitPermission[2]) {
                            case 'email':
                                auth.reallabel = vm.user['emailaddresses'].length ? vm.user['emailaddresses'][0].label : null;
                                $scope.authorizations['validatedemailaddresses'].push(auth);
                              break;
                            case 'phone':
                                auth.reallabel = vm.user['phonenumbers'].length ? vm.user['phonenumbers'][0].label : null;
                                $scope.authorizations['validatedphonenumbers'].push(auth);
                              break;
                        }
                    }
                    else if (scope.startsWith('user:ownerof')) {
                        switch (splitPermission[2]) {
                            case 'email':
                                var emailAddress = splitPermission[3];
                                if (emailAddress && !$scope.authorizations.ownerof.emailaddresses.includes(emailAddress)) {
                                    $scope.authorizations.ownerof.emailaddresses.push(emailAddress);
                                }
                                break;
                        }
                    }
                });
            }
        }

        function loadVerifiedPhones() {
            return $q(function (resolve, reject) {
                if (vm.isLoadingVerifiedPhones) {
                    return;
                }
                vm.isLoadingVerifiedPhones = true;
                UserService
                    .getVerifiedPhones(true)
                    .then(function (confirmedPhones) {
                        angular.forEach(vm.validatedphonenumbers, function (vp) {
                            vp.verified = undefined;
                        });
                        confirmedPhones.map(function (p) {
                            findByLabel('phonenumbers', p.label).verified = true;
                        });
                        vm.loaded.verifiedPhones = true;
                        resolve(confirmedPhones);
                    }, reject);
                    vm.isLoadingVerifiedPhones = false;
            });
        }

        function loadVerifiedEmails() {
            return $q(function (resolve, reject) {
                if (vm.isLoadingVerifiedEmails) {
                    return;
                }
                vm.isLoadingVerifiedEmails = true;
                UserService
                    .getVerifiedEmailAddresses(true)
                    .then(function (confirmedEmails) {
                        angular.forEach(vm.validatedemailaddresses, function (ve) {
                            ve.verified = undefined;
                        });
                        confirmedEmails.map(function (p) {
                            findByLabel('emailaddresses', p.label).verified = true;
                        });
                        vm.loaded.verifiedEmails = true;
                        resolve(confirmedEmails);
                    }, reject);
                vm.isLoadingVerifiedEmails = false;
            });
        }

        function findByLabel(property, label) {
            return vm.user[property].filter(function (val) {
                return val.label === label;
            })[0];
        }

        function getInitialData() {
            if (Object.keys($scope.authorizations.ownerof.emailaddresses).length) {
                getVerifiedEmails();
            }


        }

        function getMissingRequiredEmails() {
            return $scope.authorizations.ownerof.emailaddresses.filter(function (email) {
                return vm.verifiedEmails && !vm.verifiedEmails.includes(email);
            });
        }

        function checkOwner() {
            return $q(function (resolve, reject) {
                if ($scope.authorizations.ownerof.emailaddresses.length === 0 || getMissingRequiredEmails().length === 0) {
                    resolve();
                } else {
                    UserService.getVerifiedEmailAddresses(true).then(function (confirmedEmails) {
                        setConfirmedEmails(confirmedEmails);
                        var missingRequiredEmails = getMissingRequiredEmails();
                        if (missingRequiredEmails.length === 0) {
                            resolve();
                        } else {
                            reject({key: 'please_verify_email_x', values: {email: missingRequiredEmails[0]}});
                        }
                    });
                }
            });
        }

        // called by the authorizationDetailsDirective
        function update(event) {
            // validation
            checkOwner().then(function () {
                $scope.authorizations.username = vm.username;
                $scope.authorizations.grantedTo = vm.requestingorganization;
                UserService
                    .saveAuthorization($scope.authorizations)
                    .then(
                        function () {
                            var u = URI($location.absUrl());
                            var endpoint = queryParams["endpoint"];
                            delete queryParams.endpoint;
                            u.pathname(endpoint);
                            u.search(queryParams);
                            $window.location.href = u.toString();
                        }
                    );
            }, function (translation) {
                $translate(['error', 'close', translation.key], translation.values).then(function (translations) {
                    UserDialogService.showSimpleDialog(translations[translation.key], translations['error'], translations['close'], event);
                });
            });
        }

        function addAvatar(event, auth) {
            selectDefault(UserDialogService.avatar, event, auth, 'avatars')
        }

        function addEmail(event, auth) {
            selectDefault(UserDialogService.emailDetail, event, auth, 'emailaddresses');
        }

        function addPhone(event, auth) {
            selectDefault(UserDialogService.phonenumberDetail, event, auth, 'phonenumbers');
        }

        function addAddress(event, auth) {
            selectDefault(UserDialogService.addressDetail, event, auth, 'addresses');
        }

        function bank(event, auth) {
            selectDefault(UserDialogService.bankAccount, event, auth, 'bankaccounts');
        }

        function addPublicKey(event, auth) {
            selectDefault(UserDialogService.publicKey, event, auth, 'publicKeys');
        }

        function verifyEmail(event, email) {
            var userEmail = vm.user.emailaddresses.filter(function (e) {
                return e.emailaddress === email;
            })[0];
            if (userEmail) {
                verify(userEmail);
            } else {
                var data = {
                    label: email.replace('@', ' - ').replace('+', ' ').split('.')[0],
                    emailaddress: email
                };
                UserService.registerNewEmailAddress(vm.username, data).then(function (newEmail) {
                    vm.user.emailaddresses.push(newEmail);
                    verify(newEmail);
                });
            }

            function verify(email) {
                UserDialogService.verifyEmailAddress(event, email).then(getVerifiedEmails);
            }
        }

        function submit(event) {
            if ($scope.authorizeform.$invalid) {
                return;
            }

            var requests = [];
            if ($scope.authorizations['validatedemailaddresses'] && $scope.authorizations['validatedemailaddresses'].length > 0) {
                requests.push(loadVerifiedEmails());
            }
            if ($scope.authorizations['validatedphonenumbers'] && $scope.authorizations['validatedphonenumbers'].length > 0) {
                requests.push(loadVerifiedPhones());
            }

            
            // Remove the messages that the email address and phone number needs to be validated first
            vm.emailNotValidated = {};
            vm.phoneNotValidated = {};
            
            $q.all(requests)
                .then(function (results) {
               
                    // If we have validated email address scopes check if given values are all verified
                    // If not, set a value to notify the UI and return
                    if ($scope.authorizations['validatedemailaddresses'] && $scope.authorizations['validatedemailaddresses'].length > 0) {
                        angular.forEach($scope.authorizations['validatedemailaddresses'], function(validatedemailrequest) {
                            if (!isEmailVerified(validatedemailrequest.reallabel)) {
                                var email = vm.user.emailaddresses.filter(function (e) {
                                    return e.label === validatedemailrequest.reallabel;
                                })[0];
                                vm.emailNotValidated = { email: email.emailaddress };
                                return;
                            }
                        });
                    }

                    if (vm.emailNotValidated.email) {
                        return;
                    }

                    // Likewise, check if the phone numbers are validated 
                    if ($scope.authorizations['validatedphonenumbers'] && $scope.authorizations['validatedphonenumbers'].length > 0) {
                        angular.forEach($scope.authorizations['validatedphonenumbers'], function(validatedphonerequest) {
                            if (!isPhoneVerified(validatedphonerequest.reallabel)) {
                                var phone = vm.user.phonenumbers.filter(function (p) {
                                    return p.label === validatedphonerequest.reallabel;
                                })[0];
                                vm.phoneNotValidated = { phone: phone.phonenumber };
                                return;
                            }
                        });
                    }

                    if (vm.phoneNotValidated.phone) {
                        return;
                    }

                    confirm(event);

                });
        
        }

        function confirm(event) {
            // accept all the invites first
            var requests = [];

            vm.pendingNotifications.forEach(function (invitation) {
                requests.push(NotificationService.accept(invitation, vm.username));
            });

            $q.all(requests)
                .then(function () {
                    $scope.save(event);
                });
        }

        function digitalWalletAddress(event, auth) {
            selectDefault(UserDialogService.digitalWalletAddressDetail, event, auth, 'digitalwallet');
        }

        function isEmailVerified(authlabel) {
            var email = vm.user.emailaddresses.filter(function (e) {
                return e.label === authlabel;
            })[0];
            return email && !!email.verified
        }

        function isPhoneVerified(authlabel) {
            var phone = vm.user.phonenumbers.filter(function (p) {
                return p.label === authlabel;
            })[0];
            return phone && !!phone.verified
        }

        function verifyEmailByLabel(event, authlabel) {
            var email = vm.user.emailaddresses.filter(function (e) {
                return e.label === authlabel;
            })[0];
            if (email) {
                UserDialogService.verifyEmailAddress(event, email).then(loadVerifiedEmails)
            }
        }

        function verifyPhoneByLabel(event, authlabel) {
            var phone = vm.user.phonenumbers.filter(function (p) {
                return p.label === authlabel;
            })[0];
            if (phone) {
                UserDialogService.verifyPhone(event, phone).then(loadVerifiedPhones)
            }
        }

        function selectDefault(fx, event, auth, property) {
            fx(event).then(function (data) {
                auth.reallabel = data.data.label;
            }, function () {
                auth.reallabel = vm.user[property][0] ? vm.user[property][0].label : '';
            });
        }
    }
})();
