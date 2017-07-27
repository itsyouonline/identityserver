(function () {
    'use strict';
    angular.module('loginApp')
        .controller('organizationInviteController', ['$http', '$window', '$routeParams', '$mdDialog', organizationInviteController]);

    function organizationInviteController($http, $window, $routeParams, $mdDialog) {
        var vm = this;
        var code = encodeURIComponent($routeParams.code);
        vm.loginUrl = '?invitecode=' + code + '#/';
        vm.registerUrl = '/register?invitecode=' + code + '#/';

        init();

        function init() {
            $http.get('/login/organizationinvitation/' + code).then(function (response) {
                    vm.invite = response.data;
                },
                function (response) {
                    switch (response.status) {
                        case 404:
                            var msg = 'The organization invite has expired or could not be found.';
                            showErrorMessage(msg);
                            break;
                        default:
                            $window.location.href = 'error' + response.status;
                    }
                }
            );
        }

        function showErrorMessage(msg) {
            $mdDialog.show($mdDialog.alert()
                .clickOutsideToClose(true)
                .title('Error')
                .textContent(msg)
                .ariaLabel('Error: ' + msg)
                .ok('ok'));
        }
    }
})();
