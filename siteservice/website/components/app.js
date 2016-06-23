(function () {
    'use strict';
    angular
        .module('itsyouonlineApp', ['ngCookies', 'ngMaterial', 'ngRoute', 'ngMessages', 'ngSanitize',
            'monospaced.qrcode',
            'itsyouonline.shared', 'itsyouonline.header', 'itsyouonline.footer', 'itsyouonline.user'])
        .config(['$mdThemingProvider', themingConfig])
        .config(['$httpProvider', httpConfig])
        .config(['$routeProvider', routeConfig])
        .factory('authenticationInterceptor', ['$q', '$window', authenticationInterceptor])
        .directive('pagetitle', ['$rootScope', '$timeout', pagetitle])
        .run(['$route', '$cookies', '$rootScope', '$location', runFunction]);

    function themingConfig($mdThemingProvider) {
        $mdThemingProvider.definePalette('blueish', {
            '50': '#f7fbfd',
            '100': '#badeed',
            '200': '#8ec8e2',
            '300': '#55add3',
            '400': '#3ca1cd',
            '500': '#3091bb',
            '600': '#2a7ea3',
            '700': '#236b8a',
            '800': '#1d5872',
            '900': '#17455a',
            'A100': '#f7fbfd',
            'A200': '#badeed',
            'A400': '#3ca1cd',
            'A700': '#236b8a',
            'contrastDefaultColor': 'light',
            'contrastDarkColors': '50 100 200 300 400 A100 A200 A400'
        });
        $mdThemingProvider.theme('default')
            .primaryPalette('blueish');
    }

    function httpConfig($httpProvider) {
        $httpProvider.interceptors.push('authenticationInterceptor');
        //initialize get if not there
        if (!$httpProvider.defaults.headers.get) {
            $httpProvider.defaults.headers.get = {};
        }
        //disable IE ajax request caching
        $httpProvider.defaults.headers.get['If-Modified-Since'] = '0';
    }

    function authenticationInterceptor($q, $window) {
        return {
            'request': function (config) {
                return config || $q.when(config);
            },
            'response': function (response) {
                return response || $q.when(response);
            },

            'responseError': function (rejection) {
                if (rejection.status == 401 || rejection.status == 403 || rejection.status == 419) {
                    $window.location.href = "";
                }

                return $q.reject(rejection);
            }
        };
    }

    function routeConfig($routeProvider) {
        $routeProvider
            .when('/home/:tab?', {
                templateUrl: 'components/user/views/home.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Home'
                }
            })
            .when('/authorize', {
                templateUrl: 'components/user/views/authorize.html',
                controller: 'AuthorizeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Authorize'
                }
            })
            .when('/company/new', {
                templateUrl: 'components/company/views/new.html',
                controller: 'CompanyController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'New company'
                }
            })
            .when('/organization/:globalid', {
                templateUrl: 'components/organization/views/detail.html',
                controller: 'OrganizationDetailController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Organization detail'
                }
            })
            .otherwise('/home');
    }

    function pagetitle($rootScope, $timeout) {
        return {
            link: function (scope, element) {
                var listener = function (event, current, previous) {
                    var pageTitle = 'It\'s You Online';
                    if (current.$$route && current.$$route.data && current.$$route.data.pageTitle) {
                        pageTitle = current.$$route.data.pageTitle + ' - ' + pageTitle;
                    }
                    $timeout(function () {
                        element.text(pageTitle);
                    }, 0, false);
                };

                $rootScope.$on('$routeChangeSuccess', listener);
            }
        };
    }

    function runFunction($route, $cookies, $rootScope, $location) {
        $rootScope.user = $cookies.get('itsyou.online.user');
        var original = $location.path;
        // prevent controller reload when changing route params in code because we aren't using states
        $location.path = function (path, reload) {
            if (reload === false) {
                var lastRoute = $route.current;
                var un = $rootScope.$on('$locationChangeSuccess', function () {
                    $route.current = lastRoute;
                    un();
                });
            }
            return original.apply($location, [path]);
        };
        if (!String.prototype.startsWith) {
            String.prototype.startsWith = function (hashstack, needle) {
                return haystack.lastIndexOf(needle, 0) === 0;
            };
        }
    }
})();
