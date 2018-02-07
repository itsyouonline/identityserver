/**
 * Created by lucas on 13/06/16.
 */
(function () {
    'use strict';

    angular
        .module('itsyouonline.user', [])
        .factory('UserDialogService', ['$window', '$q', '$interval', '$mdMedia', '$mdDialog', '$translate',
            'UserService', 'configService', UserDialogService]);

    function UserDialogService($window, $q, $interval, $mdMedia, $mdDialog, $translate, UserService, configService) {
        var vm;
        var genericDetailControllerParams = ['$scope', '$mdDialog', 'user', 'data',
            'createFunction', 'updateFunction', 'deleteFunction', 'preFilledLabels', GenericDetailDialogController];
        return {
            init: init,
            emailDetail: emailDetail,
            addressDetail: addressDetail,
            phonenumberDetail: phonenumberDetail,
            verifyPhone: verifyPhone,
            verifyEmailAddress: verifyEmailAddress,
            bankAccount: bankAccount,
            facebook: facebook,
            addFacebook: addFacebook,
            github: github,
            addGithub: addGithub,
            showSimpleDialog: showSimpleDialog,
            createOrganization: createOrganization,
            digitalWalletAddressDetail: digitalWalletAddressDetail,
            publicKey: publicKey,
            avatar: avatar
        };

        function init(scope) {
            vm = scope;
        }

        function doNothing() {
        }

        function findByLabel(property, label) {
            return vm.user[property].filter(function (val) {
                return val.label === label;
            })[0];
        }

        function emailDetail(ev, email) {

            return $q(function (resolve, reject) {
                var preFilledLabels = ["Work", "Personal", "Spam", "Other"];
                email = email || {};
                var originalEmail = JSON.parse(JSON.stringify(email));
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/emailaddressdialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        user: vm.user,
                        data: email,
                        createFunction: UserService.registerNewEmailAddress,
                        updateFunction: UserService.updateEmailAddress,
                        deleteFunction: UserService.deleteEmailAddress,
                        preFilledLabels: preFilledLabels
                    }
                })
                    .then(
                        function (data) {
                            if (data.fx === 'create') {
                                vm.user.emailaddresses.push(data.data);
                            } else if (data.fx === 'delete') {
                                vm.user.emailaddresses.splice(vm.user.emailaddresses.indexOf(email), 1);
                            }
                            resolve(data);
                        }, function (response) {
                            angular.forEach(originalEmail, function (value, key) {
                                email[key] = value;
                            });
                            reject(response);
                        });
            });
        }

        function phonenumberDetail(ev, phone) {
            var preFilledLabels = ["Home", "Work", "Mobile", "General", "Fax", "Pager", "Other"];
            return $q(function (res, rej) {
                phone = phone || {};
                var originalPhone = JSON.parse(JSON.stringify(phone));

                function deletePhoneNumber(username, label) {
                    return $q(function (resolve, reject) {
                        rollbackPhoneNumber();
                        UserService
                            .deletePhonenumber(username, phone.label, false)
                            .then(resolve, function (response) {
                                if (response.status === 409) {
                                    var errorMsg, title, dialog;
                                    if (response.data.error === 'cannot_delete_last_verified_phone_number') {
                                        errorMsg = 'You cannot delete your last verified phone number.';
                                        title = "Can't delete last verified phone number";
                                        showSimpleDialog(errorMsg, title, 'ok', ev);
                                    }
                                } else {
                                    reject();
                                    errorMsg = 'Could not delete phone number. Please try again later.';
                                    showSimpleDialog(errorMsg, 'Error', 'Ok', ev);
                                }
                            });
                    });
                }

                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/phonenumberdialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        user: vm.user,
                        data: phone,
                        createFunction: UserService.registerNewPhonenumber,
                        updateFunction: UserService.updatePhonenumber,
                        deleteFunction: deletePhoneNumber,
                        preFilledLabels: preFilledLabels
                    }
                })
                    .then(updatePhoneNumber, rollbackPhoneNumber);

                function rollbackPhoneNumber(response) {
                    angular.forEach(originalPhone, function (value, key) {
                        phone[key] = value;
                    });
                    rej(response);
                }

                function updatePhoneNumber(data) {
                    // no data is provided when dialog is closed because another dialog opened (in case a confirmation is asked)
                    if (data) {
                        var newPhone = data.data;
                        if (data.fx === 'delete') {
                            vm.user.phonenumbers.splice(vm.user.phonenumbers.indexOf(phone), 1);
                        }
                        else {
                            // Mark a phonenumber as verified if it's the same number as an already verified one.
                            // Executed when updating and adding
                            vm.user.phonenumbers.map(function (number) {
                                if (number.phonenumber === newPhone.phonenumber) {
                                    newPhone.verified = true;
                                }
                            });
                            if (data.fx === 'create') {
                                vm.user.phonenumbers.push(newPhone);
                            }
                            if (!newPhone.verified) {
                                verifyPhone(ev, newPhone);
                            }
                        }
                    }
                    res(data);
                }

            });
        }

        function verifyPhone(event, phone) {
            var promise = $q.defer();
            var interval;
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', '$interval', 'user', 'phone', verifyPhoneDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/verifyPhoneDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    user: vm.user,
                    phone: phone
                }
            }).finally(function () {
                $interval.cancel(interval);
            });

            function verifyPhoneDialogController($scope, $mdDialog, $interval, user, phone) {
                var ctrl = this;
                ctrl.label = phone.label;
                ctrl.phone = phone.phonenumber;
                ctrl.close = close;
                ctrl.submit = submit;
                ctrl.validationKey = '';
                ctrl.resetValidation = resetValidation;

                init();

                function init() {
                    UserService
                        .sendPhoneVerificationCode(vm.username, phone.label)
                        .then(function (responseData) {
                            ctrl.validationKey = responseData.validationkey;
                            interval = $interval(checkconfirmation, 1000);
                        }, function () {
                            $mdDialog.show(
                                $mdDialog.alert()
                                    .clickOutsideToClose(true)
                                    .title('Error')
                                    .textContent('Failed to send verification code. Please try again later.')
                                    .ariaLabel('Error while sending verification code')
                                    .ok('Close')
                                    .targetEvent(event)
                            );
                        });
                }

                function close() {
                    $mdDialog.cancel();
                }

                function checkconfirmation() {
                    UserService
                        .getVerifiedPhones(true)
                        .then(function success(confirmedPhones) {
                            var confirmed = confirmedPhones.filter(function (p) {
                                    return p.label === ctrl.label;
                                }).length !== 0;
                            if (confirmed) {
                                findByLabel('phonenumbers', ctrl.label).verified = true;
                                close();
                                promise.resolve();
                            }
                        }, function (response) {
                            promise.reject(response);
                        });
                }

                function submit() {
                    UserService
                        .verifyPhone(user.username, ctrl.label, ctrl.validationKey, ctrl.smscode)
                        .then(function () {
                            findByLabel('phonenumbers', ctrl.label).verified = true;
                            close();
                            promise.resolve()
                        }, function (response) {
                            if (response.status === 422) {
                                $scope.form.smscode.$setValidity('invalid_code', false);
                            }
                            promise.reject(response);
                        });
                }

                function resetValidation() {
                    $scope.form.smscode.$setValidity('invalid_code', true);
                }
            }
            return promise.promise;
        }

        function verifyEmailAddress(event, email) {
            return $q(function (res, rej) {
                UserService.sendEmailAddressVerification(vm.username, email.label)
                    .then(function () {
                        $translate(['user.controller.emailsent', 'user.controller.emailsentto', 'user.controller.close'], {email: email.emailaddress}).then(function (translations) {
                            $mdDialog.show(
                                $mdDialog.alert()
                                    .clickOutsideToClose(true)
                                    .title(translations['user.controller.emailsent'])
                                    .textContent(translations['user.controller.emailsentto'])
                                    .ariaLabel(translations['user.controller.emailsent'])
                                    .ok(translations['user.controller.close'])
                                    .targetEvent(event)
                            );
                        });
                        res();
                    }, function () {
                        $translate(['user.controller.error', 'user.controller.couldnotsend', 'user.controller.errorwhilesending', 'user.controller.close']).then(function (translations) {
                            $mdDialog.show(
                                $mdDialog.alert()
                                    .clickOutsideToClose(true)
                                    .title(translations['user.controller.error'])
                                    .textContent(translations['user.controller.couldnotsend'])
                                    .ariaLabel(translations['user.controller.errorwhilesending'])
                                    .ok(translations['user.controller.close'])
                                    .targetEvent(event)
                            );
                            rej();
                        });
                    });
            });
        }

        function addressDetail(ev, address) {
            return $q(function (resolve, reject) {
                var preFilledLabels = ["Home", "Work", "Other"];
                address = address || {};
                var originalAddress = JSON.parse(JSON.stringify(address));
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/addressdialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        user: vm.user,
                        data: address,
                        createFunction: UserService.registerNewAddress,
                        updateFunction: UserService.updateAddress,
                        deleteFunction: UserService.deleteAddress,
                        preFilledLabels: preFilledLabels
                    }
                })
                    .then(
                        function (data) {
                            if (data.fx === 'create') {
                                vm.user.addresses.push(data.data);
                            }
                            else if (data.fx === 'delete') {
                                // delete
                                vm.user.addresses.splice(vm.user.addresses.indexOf(address), 1);
                            }
                            resolve(data);
                        }, function (response) {
                            // Dialog closed without saving. Rollback.
                            angular.forEach(originalAddress, function (value, key) {
                                address[key] = value;
                            });
                            reject(response);
                        });
            });
        }

        function bankAccount(ev, bank) {
            var preFilledLabels = ["Personal", "Business", "Shared", "Other"];
            var originalBankAccount = JSON.parse(JSON.stringify(bank || {}));
            return $q(function (resolve, reject) {
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/bankAccountDialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        user: vm.user,
                        data: bank,
                        createFunction: UserService.registerNewBankAccount,
                        updateFunction: UserService.updateBankAccount,
                        deleteFunction: UserService.deleteBankAccount,
                        preFilledLabels: preFilledLabels
                    }
                })
                    .then(
                        function (data) {
                            if (data.fx === 'delete') {
                                vm.user.bankaccounts.splice(vm.user.bankaccounts.indexOf(bank), 1);
                            }
                            else if (data.fx === 'create') {
                                vm.user.bankaccounts.push(data.data);
                            }
                            resolve(data);
                        }, function (response) {
                            angular.forEach(originalBankAccount, function (value, key) {
                                bank[key] = value;
                            });
                            reject(response);
                        });
            });
        }

        function avatar(ev, avatar) {
            var originalAvatar = JSON.parse(JSON.stringify(avatar || {}));
            // all avatars are stored as links in itsyou.online
            if (avatar) {
                avatar.link = avatar.source;
                avatar.fileupload = false;
            }
            return $q(function (resolve, reject) {
                $mdDialog.show({
                    controller: AvatarDialogController,
                    templateUrl: 'components/user/views/avatarDialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        user: vm.user,
                        data: avatar,
                        userService: UserService
                    }
                })
                    .then(
                        function (data) {
                            if (data.fx === 'delete') {
                                vm.user.avatars.splice(vm.user.avatars.indexOf(avatar), 1);
                            }
                            else if (data.fx === 'create') {
                                vm.user.avatars.push(data.data);
                            }
                            else if (data.fx === 'update') {
                                avatar.label = data.data.label;
                                avatar.source = data.data.source;
                            }
                            resolve(data);
                        }, function (response) {
                            angular.forEach(originalAvatar, function (value, key) {
                                avatar[key] = value;
                            });
                            reject(response);
                        });
            });
        }

        function addFacebook() {
            configService.getConfig(function (config) {
                $window.location.href = 'https://www.facebook.com/dialog/oauth?client_id='
                    + config.facebookclientid
                    + '&response_type=code&redirect_uri='
                    + $window.location.origin
                    + '/facebook_callback';
            });
        }

        function facebook(ev) {
            $mdDialog.show({
                controller: genericDetailControllerParams,
                templateUrl: 'components/user/views/facebookDialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    user: vm.user,
                    data: vm.user.facebook,
                    createFunction: doNothing,
                    updateFunction: doNothing,
                    deleteFunction: UserService.deleteFacebookAccount,
                    preFilledLabels: []
                }
            })
                .then(
                    function () {
                        vm.user.facebook = {};
                    });
        }

        function github(ev) {
            $mdDialog.show({
                controller: genericDetailControllerParams,
                templateUrl: 'components/user/views/githubDialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    user: vm.user,
                    data: vm.user.github,
                    createFunction: doNothing,
                    updateFunction: doNothing,
                    deleteFunction: UserService.deleteGithubAccount,
                    preFilledLabels: []
                }
            })
                .then(
                    function () {
                        vm.user.github = {};
                    });
        }

        function addGithub() {
            configService.getConfig(function (config) {
                $window.location.href = 'https://github.com/login/oauth/authorize/?client_id=' + config.githubclientid;
            });
        }

        /**
         *
         * @param message
         * @param title
         * @param closeText
         * @param event optional click event
         * @returns promise
         */
        function showSimpleDialog(message, title, closeText, event) {
            title = title || 'Alert';
            closeText = closeText || 'Close';
            return $mdDialog.show(
                $mdDialog.alert({
                    title: title,
                    htmlContent: message,
                    ok: closeText,
                    targetEvent: event
                })
            );
        }

        function createOrganization(ev, parentOrganization) {
            $mdDialog.show({
                controller: ['$scope', '$window', '$mdDialog', 'OrganizationService', 'UserService', 'parentOrganization', CreateOrganizationController],
                controllerAs: 'ctrl',
                templateUrl: 'components/organization/views/createOrganizationDialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    parentOrganization: parentOrganization
                }
            });
        }

        function digitalWalletAddressDetail(event, walletAddress) {
            walletAddress = walletAddress || {expire: new Date(new Date().setFullYear(new Date().getFullYear() + 1))};
            if (walletAddress.noexpiration == null) {
                walletAddress.noexpiration = true;
            }
            var originalWalletAddress = JSON.parse(JSON.stringify(walletAddress));
            walletAddress.expire = ['string'].indexOf(typeof walletAddress.expire) !== -1 ? new Date(walletAddress.expire) : walletAddress.expire;
            if (walletAddress.expire.getFullYear() < 2000) {
                walletAddress.expire = new Date();
            }
            return $q(function (resolve, reject) {
                var preFilledLabels = ["Primary", "Secondary", "Other"]
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/digitalWalletAddressDialog.html',
                    targetEvent: event,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        user: vm.user,
                        data: walletAddress,
                        createFunction: UserService.createDigitalWalletAddress,
                        updateFunction: UserService.updateDigitalWalletAddress,
                        deleteFunction: UserService.deleteDigitalWalletAddress,
                        preFilledLabels: preFilledLabels,
                    }
                }).then(
                    function (data) {
                        if (data.fx === 'delete') {
                            vm.user.digitalwallet.splice(vm.user.digitalwallet.indexOf(walletAddress), 1);
                        }
                        else if (data.fx === 'create') {
                            vm.user.digitalwallet.push(data.data);
                        }
                        resolve(data);
                    }, function (response) {
                        angular.forEach(originalWalletAddress, function (value, key) {
                            walletAddress[key] = value;
                        });
                        reject(response);
                    });
            });
        }

        function publicKey(event, key) {
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'UserService', 'username', 'key', PublicKeyDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/publicKeyDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                clickOutsideToClose: true,
                locals: {
                    UserService: UserService,
                    username: vm.username,
                    key: key
                }
            })
                .then(
                    function (data) {
                        if (data.originalLabel != data.newLabel) {
                            if (data.originalLabel) {
                                var originalKey = getPublicKey(data.originalLabel);
                                if (data.newLabel) {
                                    // update
                                    originalKey.label = data.newLabel;
                                }
                                else {
                                    // remove
                                    vm.user.publicKeys.splice(vm.user.publicKeys.indexOf(originalKey), 1);
                                }
                            }
                            else {
                                // add
                                vm.user.publicKeys.push(data.key);
                            }
                        }
                    });

            function getPublicKey(label) {
                return vm.user.publicKeys.filter(function (k) {
                    return k.label === label;
                })[0];
            }

            function PublicKeyDialogController($scope, $mdDialog, UserService, username, key) {
                var ctrl = this;
                ctrl.Key = key || {publickey: ""};
                ctrl.originalKey = key ? key.publickey : "";
                ctrl.originalLabel = key ? key.label : null;
                ctrl.savedLabel = key ? key.label : null;
                ctrl.label = key ? key.label : null;

                ctrl.cancel = cancel;
                ctrl.create = createPublicKey;
                ctrl.update = updatePublicKey;
                ctrl.delete = deletePublicKey;

                ctrl.modified = false;

                function cancel() {
                    if (ctrl.originalLabel) {
                        ctrl.Key.publickey = ctrl.originalKey;
                        ctrl.Key.label = ctrl.originalLabel;
                    }
                    $mdDialog.cancel();
                }

                function createPublicKey() {
                    ctrl.validationerrors = {};
                    UserService.createPublicKey(username, ctrl.label, ctrl.Key.publickey).then(
                        function (data) {
                            ctrl.modified = true;
                            ctrl.Key = data;
                            $mdDialog.hide({newLabel: ctrl.label, key: data});
                        },
                        function (reason) {
                            if (reason.status === 409) {
                                $scope.PublicKeyForm.label.$setValidity('duplicate', false);
                            }
                        }
                    );
                }

                function updatePublicKey() {
                    UserService.updatePublicKey(username, ctrl.savedLabel, ctrl.label, ctrl.Key.publickey).then(
                        function (data) {
                            $mdDialog.hide({originalLabel: ctrl.savedLabel, newLabel: ctrl.label});
                        },
                        function (reason) {
                            if (reason.status === 409) {
                                $scope.PublicKeyForm.label.$setValidity('duplicate', false);
                            }
                        }
                    );
                }

                function deletePublicKey() {
                    UserService.deletePublicKey(username, key.label).then(
                        function () {
                            $mdDialog.hide({originalLabel: key.label, newLabel: ""});
                        }
                    );
                }
            }
        }

        // Controller for the avatar view
        function AvatarDialogController($scope, $mdDialog, user, data, userService) {
            data = data || {};
            $scope.data = data;

            $scope.originalLabel = data.label;
            $scope.user = user;

            $scope.cancel = cancel;
            $scope.validationerrors = {};
            $scope.create = create;
            $scope.update = update;
            $scope.remove = remove;

            function cancel() {
                $mdDialog.cancel();
            }

            function create(data) {
                if (Object.keys($scope.dataform.$error).length > 0) {
                    return;
                }
                if (data.fileupload && !data.src) {
                    $scope.validationerrors.no_file_selected = true;
                    return;
                }
                $scope.validationerrors = {};
                var createFunction = data.fileupload ? userService.createAvatarFromFile : userService.createAvatarFromLink;
                createFunction(user.username, data.label,
                    data.fileupload ? data.src : data.link).then(
                    function (response) {
                        $mdDialog.hide({fx: 'create', data: response.data});
                    },
                    function (reason) {
                        if (reason.data && reason.data.error) {
                            $scope.validationerrors[reason.data.error] = true;
                        }
                        else if (reason.status === 413) {
                            $scope.validationerrors.file_too_large = true;
                        }
                    }
                );
            }

            function update(oldLabel, data) {
                if (Object.keys($scope.dataform.$error).length > 0) {
                    return;
                }
                $scope.validationerrors = {};
                var updateFunction = data.fileupload ? userService.updateAvatarFile : userService.updateAvatarLink;
                updateFunction(user.username, oldLabel, data.label,
                    data.fileupload ? data.src :  data.link).then(
                    function (response) {
                        $mdDialog.hide({fx: 'update', data: response.data});
                    },
                    function (response) {
                        if (response.data && response.data.error) {
                            $scope.validationerrors[response.data.error] = true;
                        }
                        else if (response.status === 413) {
                            $scope.validationerrors.file_too_large = true;
                        }
                    }
                );
            }

            function remove(label) {
                $scope.validationerrors = {};
                userService.deleteAvatar(user.username, label)
                    .then(
                        function () {
                            $mdDialog.hide({fx: 'delete'});
                        },
                        function (reason) {
                            if (reason.status === 409) {
                                $scope.validationerrors.delete_protected_label = true;
                            }
                        });
            }

            // we need to declare this function on the scope to statisfy the directive,
            // else it won't be able to evaluate it causing some nasty errors
            $scope.updateFile = function(event) {
                $scope.validationerrors.no_file_selected = false;
                $scope.validationerrors.file_too_large = false;
                var files = event.target.files;
                $scope.data.src = files[0];
                var url = URL.createObjectURL(files[0]);
                if (url) {
                    $scope.data.file = url;
                }
                $scope.$digest();
            };
        }

        function GenericDetailDialogController($scope, $mdDialog, user, data, createFunction, updateFunction, deleteFunction, preFilledLabels) {
            data = data || {};
            $scope.data = data;

            $scope.originalLabel = data.label;
            $scope.user = user;

            $scope.cancel = cancel;
            $scope.validationerrors = {};
            $scope.create = create;
            $scope.update = update;
            $scope.remove = remove;

            $scope.isCustomlabel = isCustomlabel;
            $scope.preFilledLabels = preFilledLabels || [];

            function cancel() {
                $mdDialog.cancel();
            }

            function create(data) {
                if (Object.keys($scope.dataform.$error).length > 0) {
                    angular.forEach($scope.dataform.$error, function (field) {
                        angular.forEach(field, function(errorField){
                            // Touch all the error fields, so we are sure the error message renders
                            errorField.$setTouched();
                        })
                    });
                    return;
                }
                $scope.validationerrors = {};
                createFunction(user.username, data).then(
                    function (response) {
                        $mdDialog.hide({fx: 'create', data: response});
                    },
                    function (reason) {
                        if (reason.status === 409) {
                            $scope.validationerrors.duplicate = true;
                        }
                    }
                );
            }

            function update(oldLabel, data) {
                if (Object.keys($scope.dataform.$error).length > 0) {
                    return;
                }
                $scope.validationerrors = {};
                updateFunction(user.username, oldLabel, data).then(
                    function (response) {
                        $mdDialog.hide({fx: 'update', data: response});
                    },
                    function (response) {
                        if (response.data && response.data.error) {
                            $scope.validationerrors[response.data.error] = true;
                        }
                        else if (response.status === 409) {
                            $scope.validationerrors.duplicate = true;
                        }
                        else if (response.status === 412) {
                            $scope.validationerrors.illegalactionforcurrentstate = true;
                        }
                    }
                );
            }

            function remove(label) {
                $scope.validationerrors = {};
                deleteFunction(user.username, label)
                    .then(function () {
                        $mdDialog.hide({fx: 'delete'});
                    });
            }

            function isCustomlabel(label) {
                if (!label) {
                    return false;
                }
                var custom = true;
                angular.forEach($scope.preFilledLabels, function(pfl) {
                    if (label === pfl) {
                        custom = false;
                    }
                })
                return custom;
            }

        }

        function CreateOrganizationController($scope, $window, $mdDialog, OrganizationService, UserService, parentOrganization) {
            var ctrl = this;
            ctrl.submit = submit;
            ctrl.cancel = cancel;
            ctrl.resetValidation = resetValidation;
            ctrl.name = '';
            ctrl.parentOrganization = parentOrganization || '';

            function submit() {
                ctrl.name = ctrl.name ? ctrl.name.trim() : '';
                if (!$scope.form.$valid) {
                    return;
                }
                OrganizationService
                    .create(ctrl.name, [], UserService.getUsername(), parentOrganization)
                    .then(
                        function (data) {
                            cancel();
                            $window.location.hash = '#/organization/' + encodeURIComponent(data.globalid) + '/people';
                        },
                        function (reason) {
                            if (reason.status === 409) {
                                $scope.form.name.$setValidity(reason.data.error, false);
                            } else if (reason.status === 400) {
                                $scope.form.name.$setValidity('pattern', false);
                            }
                            else if (reason.status === 422) {
                                cancel();
                                var msg = 'You cannot create any more organizations because you have reached the maximum amount of organizations.';
                                showSimpleDialog(msg, "Error");
                            }
                        }
                    );
            }

            function cancel() {
                $mdDialog.cancel();
            }
            function resetValidation() {
                $scope.form.name.$setValidity('organization_exists', true);
                $scope.form.name.$setValidity('user_exists', true);
            }
        }
    }
})();
