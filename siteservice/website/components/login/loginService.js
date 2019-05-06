(function () {
    'use strict';

    angular
        .module("loginApp")
        .service("LoginService", LoginService);

    LoginService.$inject = ['$http', '$q'];


    function LoginService($http, $q) {
        var apiURL = '/login';

        return {
            getTwoFactorAuthenticationMethods: getTwoFactorAuthenticationMethods,
            sendSmsCode: sendSmsCode,
            submitTotpCode: submitTotpCode,
            submitSmsCode: submitSmsCode,
            checkSmsConfirmation: checkSmsConfirmation,
            getLogo: getLogo,
            getDescription: getDescription,
            getTwoFASettings: getTwoFASettings,
            loginNoTwoFA: loginNoTwoFA,
            checkRequiresTwoFA:checkRequiresTwoFA
        };

        function genericHttpCall(httpFunction, url, data) {
            if (data) {
                return httpFunction(url, data)
                    .then(
                        function (response) {
                            return response.data;
                        },
                        function (reason) {
                            return $q.reject(reason);
                        }
                    );
            }
            else {
                return httpFunction(url)
                    .then(
                        function (response) {
                            return response.data;
                        },
                        function (reason) {
                            return $q.reject(reason);
                        }
                    );
            }
        }

        function getTwoFactorAuthenticationMethods() {
            var url = apiURL + '/twofamethods';
            return genericHttpCall($http.get, url);
        }

        function loginNoTwoFA() {
            var url = apiURL + '/loginno2fa';
            return genericHttpCall($http.post, url);
        }

        function getTwoFASettings(){
            var url = apiURL + '/twofasettings';
            return genericHttpCall($http.get, url);
        }

        function sendSmsCode(phoneLabel) {
            var url = apiURL + '/smscode/' + encodeURIComponent(phoneLabel);
            var data = {
                langkey: localStorage.getItem('langKey')
            };
            return genericHttpCall($http.post, url, data);
        }

        function submitTotpCode(code, queryString) {
            var url = apiURL + '/totpconfirmation' + queryString;
            var data = {
                totpcode: code
            };
            return genericHttpCall($http.post, url, data);
        }

        function submitSmsCode(code, queryString) {
            var url = apiURL + '/smsconfirmation' + queryString;
            var data = {
                smscode: code
            };
            return genericHttpCall($http.post, url, data);
        }

        function checkSmsConfirmation() {
            var url = apiURL + '/smsconfirmed';
            return genericHttpCall($http.get, url);
        }

        function getLogo(globalid) {
            var url = '/api/organizations/' + encodeURIComponent(globalid) + '/logo';
            return genericHttpCall($http.get, url);
        }

        function getDescription(globalId, langKey) {
            var url = '/api/organizations/' + encodeURIComponent(globalId) + '/description/' + encodeURIComponent(langKey) + '/withfallback';
            return genericHttpCall($http.get, url);
        }

        function checkRequiresTwoFA(globalId) {
            var url = 'api/organizations' + '/' + encodeURIComponent(globalId) + '/2fa/isrequired';
            return genericHttpCall($http.get, url);
        }
    }
})();
