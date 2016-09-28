(function () {
    'use strict';
    angular.module('loginApp')
        .controller('twoFactorAuthenticationController', ['$scope', '$window', '$interval', 'LoginService',
            twoFactorAuthenticationController]);

    function twoFactorAuthenticationController($scope, $window, $interval, LoginService) {
        var STEP_CHOICE = 'choice',
            STEP_CODE = 'code';
        var vm = this;
        vm.resetValidation = resetValidation;
        vm.shouldShowSendButton = shouldShowSendButton;
        vm.sendSmsCode = sendSmsCode;
        vm.login = login;
        vm.getHelpText = getHelpText;
        vm.nextStep = nextStep;
        vm.selectedTwoFaMethod = null;
        vm.hasMoreThanOneTwoFaMethod = false;
        var steps = [STEP_CHOICE, STEP_CODE];
        vm.step = steps[0];
        var interval;
        var queryString = $window.location.search;
        init();

        function init() {
            LoginService
                .getTwoFactorAuthenticationMethods()
                .then(function (data) {
                    vm.possibleTwoFaMethods = {};
                    if (data['totp']) {
                        vm.possibleTwoFaMethods['totp'] = 'Authenticator application';
                    }
                    if (data['sms'] && Object.keys(data['sms']).length) {
                        angular.forEach(data['sms'], function (sms, label) {
                            vm.possibleTwoFaMethods['sms-' + label] = 'SMS - ' + sms + ' (' + label + ')';
                        });
                    }
                    var methods = Object.keys(vm.possibleTwoFaMethods);
                    if (!methods.length) {
                        // Redirect to resend sms page
                        $window.location.hash = '#/resendsms';
                        return;
                    }
                    vm.hasMoreThanOneTwoFaMethod = methods.length > 1;
                    if (!vm.hasMoreThanOneTwoFaMethod) {
                        vm.selectedTwoFaMethod = methods[0];
                        nextStep();
                    } else {
                        // Preselect based on what was selected when previously logging in
                        var method = localStorage.getItem('itsyouonline.last2falabel');
                        if (!method || methods.indexOf(method) === -1) {
                            method = methods[0];
                        }
                        vm.selectedTwoFaMethod = method;
                    }
                });
        }

        function nextStep() {
            vm.step = steps[steps.indexOf(vm.step) + 1];
            if (vm.step === STEP_CODE && vm.selectedTwoFaMethod.indexOf('sms-') === 0) {
                sendSmsCode();
            }
        }

        function getHelpText() {
            var text = '';
            if (vm.step === STEP_CODE) {
                if (vm.selectedTwoFaMethod.indexOf('sms-') === 0) {
                    text = 'Click the link in the sms sent to your phone or enter the code from the sms here to continue.';
                }
                if (vm.selectedTwoFaMethod === 'totp') {
                    text = 'Fill in the 6 digit code from the authenticator application on your phone.';
                }
            }
            return text;
        }

        function shouldShowSendButton() {
            return vm.selectedTwoFaMethod && vm.selectedTwoFaMethod.indexOf('sms-') === 0 && vm.step === STEP_CODE;
        }

        function resetValidation() {
            $scope.twoFaForm.code.$setValidity("invalid_code", true);
        }

        function sendSmsCode() {
            if (interval) {
                $interval.cancel(interval);
            }
            var phoneLabel = vm.selectedTwoFaMethod.replace('sms-', '');
            LoginService
                .sendSmsCode(phoneLabel)
                .then(function () {
                    interval = $interval(checkSmsConfirmation, 1000);
                });
        }

        function login() {
            var method;
            if (vm.selectedTwoFaMethod === 'totp') {
                method = LoginService.submitTotpCode;
            } else if (vm.selectedTwoFaMethod.indexOf('sms-') === 0) {
                method = LoginService.submitSmsCode;
            }
            method(vm.code, queryString)
                .then(
                    function (data) {
                        localStorage.setItem('itsyouonline.last2falabel', vm.selectedTwoFaMethod);
                        if (interval) {
                            $interval.cancel(interval);
                        }
                        goToPage(data.redirecturl);
                    },
                    function (response) {
                        switch (response.status) {
                            case 422:
                                $scope.twoFaForm.code.$setValidity("invalid_code", false);
                                break;
                        }
                    });
        }

        function checkSmsConfirmation() {
            LoginService.checkSmsConfirmation()
                .then(function (data) {
                    if (data.confirmed) {
                        login();
                    }
                });
        }

        function goToPage(url) {
            if (interval) {
                $interval.cancel(interval);
            }
            $window.location.href = url;
        }

    }
})();
