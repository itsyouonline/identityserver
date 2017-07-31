(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("UserHomeController", UserHomeController);


    UserHomeController.$inject = [
        '$q', '$rootScope', '$stateParams', '$location', '$window', '$filter', '$mdMedia', '$mdDialog', '$translate',
        'NotificationService', 'OrganizationService', 'UserService', 'UserDialogService'];

    function UserHomeController($q, $rootScope, $stateParams, $location, $window, $filter, $mdMedia, $mdDialog, $translate,
                                NotificationService, OrganizationService, UserService, UserDialogService) {
        var vm = this;
        vm.username = $rootScope.user;
        vm.notifications = {
            invitations: [],
            approvals: [],
            contractRequests: [],
            missingscopes: []
        };
        vm.notificationMessage = '';
        var authorizationArrayProperties = ['addresses', 'emailaddresses', 'phonenumbers', 'bankaccounts', 'digitalwallet', 'publicKeys'];
        var authorizationBoolProperties = ['facebook', 'github', 'name'];

        var TAB_YOU = 'you';
        var TAB_NOTIFICATIONS = 'notifications';
        var TAB_ORGANIZATIONS = 'organizations';
        var TAB_AUTHORIZATIONS = 'authorizations';
        var TAB_SETTINGS = 'settings';
        var TABS = [TAB_YOU, TAB_NOTIFICATIONS, TAB_ORGANIZATIONS,TAB_AUTHORIZATIONS, TAB_SETTINGS];

        vm.owner = [];
        vm.ownerTree = {};
        vm.member = [];
        vm.memberTree = {};
        vm.twoFAMethods = {};
        vm.user = {};

        vm.loaded = {};
        vm.selectedTabIndex = 0;
        vm.pendingCount = 0;

        UserDialogService.init(vm);

        /*vm.tabSelected = tabSelected;*/
        vm.pageSelected = pageSelected;
        vm.accept = accept;
        vm.reject = reject;
        vm.acceptorganizationinvite = acceptorganizationinvite;
        vm.rejectorganizationinvite = rejectorganizationinvite;
        vm.getPendingCount = getPendingCount;
        vm.showAvatarDialog = UserDialogService.avatar;
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
        vm.verifyEmailAddress = UserDialogService.verifyEmailAddress;
        vm.showAPIKeyDialog = showAPIKeyDialog;
        vm.showPublicKeyDetail = UserDialogService.publicKey;
        vm.createOrganization = UserDialogService.createOrganization;
        vm.showSetupAuthenticatorApplication = showSetupAuthenticatorApplication;
        vm.showExistingAuthenticatorApplication = showExistingAuthenticatorApplication;
        vm.removeAuthenticatorApplication = removeAuthenticatorApplication;
        vm.resolveMissingScopeClicked = resolveMissingScopeClicked;
        init();

        function init() {
            var index = TABS.indexOf($stateParams.tab);
            vm.selectedTabIndex = index === -1 ? 0 : index;
            loadUser().then(function () {
                loadVerifiedPhones();
                loadVerifiedEmails().then(loadNotifications);
            });
        }

        //redirect notification to right page
        function pageSelected(tabNum) {
            if(!(tabNum in TABS)) {
                return;
            }
            var path = '/' + TABS[tabNum];
            if(path !== $window.location.hash.replace('#', '')){
                $location.path(path);
            }
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
                            $translate(['user.controller.verifiedemails']).then(function(translations){
                                vm.notifications.security.push({
                                    tabIndex: 0,
                                    subject: 'verified_emails',
                                    msg: translations['user.controller.verifiedemails'],
                                    status: 'pending'
                                });
                            });
                        }
                        updatePendingNotificationsCount();
                        vm.loaded.notifications = true;
                    }
                );
        }

        function updatePendingNotificationsCount() {
            $translate(['user.controller.nonotifcations']).then(function(translations){
                vm.pendingCount = getPendingCount('all');
                vm.notificationMessage = vm.pendingCount ? '' : translations['user.controller.notifications'];
                $rootScope.notificationCount = vm.pendingCount;
            });
        }

        function loadOrganizations() {
            if (vm.loaded.organizations) {
                return;
            }
            OrganizationService
                .getUserOrganizations(vm.username)
                .then(
                    function (data) {
                        vm.owner = $filter('orderBy')(data.owner);
                        fillTree('owner');
                        vm.member = $filter('orderBy')(data.member);
                        fillTree('member');
                        vm.loaded.organizations = true;
                    }
                );
        }

        function fillTree(type) {
            var tree;
            var list;
            switch (type) {
                case 'owner':
                    tree = vm.ownerTree;
                    list = vm.owner;
                    break;
                case 'member':
                    tree = vm.memberTree;
                    list = vm.member;
                    break;
                default:
                    return;
            }

            // IMPORTANT: this depends on the list being semi-sorted: The parent organizations should come before their (direct) children.
            for (var i = 0; i < list.length; i++) {
                parseItemInTree(tree, list[i].split('.'), list[i]);
            }
        }

        function parseItemInTree(root, structure, target) {
            if (!root || !structure) {
                return;
            }
            for (var i = 0; i < structure.length; i++) {
                if (root[structure.slice(0, i + 1).join('.')]) {
                    parseItemInTree(root[structure.slice(0, i + 1).join('.')].children, structure.slice(i + 1), target);
                    return
                }
            }
            // childrenCollapsed is a property used by the tree directive to decide if the child elements should be shown.
            // Setting it here increases the reusability of the tree directive. The value is optional, and defaults to false.
            // Since it doesn't matter if there are any actual children, we can just set it on any node.
            // Set all values to true since the tree should only show root organizations by default.
            root[structure.join(".")] = { name: structure.join('.'), link: target, children: {}, childrenCollapsed: true};
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
                            angular.forEach(authorizationArrayProperties, function (prop) {
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
                count += vm.notifications.approvals ? vm.notifications.approvals.length : 0;
                count += vm.notifications.contractRequests ? vm.notifications.contractRequests.length : 0;
                count += vm.notifications.invitations ? vm.notifications.invitations.length : 0;
                count += vm.notifications.security ? vm.notifications.security.length : 0;
                vm.orgsWithInvitation = vm.notifications.invitations.map(function (invitation) {
                    return invitation.organization;
                });
                count += vm.notifications.missingscopes ? vm.notifications.missingscopes.filter(missingScopeFilter).length : 0;
                count += vm.notifications.organizationinvitations ? vm.notifications.organizationinvitations.length : 0;
                return count;
            } else {
                return obj ? obj.length : 0;
            }

            function missingScopeFilter(missingScope) {
                return vm.orgsWithInvitation.indexOf(missingScope.organization) === -1;
            }
        }

        function accept(event, invitation) {
            var missingScope = vm.notifications.missingscopes.filter(function (missingScope) {
                return missingScope.organization === invitation.organization;
            })[0];
            resolveMissingScope(event, missingScope).then(function () {
                NotificationService.accept(invitation, vm.username).then(function () {
                    invitation.status = 'accepted';
                    if (vm[invitation.role]) {
                        vm[invitation.role].push(invitation.organization);
                    }
                    updatePendingNotificationsCount();
                });
            });
        }

        function reject(invitation) {
            NotificationService
                .reject(invitation, vm.username)
                .then(function () {
                    invitation.status = 'rejected';
                    updatePendingNotificationsCount();
                });
        }

        function acceptorganizationinvite(event, invitation) {
            NotificationService.acceptorganizationinvite(invitation).then(function() {
                invitation.status = 'accepted';
                updatePendingNotificationsCount();
            });
        }

        function rejectorganizationinvite(invitation) {
            NotificationService.rejectorganizationinvite(invitation).then(function() {
                invitation.status = 'rejected';
                updatePendingNotificationsCount();
            })
        }

        function showAuthorizationDetailDialog(authorization, event, isNew) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));

            function authController($scope, $mdDialog, user, authorization, isNew) {
                angular.forEach(authorizationArrayProperties, function (prop) {
                    if (!authorization[prop]) {
                        authorization[prop] = [];
                    }

                });
                angular.forEach(authorizationBoolProperties, function (prop) {
                    if (authorization[prop] === undefined || authorization[prop] === null) {
                        authorization[prop] = false;
                    }
                });

                angular.forEach(authorization, function (auth, prop) {
                    if (Array.isArray(auth)) {
                        angular.forEach(auth, function (value) {
                            if (typeof value === 'object' && !value.reallabel) {
                                value.reallabel = vm.user[prop][0] ? vm.user[prop][0].label : '';
                            }
                        });
                    }
                });
                authorization.organizations = authorization.organizations || [];

                var ctrl = this;
                ctrl.user = user;
                ctrl.isNew = isNew;
                ctrl.delete = UserService.deleteAuthorization;
                $scope.update = update;
                ctrl.cancel = cancel;
                ctrl.remove = remove;
                $scope.requested = {
                    organizations: {}
                };
                authorization.organizations.map(function (org) {
                    $scope.requested.organizations[org] = true;
                });
                var originalAuthorization = JSON.parse(JSON.stringify(authorization));
                $scope.authorizations = authorization;

                function update(authorization) {
                    UserService.saveAuthorization($scope.authorizations)
                        .then(function (data) {
                            if (vm.authorizations) {
                                vm.authorizations.splice(vm.authorizations.indexOf(authorization), 1);
                                vm.authorizations.push(data);
                            }
                            $mdDialog.hide(data);
                        });
                }

                function cancel() {
                    if (vm.authorizations) {
                        var index = vm.authorizations.indexOf(authorization);
                        if (index !== 1) {
                            vm.authorizations.splice(index, 1);
                            vm.authorizations.push(originalAuthorization);
                        }
                    }
                    $mdDialog.cancel();
                }

                function remove() {
                    UserService.deleteAuthorization(authorization)
                        .then(function () {
                            if (vm.authorizations) {
                                vm.authorizations.splice(vm.authorizations.indexOf(authorization), 1);
                            }
                            $mdDialog.hide();
                        });
                }
            }

            return $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'user', 'authorization', 'isNew', authController],
                controllerAs: 'vm',
                templateUrl: 'components/user/views/authorizationDialog.html',
                targetEvent: event,
                fullscreen: useFullScreen,
                locals: {
                    user: vm.user,
                    authorization: authorization,
                    isNew: isNew
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
                        $translate(['user.controller.passwordupdated', 'user.controller.passwordchanged', 'user.controller.close']).then(function(translations) {
                            $mdDialog.hide();
                            $mdDialog.show(
                            $mdDialog.alert()
                                .clickOutsideToClose(true)
                                .title(translations['user.controller.passwordupdated'])
                                .textContent(translations['user.controller.passwordchanged'])
                                .ariaLabel(translations['user.controller.passwordupdated'])
                                .ok(translations['user.controller.close'])
                                .targetEvent(event)
                            );
                        })
                    }, function (response) {
                        if (response.status === 422) {
                            $scope.changepasswordform.currentPassword.$setValidity(response.data.error, false);
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

            function EditPasswordDialogController($mdDialog, user, updateName) {
                var ctrl = this;
                ctrl.save = save;
                ctrl.cancel = function () {
                    $mdDialog.cancel();
                };
                ctrl.firstname = user.firstname;
                ctrl.lastname = user.lastname;

                function save() {
                    updateName(user.username, ctrl.firstname, ctrl.lastname)
                        .then(function () {
                            $mdDialog.hide();
                            vm.user.firstname = ctrl.firstname;
                            vm.user.lastname = ctrl.lastname;
                        });
                }
            }

            $mdDialog.show({
                controller: ['$mdDialog', 'user', 'updateName', EditPasswordDialogController],
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
                            if (reason.status === 409) {
                                $scope.APIKeyForm.label.$setValidity('duplicate', false);
                            }
                        }
                    );
                }

                function updateAPIKey() {
                    UserService.updateAPIKey(username, ctrl.savedLabel, ctrl.label).then(
                        function () {
                            $mdDialog.hide({originalLabel: ctrl.savedLabel, newLabel: ctrl.label});
                        },
                        function (reason) {
                            if (reason.status === 409) {
                                $scope.APIKeyForm.label.$setValidity('duplicate', false);
                            }
                        }
                    );
                }

                function deleteAPIKey() {
                    UserService.deleteAPIKey(username, APIKey.label).then(
                        function () {
                            $mdDialog.hide({originalLabel: APIKey.label, newLabel: ""});
                        }
                    );
                }
            }
        }

        function showSetupAuthenticatorApplication(event) {
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'UserService', 'username', SetupAuthenticatorController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/setupTOTPDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                locals: {
                    username: vm.username
                }
            });

            function SetupAuthenticatorController($scope, $mdDialog, UserService, username) {
                var ctrl = this;
                ctrl.close = close;
                ctrl.submit = submit;
                ctrl.resetValidation = resetValidation;
                ctrl.username = username;
                ctrl.getQrCodeData = getQrCodeData;
                init();

                function init() {
                    UserService.getAuthenticatorSecret(vm.username)
                        .then(function (data) {
                            ctrl.totpsecret = data.totpsecret;
                            ctrl.totpissuer = encodeURIComponent(data.totpissuer);
                        });
                }


                function getQrCodeData() {
                    return 'otpauth://totp/' + ctrl.totpissuer + ':' + vm.username + '?secret=' + ctrl.totpsecret + '&issuer=' + ctrl.totpissuer;
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
                            }
                        });
                }

                function resetValidation() {
                    $scope.form.totpcode.$setValidity('invalid_totpcode', true);
                }
            }
        }

        function showExistingAuthenticatorApplication(event) {
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'UserService', 'username', ExistingAuthenticatorController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/viewTOTPDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                parent: angular.element(document.body),
                clickOutsideToClose: true,
                locals: {
                    username: vm.username
                }
            });

            function ExistingAuthenticatorController($scope, $mdDialog, UserService, username) {
                var ctrl = this;
                ctrl.close = close;
                ctrl.username = username;
                ctrl.getQrCodeData = getQrCodeData;
                init();

                function init() {
                    UserService.getAuthenticatorSecret(vm.username)
                        .then(function (data) {
                            ctrl.totpsecret = data.totpsecret;
                            ctrl.totpissuer = encodeURIComponent(data.totpissuer);
                        });
                }


                function getQrCodeData() {
                    return 'otpauth://totp/' + ctrl.totpissuer + ':' + vm.username + '?secret=' + ctrl.totpsecret + '&issuer=' + ctrl.totpissuer;
                }

                function close() {
                    $mdDialog.cancel();
                }
            }
        }

        function removeAuthenticatorApplication(event) {
            var hasConfirmedPhones = vm.user.phonenumbers.filter(function (phone) {
                    return phone.verified;
                }).length !== 0;
            if (!hasConfirmedPhones) {
                $translate(['user.controller.cantremoveauthapp', 'user.controller.cantremoveauthappmsg', 'ok']).then(function(translations){
                    $mdDialog.show(
                        $mdDialog.alert()
                            .clickOutsideToClose(true)
                            .title(translations['user.controller.cantremoveauthapp'])
                            .htmlContent(translations['user.controller.cantremoveauthappmsg'])
                            .ariaLabel(translations['user.controller.cantremoveauthapp'])
                            .ok(translations['ok'])
                            .targetEvent(event)
                    );
                });
                return;
            }
            $translate(['user.controller.removeauthenticator', 'user.controller.confirmremoveauthenticator', 'user.controller.yes', 'user.controller.no']).then(function(translations){
                var confirm = $mdDialog.confirm()
                    .title(translations['user.controller.removeauthenticator'])
                    .textContent(translations['user.controller.confirmremoveauthenticator'])
                    .ariaLabel(translations['user.controller.removeauthenticator'])
                    .targetEvent(event)
                    .ok(translations['user.controller.yes'])
                    .cancel(translations['user.controller.no']);
                $mdDialog.show(confirm).then(function () {
                    UserService.removeAuthenticator(vm.username)
                        .then(function () {
                            vm.twoFAMethods.totp = false;
                        });
                });
            });
        }

        function resolveMissingScopeClicked(event, missingScope) {
            resolveMissingScope(event, missingScope).then(updated);
            function updated() {
                vm.loaded.notifications = false;
                loadNotifications();
            }
        }

        function resolveMissingScope(event, missingScope) {
            return $q(promise);

            function promise(resolve, reject) {
                if (!missingScope) {
                    resolve();
                    return;
                }
                UserService.getAuthorization(vm.username, missingScope.organization).then(requestSuccess, requestFailed);
                function requestFailed(response) {
                    if (response.status === 404) {
                        var authorization = {
                            username: vm.username,
                            grantedTo: missingScope.organization
                        };
                        showDialog(authorization, true).then(resolve, reject);
                    } else {
                        reject(response);
                        $window.location.href = 'error' + response.status;
                    }
                }

                function requestSuccess(authorization) {
                    showDialog(authorization, false).then(resolve, reject);
                }

                function showDialog(authorization, isNew) {
                    // mostly copied from authorizeController -> parseScopes
                    var listAuthorizations = {
                        'address': 'addresses',
                        'email': 'emailaddresses',
                        'phone': 'phonenumbers',
                        'bankaccount': 'bankaccounts',
                        'publickey': 'publicKeys',
                        'avatar': 'avatars'
                    };
                    angular.forEach(missingScope.scopes, function (scope) {
                        var splitPermission = scope.split(':');
                        if (!splitPermission.length > 1) {
                            return;
                        }
                        // Empty label -> 'main'
                        var permissionLabel = splitPermission.length > 2 && splitPermission[2] ? splitPermission[2] : 'main';
                        var auth = {
                            requestedlabel: permissionLabel,
                            reallabel: ''
                        };
                        var userScope = splitPermission[1];
                        var listScope = listAuthorizations[userScope];
                        if (listScope) {
                            auth.reallabel = vm.user[listScope].length ? vm.user[listScope][0].label : '';
                            if (!authorization[listScope]) {
                                authorization[listScope] = [];
                            }
                            authorization[listScope].push(auth);
                        }
                        else if (scope === 'user:name') {
                            authorization.name = true;
                        }
                        else if (scope.startsWith('user:memberof:')) {
                            authorization.organizations.push(permissionLabel);
                        }
                        else if (scope.startsWith('user:digitalwalletaddress:')) {
                            auth.reallabel = vm.user.digitalwallet.length ? vm.user.digitalwallet[0].label : '';
                            auth.currency = splitPermission.length === 4 ? splitPermission[3] : '';
                            authorization.digitalwallet.push(auth);
                        }
                        else if (scope === 'user:github') {
                            authorization.github = true;
                        }
                        else if (scope === 'user:facebook') {
                            authorization.facebook = true;
                        }
                        else if (scope === 'user:keystore') {
                            authorizations.keystore = true;
                        }
                        else if (scope === 'user:see') {
                            authorization.see = true;
                        }
                        else if (scope.startsWith('user:validated:')){
                            switch (splitPermission[2]) {
                                case 'email':
                                    auth.reallabel = vm.user['emailaddresses'].length ? vm.user['emailaddresses'][0].label : '';
                                    $scope.authorizations['emailaddresses'].push(auth);
                                  break;
                                case 'phone':
                                    auth.reallabel = vm.user['phonenumbers'].length ? vm.user['phonenumbers'][0].label : '';
                                    $scope.authorizations['phonenumbers'].push(auth);
                                  break;
                            }
                        }
                    });
                    return showAuthorizationDetailDialog(authorization, event, isNew);
                }
            }
        }
    }

})();
