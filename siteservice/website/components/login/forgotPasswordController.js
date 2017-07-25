(function () {
    'use strict';
    angular.module('loginApp')
        .controller('forgotPasswordController', ['$http', '$scope', forgotPasswordController]);

    function forgotPasswordController($http, $scope) {
        var vm = this;
        vm.submit = submit;
        vm.clearValidation = clearValidation;
        vm.emailSend = false;
        function submit() {
            var data = {
                login: vm.login,
                langkey: localStorage.getItem('langKey')
            };
            $http.post('/login/forgotpassword', data).then(
                function () {
                    vm.emailSend = true;
                },
                function (response) {
                    switch (response.status) {
                        case 404:
                            $scope.form.login.$setValidity("invalid", false);
                            break;
                        case 412:
                            $scope.form.login.$setValidity("no_validated_email", false);
                            break;
                    }
                }
            );
        }

        function clearValidation() {
            $scope.form.login.$setValidity("invalid", true);
        }
    }
})();
