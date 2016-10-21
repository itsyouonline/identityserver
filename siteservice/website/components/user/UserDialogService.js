/**
 * Created by lucas on 13/06/16.
 */
(function () {
    'use strict';

    angular
        .module('itsyouonline.user', [])
        .factory('UserDialogService', ['$window', '$q', '$interval', '$mdMedia', '$mdDialog', 'UserService', 'configService', UserDialogService]);

    function UserDialogService($window, $q, $interval, $mdMedia, $mdDialog, UserService, configService) {
        var vm;
        var genericDetailControllerParams = ['$scope', '$mdDialog', 'user', 'data',
            'createFunction', 'updateFunction', 'deleteFunction', GenericDetailDialogController];
        return {
            init: init,
            emailDetail: emailDetail,
            addressDetail: addressDetail,
            phonenumberDetail: phonenumberDetail,
            verifyPhone: verifyPhone,
            bankAccount: bankAccount,
            facebook: facebook,
            addFacebook: addFacebook,
            github: github,
            addGithub: addGithub,
            showSimpleDialog: showSimpleDialog,
            createOrganization: createOrganization,
            digitalWalletAddressDetail: digitalWalletAddressDetail
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
                        deleteFunction: UserService.deleteEmailAddress
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
                                    var errorMsg, dialog;
                                    if (response.data.error === 'warning_delete_last_verified_phone_number') {
                                        errorMsg = 'Are you sure you want to delete this phone number? <br />' +
                                            'It is your last verified phone number, which means you will <br />' +
                                            'no longer be able to login using sms confirmations.';
                                        dialog = $mdDialog.confirm()
                                            .title('Confirm deletion')
                                            .ok('Confirm')
                                            .cancel('Cancel');
                                    }
                                    dialog = dialog.htmlContent(errorMsg)
                                        .ariaLabel('Delete phone number')
                                        .targetEvent(ev);
                                    $mdDialog.show(dialog)
                                        .then(function () {
                                            UserService
                                                .deletePhonenumber(username, label, true)
                                                .then(function () {
                                                    // Manually remove phone number since the dialog which executes the updatePhoneNumber promise callback had been closed before
                                                    vm.user.phonenumbers.splice(vm.user.phonenumbers.indexOf(phone), 1);
                                                }, function (response) {
                                                    if (response.data.error === 'cannot_delete_last_verified_phone_number') {
                                                        errorMsg = 'You cannot delete your last verified phone number. <br />' +
                                                            'Please change your 2 factor authentication method or add another verified phone number.';
                                                    } else {
                                                        errorMsg = 'Could not delete phone number. Please try again later.';
                                                    }
                                                    showSimpleDialog(errorMsg, 'Error', 'Ok', ev);
                                                });
                                        });
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
                        deleteFunction: deletePhoneNumber
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
                        .getVerifiedPhones(vm.username)
                        .then(function success(confirmedPhones) {
                            var confirmed = confirmedPhones.filter(function (p) {
                                    return p.label === ctrl.label;
                                }).length !== 0;
                            if (confirmed) {
                                findByLabel('phonenumbers', ctrl.label).verified = true;
                                close();
                            }
                        });
                }

                function submit() {
                    UserService
                        .verifyPhone(user.username, ctrl.label, ctrl.validationKey, ctrl.smscode)
                        .then(function () {
                            findByLabel('phonenumbers', ctrl.label).verified = true;
                            close();
                        }, function (response) {
                            if (response.status === 422) {
                                $scope.form.smscode.$setValidity('invalid_code', false);
                            }
                        });
                }

                function resetValidation() {
                    $scope.form.smscode.$setValidity('invalid_code', true);
                }
            }
        }

        function addressDetail(ev, address) {
            return $q(function (resolve, reject) {
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
                        deleteFunction: UserService.deleteAddress
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
                        deleteFunction: UserService.deleteBankAccount
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
                    deleteFunction: UserService.deleteFacebookAccount
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
                    deleteFunction: UserService.deleteGithubAccount
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
                controller: ['$scope', '$rootScope', '$window', '$mdDialog', 'OrganizationService', 'parentOrganization', CreateOrganizationController],
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
            walletAddress = walletAddress || {expire: new Date()};
            var originalWalletAddress = JSON.parse(JSON.stringify(walletAddress));
            walletAddress.expire = ['string'].indexOf(typeof walletAddress.expire) !== -1 ? new Date(walletAddress.expire) : walletAddress.expire;
            if (walletAddress.expire.getFullYear() < 2000) {
                walletAddress.expire = new Date();
            }
            return $q(function (resolve, reject) {
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
                        deleteFunction: UserService.deleteDigitalWalletAddress
                    }
                }).then(
                    function (data) {
                        if (data.fx === 'delete') {
                            vm.user.digitalwallet.splice(vm.user.digitalwallet.indexOf(walletAddress), 1);
                        }
                        else if (data.fx === 'create') {
                            vm.user.digitalwallet = data.data
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

        function GenericDetailDialogController($scope, $mdDialog, user, data, createFunction, updateFunction, deleteFunction) {
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

        }

        function CreateOrganizationController($scope, $rootScope, $window, $mdDialog, OrganizationService, parentOrganization) {
            var ctrl = this;
            ctrl.submit = submit;
            ctrl.cancel = cancel;
            ctrl.resetValidation = resetValidation;
            ctrl.name = '';
            ctrl.parentOrganization = parentOrganization || '';

            function submit() {
                if (!$scope.form.$valid) {
                    return;
                }
                OrganizationService
                    .create(ctrl.name, [], $rootScope.user, parentOrganization)
                    .then(
                        function (data) {
                            cancel();
                            $window.location.hash = '#/organization/' + encodeURIComponent(data.globalid);
                        },
                        function (reason) {
                            if (reason.status === 409) {
                                $scope.form.name.$setValidity('duplicate', false);
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
                $scope.form.name.$setValidity('duplicate', true);
            }
        }
    }
})();
