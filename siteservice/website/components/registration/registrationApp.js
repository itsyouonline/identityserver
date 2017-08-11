(function () {
    'use strict';
    angular
        .module('itsyouonline.registration', [
            'ngMaterial', 'ngMessages', 'ngRoute', 'ngCookies', 'ngSanitize', 'monospaced.qrcode',
            'itsyouonline.shared', 'itsyouonline.header', 'itsyouonline.footer', 'itsyouonline.validation',
            'itsyouonline.telinput', 'pascalprecht.translate'])
        .config(['$mdThemingProvider', themingConfig])
        .config(['$httpProvider', httpConfig])
        .config(['$routeProvider', routeConfig])
        .config(['$translateProvider', translateConfig])
        .config(['$mdAriaProvider', function ($mdAriaProvider) {
            $mdAriaProvider.disableWarnings();
        }])
        .factory('authenticationInterceptor', ['$q', '$window', authenticationInterceptor]);


    function routeConfig($routeProvider) {
        $routeProvider
            .when('/', {
                templateUrl: 'components/registration/views/registrationform.html',
                controller: 'registrationController',
                controllerAs: 'vm'
            })
            .when('/smsconfirmation', {
                templateUrl: 'components/registration/views/registrationsmsform.html',
                controller: 'smsController',
                controllerAs: 'vm'
            })
            .when('/resendsms', {
                templateUrl: 'components/registration/views/registrationresendsms.html',
                controller: 'resendSmsController',
                controllerAs: 'vm'
            })
            .otherwise('/');
    }


    initializePolyfills();

})();
