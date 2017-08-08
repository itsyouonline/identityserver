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

        var queryParams = $location.search();
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requestedorganizations = [];
        vm.username = UserService.getUsername();

        vm.user = {};
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
                        auth.reallabel = vm.user[listScope].length ? vm.user[listScope][0].label : '';
                        $scope.authorizations[listScope].push(auth);
                    }
                    else if (scope === 'user:name') {
                        $scope.authorizations.name = true;
                    }
                    else if (scope.startsWith('user:memberof:')) {
                        $scope.requested.organizations[permissionLabel] = true;
                    }
                    else if (scope.startsWith('user:digitalwalletaddress:')) {
                        auth.reallabel = vm.user.digitalwallet.length ? vm.user.digitalwallet[0].label : '';
                        auth.currency = splitPermission.length === 4 ? splitPermission[3] : '';
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
                                auth.reallabel = vm.user['emailaddresses'].length ? vm.user['emailaddresses'][0].label : '';
                                $scope.authorizations['validatedemailaddresses'].push(auth);
                              break;
                            case 'phone':
                                auth.reallabel = vm.user['phonenumbers'].length ? vm.user['phonenumbers'][0].label : '';
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

        function getInitialData() {
            if (Object.keys($scope.authorizations.ownerof.emailaddresses).length) {
                getVerifiedEmails();
            }
            if (Object.keys($scope.authorizations.validatedemailaddresses).length) {
                getValidatedEmails()
            }
            if (Object.keys($scope.authorizations.validatedphonenumbers).length) {
                getValidatedPhones();
            }

        }

        function getVerifiedEmails() {
            UserService.getVerifiedEmailAddresses(vm.username).then(function (confirmedEmails) {
                vm.verifiedEmails = confirmedEmails.map(function (mail) {
                    return mail.emailaddress;
                });
            });
        }

        function getValidatedEmails() {
            UserService.getVerifiedEmailAddresses(vm.username).then(function (confirmedEmails) {
                vm.validatedEmails = confirmedEmails.map(function (mail) {
                    return mail;
                });
            });
        }

        function getValidatedPhones() {
            UserService.getVerifiedPhones(vm.username).then(function (confirmedPhones) {
                vm.validatedPhones = confirmedPhones.map(function (phone) {
                    return phone;
                });
            });
        }

        function hasAllRequiredEmails() {
            for (var i = 0; i < $scope.authorizations.ownerof.emailaddresses.length; i++) {
                var email = $scope.authorizations.ownerof.emailaddresses[i];
                if (vm.verifiedEmails && !vm.verifiedEmails.includes(email)) {
                    return email;
                }
            }
            return true;
        }

        function checkOwner() {
            return $q(function (resolve, reject) {
                if ($scope.authorizations.ownerof.emailaddresses.length === 0 || hasAllRequiredEmails() === true) {
                    resolve();
                } else {
                    UserService.getVerifiedEmailAddresses(vm.username).then(function (confirmedEmails) {
                        vm.verifiedEmails = confirmedEmails.map(function (mail) {
                            return mail.emailaddress;
                        });
                        var requiredEmail = hasAllRequiredEmails();
                        if (requiredEmail === true) {
                            resolve();
                        } else {
                            reject({key: 'please_verify_email_x', values: {email: requiredEmail}});
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

        function selectDefault(fx, event, auth, property) {
            fx(event).then(function (data) {
                auth.reallabel = data.data.label;

            }, function () {
                // Select first possible value, else 'None'
                auth.reallabel = vm.user[property][0] ? vm.user[property][0].label : '';
            });
        }
    }
})();
