(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('registrationController', [
            '$scope', '$window', '$cookies', '$mdMedia','$mdUtil', '$rootScope', '$timeout', '$http', 'configService', 'registrationService',
            registrationController]);

    function registrationController($scope, $window, $cookies, $mdMedia, $mdUtil, $rootScope, $timeout, $http, configService, registrationService) {
        var vm = this,
            queryParams = URI($window.location.href).search(true);
        vm.resendValidation = resendValidation;
        vm.register = register;
        vm.resetValidation = resetValidation;
        vm.basicInfoValid = basicInfoValid;
        vm.onTabSelected = onTabSelected;
        vm.goToNextTabIfValid = goToNextTabIfValid;
        vm.goToPreviousTab = goToPreviousTab;
        vm.codeUpdated = codeUpdated;
        vm.externalSite = queryParams.client_id;
        $rootScope.loginUrl = '/login' + $window.location.search;
        vm.goToLogin = goToLogin;
        vm.goToForgotPassword = goToForgotPassword;
        vm.logo = undefined;
        vm.description = "";
        vm.selectedTab = 0;
        vm.oldSelectedTab = 0;
        vm.phone = {};
        vm.phone.validationerrors = {};
        vm.phoneduplicate = false;
        vm.emailduplicate = false;

        vm.sms = "";
        vm.smsvalidation = "";
        vm.email = "";
        vm.emailvalidation = "";
        vm.firstname = "";
        vm.firstnamevalidation = "";
        vm.lastname = "";
        vm.lastnamevalidation = "";
        vm.passwordvalidation = "";

        vm.emailConfirmed = false;
        vm.phoneConfirmed = false;

        vm.basicInfoFailed = false;

        init();

        function init() {
            $scope.$watch(function () {
                return $mdMedia('gt-sm');
            }, function (isGtSm) {
                vm.mobileView = !isGtSm;
            });
            $scope.$watch(function () {
                return vm.phone.validationerrors.phone_already_used;
            }, function(phoneused) {
                vm.phoneduplicate = phoneused;
            });
            if (queryParams && queryParams.scope && queryParams.scope.includes('ownerof:email')) {
                var scopes = queryParams.scope.split(',');
                for (var i = 0; i < scopes.length; i++) {
                    if (scopes[i].includes('ownerof:email')) {
                        var parts = scopes[i].split(':');
                        vm.email = parts[3];
                        break;
                    }
                }
            }
            if (vm.externalSite) {
                registrationService.getLogo(vm.externalSite).then(function (data) {
                    vm.logo = data.logo;
                });
                loadDescription();
            }
        }

        // Load the correct description after the user changes language
        $rootScope.$on('$translateChangeSuccess', function () {
            loadDescription();
        });

        function loadDescription() {
            registrationService.getDescription(vm.externalSite, localStorage.getItem('langKey')).then(
                function(data) {
                    vm.description = data.text;
                }
            );
        }

        function goToLogin() {
            $window.location.href = $rootScope.loginUrl;
        }

        function goToForgotPassword() {
            $window.location.href = '/login#/forgotpassword';
        }

        function register() {
            if(!$scope.signupform.$valid){
                return;
            }
            var redirectparams = $window.location.search.replace('?', '');
            registrationService
                .register(vm.firstname, vm.lastname, vm.email, vm.emailcode, vm.sms, vm.smscode, vm.password, redirectparams)
                .then(function (response) {
                    var url = response.data.redirecturl;
                    if (url === '/') {
                        $cookies.remove('registrationdetails');
                    }
                    $window.location.href = url;
                }, function (response) {
                    switch (response.status) {
                        case 422:
                            var err = response.data.error;
                            switch (err) {
                                case 'invalid_first_name':
                                    $scope.signupform.firstname.$setValidity(err, false);
                                    break;
                                case 'invalid_last_name':
                                    $scope.signupform.lastname.$setValidity(err, false);
                                    break;
                                case 'invalid_phonenumber':
                                    vm.phone.validationerrors.invalid_phone = true;
                                    break;
                                case 'invalid_totpcode':
                                    $scope.signupform.totpcode.$setValidity(err, false);
                                    break;
                                case 'invalid_password':
                                    $scope.signupform.password.$setValidity(err, false);
                                    break;
                                case 'invalid_email_code':
                                    $scope.signupform.emailcode.$setValidity(err, false);
                                    break;
                                case 'invalid_sms_code':
                                    $scope.signupform.smscode.$setValidity(err, false);
                                    break;
                                case 'invalid_email_format':
                                    $scope.signupform.email.$setValidity('email', false);
                                    break;
                                case 'email_already_used':
                                    $scope.signupform.email.$setValidity(err, false);
                                    vm.emailduplicate = true;
                                    break;
                                case 'phone_already_used':
                                    vm.phone.validationerrors[err] = true;
                                    vm.phoneduplicate = true;
                                    break;
                                default:
                                    console.error('Unconfigured error:', response.data.error);
                            }
                            break;
                        case 409:
                            $scope.signupform.login.$setValidity('user_exists', false);
                            break;
                    }
                });
        }

        function resetValidation(prop) {
            switch (prop) {
                case 'firstname':
                    $scope.signupform[prop].$setValidity("invalid_first_name", true);
                    break;
                case 'lastname':
                    $scope.signupform[prop].$setValidity("invalid_last_name", true);
                    break;
                case 'smscode':
                    $scope.signupform[prop].$setValidity("invalid_sms_code", true);
                    break;
                case 'email':
                    $scope.signupform[prop].$setValidity("email", true);
                    $scope.signupform[prop].$setValidity("email_already_used", true);
                    vm.emailduplicate = false;
                    break;
                case 'phonenumber':
                    $scope.signupform[prop].$setValidity("invalid_phonenumber", true);
                    vm.phone.validationerrors.invalid_phone = false;
                    vm.phone.validationerrors.phone_already_used = false;
                    vm.phoneduplicate = false;
                    break;
                case 'totpcode':
                    $scope.signupform[prop].$setValidity("invalid_totpcode", true);
                    break;
                case 'twoFAMethod':
                    if ($scope.signupform.totpcode) {
                        $scope.signupform.totpcode.$setValidity("totpcode", true);
                    }
                    if ($scope.signupform.phonenumber) {
                        $scope.signupform.phonenumber.$setValidity("invalid_phonenumber", true);
                        $scope.signupform.phonenumber.$setValidity("pattern", true);
                    }
                    break;
            }
        }

        function basicInfoValid() {
            return $scope.signupform.firstname
                && $scope.signupform.firstname.$valid
                && $scope.signupform.lastname.$valid
                && vm.sms.length > 5
                // double boolean negation to cast to the real boolean form, undefined -> false
                && !!vm.phone.validationerrors.pattern === false
                && !!vm.phone.validationerrors.invalid_phone === false
                && !!vm.phone.validationerrors.phone_already_used === false
                && $scope.signupform.email.$valid
                && $scope.signupform.password.$valid
                && $scope.signupform.passwordvalidation.$valid;
        }

        function goToNextTabIfValid() {
            vm.selectedTab += 1;
        }

        function goToPreviousTab() {
            vm.oldSelectedTab = vm.selectedTab;
            vm.selectedTab -= 1;
        }

        function onTabSelected() {
            if ((vm.oldSelectedTab === 0 && vm.selectedTab === 1 && vm.selectedTab != vm.oldSelectedTab) || vm.basicInfoFailed ) {
                if (basicInfoValid()) {
                    // blur the password validaton field 
                    document.getElementById('passwordvalidation').blur();
                    requestValidationInfo()
                }
            }
            vm.oldSelectedTab = vm.selectedTab;
        }

        function requestValidationInfo() {
            if (basicInfoValid() && ((vm.sms != vm.smsvalidation || vm.email != vm.emailvalidation ||
                vm.firstname != vm.firstnamevalidation || vm.lastname != vm.lastnamevalidation) || vm.basicInfoFailed)) {

                vm.basicInfoFailed = false;
                registrationService.requestValidation(vm.firstname, vm.lastname, vm.email, vm.sms, vm.password).then(
                    function(success) {
                        startCodePolling();
                    },
                    function(failure) {
                        var err = failure.data.error;
                        vm.basicInfoFailed = true;
                        switch (err) {
                            case 'invalid_first_name':
                                $scope.signupform.firstname.$setValidity(err, false);
                                break;
                            case 'invalid_last_name':
                                $scope.signupform.lastname.$setValidity(err, false);
                                break;
                            case 'invalid_phonenumber':
                                vm.phone.validationerrors.invalid_phone = true;
                                break;
                            case 'invalid_totpcode':
                                $scope.signupform.totpcode.$setValidity(err, false);
                                break;
                            case 'invalid_password':
                                $scope.signupform.password.$setValidity(err, false);
                                break;
                            case 'invalid_sms_code':
                                $scope.signupform.smscode.$setValidity(err, false);
                                break;
                            case 'invalid_email_format':
                                $scope.signupform.email.$setValidity('email', false);
                                break;
                            case 'email_already_used':
                                $scope.signupform.email.$setValidity(err, false);
                                vm.emailduplicate = true;
                                break;
                            case 'phone_already_used':
                                vm.phone.validationerrors[err] = true;
                                vm.phoneduplicate = true;
                                break;
                            default:
                                console.error('Unconfigured error:', failure.data.error);
                        }
                    }
                )
                if (vm.smsvalidation != vm.sms) {
                    vm.phoneConfirmed = false;
                    // Clear the input field
                    vm.smscode = "";
                }
                vm.smsvalidation = vm.sms;
                if (vm.emailvalidation != vm.email) {
                    vm.emailConfirmed = false;
                }
                vm.emailvalidation = vm.email;
                vm.firstnamevalidation = vm.firstname;
                vm.lastnamevalidation = vm.lastname;
                vm.passwordvalidation = vm.password;
            }
        }

        function startCodePolling() {
            $timeout(checkPhoneConfirmation, 1000);
            $timeout(checkEmailConfirmation, 1000);
        }

        function codeUpdated() {
            vm.resetValidation('smscode');
            if ($scope.signupform.smscode.$invalid) {
                return;
            }
            submitSMSCode();
        }

        function submitSMSCode() {
            registrationService.submitSMSCode(vm.smscode).then(
                function(response) {
                    if (response.data.confirmed) {
                        vm.phoneConfirmed = response.data.confirmed;
                        // Don't trigger the email verification. It will be triggered by the next return of 
                        // checkPhoneConfirmation
                    } 
                },
                function(failure) {
                    var err = failure.data.error;
                    $scope.signupform.smscode.$setValidity(err, false);
                }
            );
        }

        function checkPhoneConfirmation() {
            $http.get('register/smsconfirmed' + $window.location.search).then(
                function(response) {
                    if (response.data.confirmed) {
                        vm.phoneConfirmed = response.data.confirmed;
                        // trigger email validation
                        resendValidation();
                    } else {
                        $timeout(checkPhoneConfirmation, 1000);
                    }
                },
                function() {
                    $timeout(checkPhoneConfirmation, 1000);
                }
            );
        }

        function checkEmailConfirmation() {
            $http.get('register/emailconfirmed' + $window.location.search).then(
                function(response) {
                    if (response.data.confirmed) {
                        vm.emailConfirmed = response.data.confirmed;
                    } else {
                        $timeout(checkEmailConfirmation, 1000);
                    }
                },
                function() {
                    $timeout(checkEmailConfirmation, 1000);
                }
            );
        }

        function resendValidation() {
            registrationService.resendValidation(vm.email, vm.sms);
        }

    }
})();
