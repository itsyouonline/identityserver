(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('registrationController', [
            '$scope', '$window', '$cookies', '$mdUtil', '$rootScope', 'configService', 'registrationService',
            registrationController]);

    function registrationController($scope, $window, $cookies, $mdUtil, $rootScope, configService, registrationService) {
        var vm = this,
            queryParams = URI($window.location.href).search(true);
        configService.getConfig(function (config) {
            vm.totpsecret = config.totpsecret;
            vm.totpissuer = encodeURIComponent(config.totpissuer);
        });
        vm.register = register;
        vm.resetValidation = resetValidation;
        vm.basicInfoValid = basicInfoValid;
        vm.getQrCodeData = getQrCodeData;
        vm.goToNextTabIfValid = goToNextTabIfValid;
        vm.externalSite = queryParams.client_id;
        $rootScope.loginUrl = '/login' + $window.location.search;
        vm.logo = "";
        vm.twoFAMethod = 'sms';
        vm.description = "";
        vm.selectedTab = 0;
        vm.validateUsername = $mdUtil.debounce(function () {
            $scope.signupform.login.$setValidity("organization_exists", true);
            $scope.signupform.login.$setValidity("user_exists", true);
            $scope.signupform.login.$setValidity("invalid_username_format", true);
            if ($scope.signupform.login.$valid) {
                registrationService
                    .validateUsername(vm.login.toLowerCase())
                    .then(function (response) {
                        $scope.signupform.login.$setValidity(response.data.error, response.data.valid);
                    });
            }
        }, 500, true);

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
            if (vm.externalSite) {
                registrationService.getLogo(vm.externalSite).then(
                    function(data) {
                        vm.logo = data.logo;
                        resizeLogo();
                    }
                );
                window.addEventListener('resize', resizeLogo, false);
                window.addEventListener('orientationchange', resizeLogo, false);
                loadDescription();
            }
        }

        function getQrCodeData() {
            return 'otpauth://totp/' + vm.totpissuer + ':' + vm.login + '?secret=' + vm.totpsecret + '&issuer=' + vm.totpissuer;
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

        function renderLogo() {
            if (vm.logo !== "") {
                var img = new Image();
                img.onload = function() {
                    var c = document.getElementById("register-logo");
                    if (!c) {
                        return;
                    }
                    var ctx = c.getContext("2d");
                    ctx.clearRect(0, 0, c.width, c.height);
                    ctx.drawImage(img, 0, 0, c.width, c.height);
                };
                img.src = vm.logo;
            }
        }

        function resizeLogo(e) {
            var formArea = document.getElementById("form-area");
            var logoArea = document.getElementById("register-logo");
            var widthToHeight = 25 / 12;
            var newWidth = formArea.clientWidth - 20;
            if (newWidth < 500) {
                logoArea.width = newWidth;
                logoArea.height = (newWidth) / widthToHeight;
            } else if (newWidth >= 500 && logoArea.width < 500) {
                logoArea.width = 500;
                logoArea.height = 240;
            }
            renderLogo();
        }

        function register() {
            if(!$scope.signupform.$valid){
                return;
            }
            var redirectparams = $window.location.search.replace('?', '');
            registrationService
                .register(vm.twoFAMethod, vm.login.toLowerCase(), vm.email, vm.password, vm.totpcode, vm.sms, redirectparams)
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
                case 'phonenumber':
                    $scope.signupform[prop].$setValidity("invalid_phonenumber", true);
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
