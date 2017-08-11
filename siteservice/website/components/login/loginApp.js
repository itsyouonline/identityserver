(function () {
    'use strict';
    angular.module('loginApp', ['ngMaterial', 'ngCookies', 'ngMessages', 'ngRoute', 'ngSanitize', 'monospaced.qrcode', 'itsyouonline.shared',
        'itsyouonline.header', 'itsyouonline.footer', 'itsyouonline.user', 'itsyouonline.validation',
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
                templateUrl: 'components/login/views/loginform.html',
                controller: 'loginController',
                controllerAs: 'vm'
            })
            .when('/2fa', {
                templateUrl: 'components/login/views/twoFactorAuthentication.html',
                controller: 'twoFactorAuthenticationController',
                controllerAs: 'vm'
            })
            .when('/forgotpassword', {
                templateUrl: 'components/login/views/forgotpassword.html',
                controller: 'forgotPasswordController',
                controllerAs: 'vm'
            })
            .when('/resetpassword/:code', {
                templateUrl: 'components/login/views/resetpassword.html',
                controller: 'resetPasswordController',
                controllerAs: 'vm'
            })
            .when('/resendsms', {
                templateUrl: 'components/registration/views/registrationresendsms.html',
                controller: 'resendSmsController',
                controllerAs: 'vm'
            })
            .when('/smsconfirmation', {
                templateUrl: 'components/registration/views/registrationsmsform.html',
                controller: 'smsConfirmationController',
                controllerAs: 'vm'
            }).when('/organizationinvite/:code', {
                templateUrl: 'components/login/views/organizationinvite.html',
                controller: 'organizationInviteController',
                controllerAs: 'vm'
            }).when('/validateemail', {
                templateUrl: 'components/login/views/validateEmailAddress.html',
                controller: 'validateEmailController',
                controllerAs: 'vm'
            })
            .otherwise('/');
    }

    initializePolyfills();

})();
