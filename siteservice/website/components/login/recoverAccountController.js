(function () {
    'use strict';
    angular.module('loginApp')
        .controller('resetPasswordController', ['$http', '$window', '$routeParams', '$mdDialog', resetPasswordController]);

    function resetPasswordController($http, $window, $routeParams, $mdDialog) {
        var vm = this;
        vm.submit = submit;
        var code = $routeParams.code;

        function submit() {
            var data = {
                password: vm.password,
                token: code
            };
            $http
                .post('/login/resetpassword', data)
                .then(function () {
                        // redirect to login
                        $window.location.hash = '';
                    },
                    function (response) {
                        switch (response.status) {
                            case 404:
                                var msg = 'The password reset token was already used or was not found.';
                                showErrorMessage(msg);
                                break;
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
