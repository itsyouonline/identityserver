(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('registrationController', [
            '$scope', '$window', '$cookies', '$mdUtil', '$rootScope', '$timeout', '$http', 'configService', 'registrationService',
            registrationController]);

    function registrationController($scope, $window, $cookies, $mdUtil, $rootScope, $timeout, $http, configService, registrationService) {
        var vm = this,
            queryParams = URI($window.location.href).search(true);
        vm.resendValidation = resendValidation;
        vm.register = register;
        vm.resetValidation = resetValidation;
        vm.basicInfoValid = basicInfoValid;
        vm.onTabSelected = onTabSelected;
        vm.goToNextTabIfValid = goToNextTabIfValid;
        vm.externalSite = queryParams.client_id;
        $rootScope.loginUrl = '/login' + $window.location.search;
        vm.logo = undefined;
        vm.description = "";
        vm.selectedTab = 0;
        vm.oldSelectedTab = 0;
        vm.phone = {};
        vm.phone.validationerrors = {};

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
        // Do we need a validated email address to register?
        vm.needDoubleValidation = false;


        init();

        function init() {
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
            // require a validated email to register if:
            //  - a validated email scope is required to log in to an external org
            //  - the user is registering against IYO (not in an oauth flow) -> no externalsite
            //  - the `requirevalidatedemail` queryparameter is set.
            if ((queryParams && queryParams.scope && queryParams.scope.includes('user:validated:email'))
                || !vm.externalSite || (queryParams && queryParams.requirevalidatedemail)) {
                vm.needDoubleValidation = true;
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
                                    break
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
                case 'smscode':
                    $scope.signupform[prop].$setValidity("invalid_sms_code", true);
                    break;
                case 'emailcode':
                    $scope.signupform[prop].$setValidity("invalid_email_code", true);
                    break;
                case 'phonenumber':
                    $scope.signupform[prop].$setValidity("invalid_phonenumber", true);
                    vm.phone.validationerrors.invalid_phone = false;
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
                && $scope.signupform.email.$valid
                && $scope.signupform.password.$valid
                && $scope.signupform.passwordvalidation.$valid;
        }

        function goToNextTabIfValid() {
            vm.selectedTab = 1;
        }

        function onTabSelected() {
            if (vm.selectedTab === 1 && vm.selectedTab != vm.oldSelectedTab) {
                requestValidationInfo()
            }
            vm.oldSelectedTab = vm.selectedTab;
        }

        function requestValidationInfo() {
            if (basicInfoValid() && (vm.sms != vm.smsvalidation || vm.email != vm.emailvalidation ||
                vm.firstname != vm.firstnamevalidation || vm.lastname != vm.lastnamevalidation)) {
                registrationService.requestValidation(vm.firstname, vm.lastname, vm.email, vm.sms, vm.password).then(
                    function(success) {
                        startCodePolling();
                    },
                    function(failure) {
                        var err = failure.data.error;
                        switch (err) {
                          case 'invalid_phonenumber':
                              vm.phone.validationerrors.invalid_phone = true;
                              break;
                          case 'invalid_password':
                              $scope.signupform.password.$setValidity('password', false);
                              break;
                          case 'invalid_email_format':
                              $scope.signupform.email.$setValidity('email', false);
                              break
                          default:
                              console.error('Unconfigured error:', failure.data.error);
                        }
                    }
                )
                vm.smsvalidation = vm.sms;
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

        function checkPhoneConfirmation() {
            $http.get('register/smsconfirmed' + $window.location.search).then(
                function(response) {
                    if (response.data.confirmed) {
                        vm.phoneConfirmed = response.data.confirmed;
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
          console.log("sending validaton")
            registrationService.resendValidation(vm.email, vm.sms);
        }

    }
})();
