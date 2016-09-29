(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('registrationController', [
            '$scope', '$window', '$cookies', '$mdUtil', 'configService', 'registrationService',
            registrationController]);

    function registrationController($scope, $window, $cookies, $mdUtil, configService, registrationService) {
        var vm = this;
        configService.getConfig(function (config) {
            vm.totpsecret = config.totpsecret;
        });
        vm.register = register;
        vm.resetValidation = resetValidation;
        vm.basicInfoValid = basicInfoValid;
        vm.goToNextTabIfValid = goToNextTabIfValid;
        vm.twoFAMethod = 'sms';
        vm.selectedTab = 0;
        vm.validateUsername = $mdUtil.debounce(function () {
            $scope.signupform.login.$setValidity("duplicate_username", true);
            $scope.signupform.login.$setValidity("invalid_username_format", true);
            if ($scope.signupform.login.$valid) {
                registrationService
                    .validateUsername(vm.login)
                    .then(function (response) {
                        $scope.signupform.login.$setValidity(response.data.error, response.data.valid);
                    });
            }
        }, 500, true);

        function register() {
            var redirectparams = $window.location.search.replace('?', '');
            registrationService
                .register(vm.twoFAMethod, vm.login, vm.email, vm.password, vm.totpcode, vm.sms, redirectparams)
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
                                    $scope.signupform.phonenumber.$setValidity(err, false);
                                    break;
                                case 'invalid_totpcode':
                                    $scope.signupform.totpcode.$setValidity(err, false);
                                    break;
                                case 'invalid_password':
                                    $scope.signupform.password.$setValidity(err, false);
                                    break;
                                case 'invalid_username_format':
                                    $scope.signupform.login.$setValidity(err, false);
                                    break;
                                default:
                                    console.error('Unconfigured error:', response.data.error);
                            }
                            break;
                        case 409:
                            $scope.signupform.login.$setValidity('duplicate_username', false);
                            break;
                    }
                });
        }

        function resetValidation(prop) {
            switch (prop) {
                case 'phonenumber':
                    $scope.signupform[prop].$setValidity("invalid_phonenumber", true);
                    break;
                case 'totpcode':
                    $scope.signupform[prop].$setValidity("invalid_totpcode", true);
                    break;
                case 'twoFAMethod':
                    $scope.signupform.totpcode.$setValidity("totpcode", true);
                    $scope.signupform.phonenumber.$setValidity("invalid_phonenumber", true);
                    $scope.signupform.phonenumber.$setValidity("pattern", true);
                    break;
            }
        }

        function basicInfoValid() {
            return $scope.signupform.login
                && $scope.signupform.login.$valid
                && $scope.signupform.email.$valid
                && $scope.signupform.password.$valid
                && $scope.signupform.passwordvalidation.$valid;
        }

        function goToNextTabIfValid() {
            vm.selectedTab = 1;
        }
    }
})();
