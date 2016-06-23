(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("UserHomeController", UserHomeController);


    UserHomeController.$inject = [
        '$q', '$rootScope', '$routeParams', '$location', '$window', '$mdToast', '$mdMedia', '$mdDialog',
        'NotificationService', 'OrganizationService', 'UserService', 'UserDialogService'];

    function UserHomeController($q, $rootScope, $routeParams, $location, $window, $mdToast, $mdMedia, $mdDialog,
                                NotificationService, OrganizationService, UserService, UserDialogService) {
        var vm = this;
        vm.username = $rootScope.user;
        vm.notifications = {
            invitations: [],
            approvals: [],
            contractRequests: []
        };
        vm.notificationMessage = '';
        var arrayProperties = ['addresses', 'emailaddresses', 'phonenumbers', 'bankaccounts', 'digitalwallet'];

        vm.owner = [];
        vm.member = [];
        vm.twoFAMethods = {};
        vm.user = {};

        vm.loaded = {};
        vm.selectedTabIndex = 0;
        vm.pendingCount = 0;

        UserDialogService.init(vm);

        vm.checkSelected = checkSelected;
        vm.tabSelected = tabSelected;
        vm.accept = accept;
        vm.reject = reject;
        vm.getPendingCount = getPendingCount;
        vm.showEmailDetailDialog = UserDialogService.emailDetail;
        vm.showPhonenumberDetailDialog = UserDialogService.phonenumberDetail;
        vm.showAddressDetailDialog = UserDialogService.addressDetail;
        vm.showAddressDetailDialog = UserDialogService.addressDetail;
        vm.showBankAccountDialog = UserDialogService.bankAccount;
        vm.showFacebookDialog = UserDialogService.facebook;
        vm.showGithubDialog = UserDialogService.github;
        vm.addFacebookAccount = UserDialogService.addFacebook;
        vm.addGithubAccount = UserDialogService.addGithub;
        vm.showDigitalWalletAddressDetail = UserDialogService.digitalWalletAddressDetail;
        vm.loadNotifications = loadNotifications;
        vm.loadOrganizations = loadOrganizations;
        vm.loadUser = loadUser;
        vm.loadAuthorizations = loadAuthorizations;
        vm.loadVerifiedPhones = loadVerifiedPhones;
        vm.loadSettings = loadSettings;
        vm.showAuthorizationDetailDialog = showAuthorizationDetailDialog;
        vm.showChangePasswordDialog = showChangePasswordDialog;
        vm.showEditNameDialog = showEditNameDialog;
        vm.verifyPhone = UserDialogService.verifyPhone;
        vm.verifyEmailAddress = verifyEmailAddress;
        vm.showAPIKeyDialog = showAPIKeyDialog;
        vm.createOrganization = UserDialogService.createOrganization;
        vm.showSetupAuthenticatorApplication = showSetupAuthenticatorApplication;
        vm.removeAuthenticatorApplication = removeAuthenticatorApplication;
        init();

        function init() {
            vm.selectedTabIndex = parseInt($routeParams.tab) || 0;
            loadUser()
                .then(function () {
                    loadVerifiedPhones();
                    loadVerifiedEmails()
                        .then(function () {
                            loadNotifications();
                        });
                });
        }

        function tabSelected(fx) {
            fx();
            $location.path('/home/' + vm.selectedTabIndex, false);
        }

        function loadNotifications() {
            if (vm.loaded.notifications) {
                return;
            }
            NotificationService
                .get(vm.username)
                .then(
                    function (data) {
                        vm.notifications = data;
                        vm.notifications.security = [];
                        var hasVerifiedEmail = vm.user.emailaddresses.filter(function (email) {
                                return email.verified;
                            }).length > 0;
                        if (!hasVerifiedEmail) {
                            vm.notifications.security.push({
                                tabIndex: 0,
                                subject: 'verified_emails',
                                msg: 'You do not have any verified email addresses. You will not be able to recover your account if you do not verify at least one of them.',
                                status: 'pending'
                            });
                        }
                        vm.pendingCount = getPendingCount('all');
                        vm.notificationMessage = vm.pendingCount ? '' : 'No unhandled notifications';
                        vm.loaded.notifications = true;
                        $rootScope.openRequests = vm.pendingCount;

                    }
                );
        }

        function loadOrganizations() {
            if (vm.loaded.organizations) {
                return;
            }
            OrganizationService
                .getUserOrganizations(vm.username)
                .then(
                    function (data) {
                        vm.owner = data.owner;
                        vm.member = data.member;
                        vm.loaded.organizations = true;
                    }
                );
        }

        function loadAuthorizations() {
            if (vm.loaded.authorizations) {
                return;
            }
            UserService.getAuthorizations(vm.username)
                .then(
                    function (data) {
                        vm.authorizations = data;
                        vm.loaded.authorizations = true;
                    }
                );
        }

        function loadUser() {
            return $q(function (resolve, reject) {
                if (vm.loaded.user) {
                    return;
                }
                UserService
                    .get(vm.username)
                    .then(
                        function (data) {
                            angular.forEach(arrayProperties, function (prop) {
                                if (!data[prop]) {
                                    data[prop] = [];
                                }
                            });
                            vm.user = data;
                            vm.loaded.user = true;
                            resolve(data);
                        }, reject
                    );
            });
        }

        function loadVerifiedPhones() {
            if (vm.loaded.verifiedPhones) {
                return;
            }
            UserService
                .getVerifiedPhones(vm.username)
                .then(function (confirmedPhones) {
                    confirmedPhones.map(function (p) {
                        findByLabel('phonenumbers', p.label).verified = true;
                    });
                    vm.loaded.verifiedPhones = true;
                });
        }

        function loadVerifiedEmails() {
            return $q(function (resolve, reject) {
                if (vm.loaded.verifiedEmails) {
                    return;
                }
                UserService
                    .getVerifiedEmailAddresses(vm.username)
                    .then(function (confirmedEmails) {
                        confirmedEmails.map(function (p) {
                            findByLabel('emailaddresses', p.label).verified = true;
                        });
                        vm.loaded.verifiedEmails = true;
                        resolve(confirmedEmails);
                    }, reject);
            });
        }

        function findByLabel(property, label) {
            return vm.user[property].filter(function (val) {
                return val.label === label;
            })[0];
        }

        function loadSettings() {
            if (vm.loaded.APIKeys) {
                return;
            }
            UserService
                .getAPIKeys(vm.username)
                .then(function (data) {
                    vm.APIKeys = data;
                    vm.loaded.APIKeys = true;
                });
            UserService
                .getTwoFAMethods(vm.username)
                .then(function (data) {
                    vm.twoFAMethods = data;
                });
        }

        function getPendingCount(obj) {
            var count = 0;
            if (obj === 'all') {
                count += vm.notifications.approvals.filter(pendingFilter).length;
                count += vm.notifications.contractRequests.filter(pendingFilter).length;
                count += vm.notifications.invitations.filter(pendingFilter).length;
                count += vm.notifications.security.length;
                return count;
            } else {
                return obj ? obj.filter(pendingFilter).length : 0;
            }
            function pendingFilter(prop) {
                return prop.status === 'pending';
            }
        }

        function checkSelected() {
            var selected = false;

            vm.notifications.invitations.forEach(function(invitation) {
                if (invitation.selected === true) {
                    selected = true;
                }
            });

            return selected;
        }

        function accept() {
            var requests = [];

            vm.notifications.invitations.forEach(function(invitation) {
                if (invitation.selected === true) {
                    requests.push(NotificationService.accept(invitation));
                }
            });

            $q
                .all(requests)
                .then(
                    function(responses) {
                        toast('Accepted ' + responses.length + ' invitations!');
                        vm.loaded.notifications = false;
                        loadNotifications();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function reject() {
            var requests = [];

            vm.notifications.invitations.forEach(function(invitation) {
                if (invitation.selected === true) {
                    requests.push(NotificationService.reject(invitation));
                }
            });

            $q
                .all(requests)
                .then(
                    function(responses) {
                        toast('Rejected ' + responses.length + ' invitations!');
                        vm.loaded.notifications = false;
                        loadNotifications();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function toast(message) {
            var toast = $mdToast
                .simple()
                .textContent(message)
                .hideDelay(2500)
                .position('top right');

            // Show toast!
            $mdToast.show(toast);
        }

        function showAuthorizationDetailDialog(authorization, event) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));

            function authController($scope, $mdDialog, user, authorization) {
                var ctrl = this;
                ctrl.user = user;
                $scope.delete = UserService.deleteAuthorization;
                $scope.update = update;
                $scope.cancel = cancel;
                $scope.remove = remove;
                $scope.requested = {};
                var originalAuthorization = JSON.parse(JSON.stringify(authorization));
                angular.forEach(authorization, function (value, key) {
                    if (Array.isArray(value)) {
                        $scope.requested[key] = {};
                        angular.forEach(value, function (v, i) {
                            if (typeof v !== 'object') {
                                $scope.requested[key][v] = true;
                            } else {
                                if (!Array.isArray($scope.requested[key])) {
                                    $scope.requested[key] = [];
                                }
                                $scope.requested[key].push(v.requestedlabel);
                            }
                        });
                    }
                    else {
                        $scope.requested[key] = value;

                    }
                });
                $scope.authorizations = authorization;

                function update(authorization) {
                    UserService.saveAuthorization($scope.authorizations)
                        .then(function (data) {
                            $mdDialog.cancel();
                            vm.authorizations.splice(vm.authorizations.indexOf(authorization), 1);
                            vm.authorizations.push(data);
                        });
                }

                function cancel() {
                    vm.authorizations.splice(vm.authorizations.indexOf(authorization), 1);
                    vm.authorizations.push(originalAuthorization);
                    $mdDialog.cancel();
                }

                function remove() {
                    UserService.deleteAuthorization(authorization)
                        .then(function (data) {
                            vm.authorizations.splice(vm.authorizations.indexOf(authorization), 1);
                            $mdDialog.cancel();
                        });
                }
            }

            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'user', 'authorization', authController],
                controllerAs: 'vm',
                templateUrl: 'components/user/views/authorizationDialog.html',
                targetEvent: event,
                fullscreen: useFullScreen,
                locals: {
                    user: vm.user,
                    authorization: authorization
                }
            });
        }

        function showChangePasswordDialog(event) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));

            function showPasswordDialogController($scope, $mdDialog, username, updatePassword) {
                var ctrl = this;
                ctrl.resetValidation = resetValidation;
                ctrl.updatePassword = updatepwd;
                ctrl.cancel = function () {
                    $mdDialog.cancel();
                };

                function resetValidation() {
                    $scope.changepasswordform.currentPassword.$setValidity('incorrect_password', true);
                    $scope.changepasswordform.currentPassword.$setValidity('invalid_password', true);
                }

                function updatepwd() {
                    updatePassword(username, ctrl.currentPassword, ctrl.newPassword).then(function () {
                        $mdDialog.hide();
                        $mdDialog.show(
                            $mdDialog.alert()
                                .clickOutsideToClose(true)
                                .title('Password updated')
                                .textContent('Your password has been changed.')
                                .ariaLabel('Password updated')
                                .ok('Close')
                                .targetEvent(event)
                        );
                    }, function (response) {
                        switch (response.status) {
                            case 422:
                                switch (response.data.error) {
                                    case 'incorrect_password':
                                        $scope.changepasswordform.currentPassword.$setValidity('incorrect_password', false);
                                        break;
                                    case 'invalid_password':
                                        $scope.changepasswordform.currentPassword.$setValidity('invalid_password', false);
                                        break;
                                }
                                break;
                            default:
                                $window.location.href = 'error' + response.status;
                                break;
                        }
                    });
                }
            }

            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'username', 'updatePassword', showPasswordDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/resetPasswordDialog.html',
                targetEvent: event,
                fullscreen: useFullScreen,
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                locals: {
                    username: vm.username,
                    updatePassword: UserService.updatePassword
                }
            });
        }

        function showEditNameDialog(event) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));

            function showPasswordDialogController($scope, $mdDialog, user, updateName) {
                var ctrl = this;
                ctrl.save = save;
                ctrl.cancel = function () {
                    $mdDialog.cancel();
                };
                ctrl.firstname = user.firstname;
                ctrl.lastname = user.lastname;

                function save() {
                    updateName(user.username, ctrl.firstname, ctrl.lastname)
                        .then(function (response) {
                            $mdDialog.hide();
                            vm.user.firstname = ctrl.firstname;
                            vm.user.lastname = ctrl.lastname;
                        });
                }
            }

            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'user', 'updateName', showPasswordDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/nameDialog.html',
                targetEvent: event,
                fullscreen: useFullScreen,
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                locals: {
                    user: vm.user,
                    updateName: UserService.updateName
                }
            });
        }

        function verifyEmailAddress(event, email) {
            UserService.sendEmailAddressVerification(vm.username, email.label)
                .then(function () {
                    $mdDialog.show(
                        $mdDialog.alert()
                            .clickOutsideToClose(true)
                            .title('Verification email sent')
                            .textContent('A verification email has been sent to ' + email.emailaddress + '.')
                            .ariaLabel('Verification email sent')
                            .ok('close')
                            .targetEvent(event)
                    );
                }, function () {
                    $mdDialog.show(
                        $mdDialog.alert()
                            .clickOutsideToClose(true)
                            .title('Error')
                            .textContent('Could not send verification email. Please try again later.')
                            .ariaLabel('Error while sending verification email')
                            .ok('close')
                            .targetEvent(event)
                    );
                });
        }

        function showAPIKeyDialog(event, APIKey) {
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'UserService', 'username', 'APIKey', APIKeyDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/APIKeyDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                clickOutsideToClose: true,
                locals: {
                    UserService: UserService,
                    username: vm.username,
                    APIKey: APIKey
                }
            })
                .then(
                    function (data) {
                        if (data.originalLabel != data.newLabel) {
                            if (data.originalLabel) {
                                var originalKey = getApiKey(data.originalLabel);
                                if (data.newLabel) {
                                    // update
                                    originalKey.label = data.newLabel;
                                }
                                else {
                                    // remove
                                    vm.APIKeys.splice(vm.APIKeys.indexOf(originalKey), 1);
                                }
                            }
                            else {
                                // add
                                vm.APIKeys.push(data.APIKey);
                            }
                        }
                    });

            function getApiKey(label) {
                return vm.APIKeys.filter(function (k) {
                    return k.label === label;
                })[0];
            }

            function APIKeyDialogController($scope, $mdDialog, UserService, username, APIKey) {
                var ctrl = this;
                ctrl.APIKey = APIKey || {secret: ""};
                ctrl.originalLabel = APIKey ? APIKey.label : null;
                ctrl.savedLabel = APIKey ? APIKey.label : null;
                ctrl.label = APIKey ? APIKey.label : null;

                ctrl.cancel = cancel;
                ctrl.create = createAPIKey;
                ctrl.update = updateAPIKey;
                ctrl.delete = deleteAPIKey;

                ctrl.modified = false;

                function cancel() {
                    if (ctrl.modified) {
                        $mdDialog.hide({originalLabel: ctrl.originalLabel, newLabel: ctrl.label, APIKey: ctrl.APIKey});
                    }
                    else {
                        $mdDialog.cancel();
                    }
                }

                function createAPIKey() {
                    ctrl.validationerrors = {};
                    UserService.createAPIKey(username, ctrl.label).then(
                        function (data) {
                            ctrl.modified = true;
                            ctrl.APIKey = data;
                            ctrl.savedLabel = data.label;
                        },
                        function (reason) {
                            if (reason.status == 409) {
                                $scope.APIKeyForm.label.$setValidity('duplicate', false);
                            }
                            else {
                                $window.location.href = "error" + reason.status;
                            }
                        }
                    );
                }

                function updateAPIKey() {
                    UserService.updateAPIKey(username, ctrl.savedLabel, ctrl.label).then(
                        function (data) {
                            $mdDialog.hide({originalLabel: ctrl.savedLabel, newLabel: ctrl.label});
                        },
                        function (reason) {
                            if (reason.status == 409) {
                                $scope.APIKeyForm.label.$setValidity('duplicate', false);
                            }
                            else {
                                $window.location.href = "error" + reason.status;
                            }
                        }
                    );
                }

                function deleteAPIKey() {
                    UserService.deleteAPIKey(username, APIKey.label).then(
                        function (data) {
                            $mdDialog.hide({originalLabel: APIKey.label, newLabel: ""});
                        },
                        function (reason) {
                            $window.location.href = "error" + reason.status;
                        }
                    );
                }
            }
        }

        function showSetupAuthenticatorApplication(event) {
            $mdDialog.show({
                controller: ['$scope', '$window', '$mdDialog', 'UserService', SetupAuthenticatorController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/setupTOTPDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                parent: angular.element(document.body),
                clickOutsideToClose: true
            });

            function SetupAuthenticatorController($scope, $window, $mdDialog, UserService) {
                var ctrl = this;
                ctrl.close = close;
                ctrl.submit = submit;
                ctrl.resetValidation = resetValidation;
                vm.config = {};
                init();

                function init() {
                    UserService.getAuthenticatorSecret(vm.username)
                        .then(function (data) {
                            ctrl.totpsecret = data.totpsecret;
                        });
                }

                function close() {
                    $mdDialog.cancel();
                }

                function submit() {
                    UserService.setAuthenticator(vm.username, ctrl.totpsecret, ctrl.totpcode)
                        .then(function () {
                            vm.twoFAMethods.totp = true;
                            $mdDialog.hide();
                        }, function (response) {
                            if (response.status === 422) {
                                $scope.form.totpcode.$setValidity('invalid_totpcode', false);
                            } else {
                                $window.location.href = 'error' + response.status;
                            }
                        });
                }

                function resetValidation() {
                    $scope.form.totpcode.$setValidity('invalid_totpcode', true);
                }
            }
        }

        function removeAuthenticatorApplication(event) {
            var hasConfirmedPhones = vm.user.phonenumbers.filter(function (phone) {
                    return phone.verified;
                }).length !== 0;
            if (!hasConfirmedPhones) {
                var msg = 'You cannot remove your authenticator application because this is your last two-factor' +
                    ' authentication method.<br />Add a phone number and verify it to be able to remove your authenticator application.';
                $mdDialog.show(
                    $mdDialog.alert()
                        .clickOutsideToClose(true)
                        .title('Cannot remove authenticator')
                        .htmlContent(msg)
                        .ariaLabel('Cannot remove authenticator')
                        .ok('Ok')
                        .targetEvent(event)
                );
                return;
            }
            var confirm = $mdDialog.confirm()
                .title('Unauthorize authenticator')
                .textContent('Are you sure you want to unauthorize your authenticator application?')
                .ariaLabel('Unauthorize authenticator')
                .targetEvent(event)
                .ok('Yes')
                .cancel('No');
            $mdDialog.show(confirm).then(function () {
                UserService.removeAuthenticator(vm.username)
                    .then(function () {
                        vm.twoFAMethods.totp = false;
                    });
            });
        }
    }

})();
