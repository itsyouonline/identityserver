(function () {
    'use strict';
    angular.module('loginApp')
        .controller('loginController', ['$http', '$window', '$scope', '$rootScope', 'LoginService', loginController]);

    function loginController($http, $window, $scope, $rootScope, LoginService) {
        var vm = this;
        vm.submit = submit;
        vm.clearValidation = clearValidation;
        vm.externalSite = URI($window.location.href).search(true).client_id;
        $rootScope.registrationUrl = '/register' + $window.location.search;
        vm.logo = "";

        activate();

        function activate() {
            if (vm.externalSite) {
                LoginService.getLogo(vm.externalSite).then(
                    function(data) {
                        vm.logo = data.logo;
                        renderLogo();
                    }
                );
            }
        }

        function renderLogo() {
            if (vm.logo !== "") {
                var img = new Image();
                img.src = vm.logo;

                var c = document.getElementById("login-logo");
                if (!c) {
                    console.log("aborting logo render - canvas not loaded");
                    return;
                }
                var ctx = c.getContext("2d");
                ctx.clearRect(0, 0, c.width, c.height);
                ctx.drawImage(img, 0, 0);
            }
        }

        function submit() {
            var data = {
                login: vm.login,
                password: vm.password
            };
            var url = '/login' + $window.location.search;
            $http.post(url, data).then(
                function (data) {
                  if (data.data.redirecturl) {
                      // Skip 2FA when logging in from an external site if the 2FA validity period hasn't passed
                      $window.location.href = data.data.redirecturl;
                  } else {
                      // Redirect 2 factor authentication page
                      $window.location.hash = '#/2fa';
                  }
                },
                function (response) {
                    if (response.status === 422) {
                        $scope.loginform.password.$setValidity("invalidcredentials", false);
                    }
                }
            );
        }

        function clearValidation() {
            $scope.loginform.password.$setValidity("invalidcredentials", true);
        }
    }
})();
