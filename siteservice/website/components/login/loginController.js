(function () {
    'use strict';
    angular.module('loginApp')
        .controller('loginController', ['$http', '$window', '$scope', '$rootScope', '$mdUtil', '$interval', 'configService', 'LoginService', loginController]);

    function loginController($http, $window, $scope, $rootScope, $mdUtil, $interval, configService, LoginService) {
        var vm = this;
        configService.getConfig(function (config) {
            vm.totpsecret = config.totpsecret;
        });
        vm.submit = submit;
        vm.register = register;
        vm.clearValidation = clearValidation;
        vm.validateUsername = validateUsername;
        vm.registerUser = registerUser;
        vm.resetValidation = resetValidation;
        vm.loginInfoValid = loginInfoValid;
        vm.basicInfoValid = basicInfoValid;
        vm.signupInfoValid = signupInfoValid;
        vm.moveOn = moveOn;
        vm.externalSite = URI($window.location.href).search(true).client_id;
        $rootScope.registrationUrl = '/register' + $window.location.search;
        vm.logo = "";
        vm.twoFAMethod = 'sms';
        vm.login = "";
        vm.password = "";
        vm.validateUsername = $mdUtil.debounce(function () {
            $scope.loginform.registerlogin.$setValidity("duplicate_username", true);
            $scope.loginform.registerlogin.$setValidity("invalid_username_format", true);
            if ($scope.loginform.registerlogin.$valid) {
                validateUsername(vm.registerlogin)
                    .then(function (response) {
                        $scope.loginform.registerlogin.$setValidity(response.data.error, response.data.valid);
                    });
            }
        }, 500, true);

        var listener;
        activate();

        function activate() {
            if (vm.externalSite) {
                LoginService.getLogo(vm.externalSite).then(
                    function(data) {
                        vm.logo = data.logo;
                        renderLogo();
                    }
                );
                window.addEventListener('resize', resizeLogo, false);
                window.addEventListener('orientationchange', resizeLogo, false);
            }
            autoFillListener();
            $scope.$on('$destroy', function() {
                  // Make sure that the interval is destroyed too
                  stopListening();
            });
        }

        function renderLogo() {
            if (vm.logo !== "") {
                var img = new Image();
                img.onload = function() {
                    var c = document.getElementById("login-logo");
                    if (!c) {
                        console.log("aborting logo render - canvas not loaded");
                        return;
                    }
                    var ctx = c.getContext("2d");
                    ctx.clearRect(0, 0, c.width, c.height);
                    ctx.drawImage(img, 0, 0, c.width, c.height);
                }
                img.src = vm.logo;

            }
        }

        function autoFillListener() {
            listener = $interval(function() {
                var login = document.getElementById("username");
                var password = document.getElementById("password");
                if (login.value !== vm.login) {
                    vm.login = login.value;
                }
                if (password.value !== vm.password) {
                    vm.password = password.value;
                }
            }, 100);
        }

        function stopListening() {
            if (angular.isDefined(listener)) {
                $interval.cancel(listener);
                listener = undefined;
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

        function register() {
          var redirectparams = $window.location.search.replace('?', '');
              registerUser(vm.twoFAMethod, vm.registerlogin, vm.registeremail, vm.registerpassword, vm.totpcode, vm.sms, redirectparams)
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
                                  $scope.loginform.phonenumber.$setValidity(err, false);
                                  break;
                              case 'invalid_totpcode':
                                  $scope.loginform.totpcode.$setValidity(err, false);
                                  break;
                              case 'invalid_password':
                                  vm.registration2fa = false;
                                  $scope.loginform.registerpassword.$setValidity(err, false);
                                  break;
                              case 'invalid_username_format':
                                  vm.registration2fa = false;
                                  $scope.loginform.registerlogin.$setValidity(err, false);
                                  break;
                              default:
                                  console.error('Unconfigured error:', response.data.error);
                          }
                          break;
                      case 409:
                          vm.registration2fa = false;
                          $scope.loginform.registerlogin.$setValidity('duplicate_username', false);
                          break;
                  }
              });
        }

        function clearValidation() {
            $scope.loginform.password.$setValidity("invalidcredentials", true);
        }

        function validateUsername(username) {
            var options = {
                params: {
                    username: username
                }
            };
            return $http.get('/validateusername', options);
        }

        function registerUser(twoFAMethod, login, email, password, totpcode, sms, redirectparams) {
            var url = '/register';
            var data = {
                twofamethod: twoFAMethod,
                login: login.trim(),
                email: email.trim(),
                password: password,
                totpcode: totpcode,
                phonenumber: sms,
                redirectparams: redirectparams
            };
            return $http.post(url, data);
        }

        function resetValidation(prop) {
            switch (prop) {
                case 'phonenumber':
                    $scope.loginform[prop].$setValidity("invalid_phonenumber", true);
                    break;
                case 'totpcode':
                    $scope.loginform[prop].$setValidity("invalid_totpcode", true);
                    break;
                case 'twoFAMethod':
                    $scope.loginform.phonenumber.$setValidity("invalid_phonenumber", true);
                    $scope.loginform.phonenumber.$setValidity("pattern", true);
                    $scope.loginform.totpcode.$setValidity("totpcode", true);
                    break;
            }
        }

        function loginInfoValid() {
            return $scope.loginform.username
                && $scope.loginform.username.$valid
                && $scope.loginform.password.$valid;
        }

        function basicInfoValid() {
            return $scope.loginform.registerlogin
                && $scope.loginform.registerlogin.$valid
                && $scope.loginform.registeremail.$valid
                && $scope.loginform.registerpassword.$valid
                && $scope.loginform.passwordvalidation.$valid;
        }

        function signupInfoValid() {
            switch (vm.twoFAMethod) {
                case 'sms':
                    return basicInfoValid() && $scope.loginform.phonenumber.$valid;
                    break;
                case 'totp':
                    return basicInfoValid() && $scope.loginform.totpcode.$valid;
                    break;
            }
        }

        function moveOn() {
            vm.registration2fa = true;
        }

        function resizeLogo(e) {
            var formArea = document.getElementById("form-area");
            var logoArea = document.getElementById("login-logo");
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
    }
})();
