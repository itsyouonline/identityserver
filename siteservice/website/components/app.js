(function () {
    'use strict';
    angular
        .module('itsyouonlineApp', ['ngCookies', 'ngMaterial', 'ngRoute', 'ngMessages', 'ngSanitize',
            'monospaced.qrcode',
            'itsyouonline.shared', 'itsyouonline.header', 'itsyouonline.footer', 'itsyouonline.user',
            'itsyouonline.validation', 'itsyouonline.telinput', 'pascalprecht.translate'])
        .config(['$mdThemingProvider', themingConfig])
        .config(['$httpProvider', httpConfig])
        .config(['$routeProvider', routeConfig])
        .config(['$translateProvider', translateConfig])
        .config([init])
        .factory('authenticationInterceptor', ['$q', '$window', authenticationInterceptor])
        .directive('pagetitle', ['$rootScope', '$timeout', 'footerService', pagetitle])
        .run(['$route', '$cookies', '$rootScope', '$location', runFunction]);

    function routeConfig($routeProvider) {
        $routeProvider
            .when('/authorize', {
                templateUrl: 'components/user/views/authorize.html',
                controller: 'AuthorizeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Authorize',
                    showFooter: false
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
            .when('/profile', {
                templateUrl: 'components/user/views/profile.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Profile'
                }
            })
            .when('/notifications', {
                templateUrl: 'components/user/views/notifications.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Notifications'
                }
            })
            .when('/organizations', {
                templateUrl: 'components/user/views/organizations.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Organizations'
                }
            })
            .when('/authorizations', {
                templateUrl: 'components/user/views/authorizations.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Authorizations'
                }
            })
            .when('/see', {
                templateUrl: 'components/user/views/see.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'See'
                }
            })
            .when('/see/:globalid/:uniqueid/detail', {
                templateUrl: 'components/user/views/seeDetail.html',
                controller: 'SeeDetailController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'See detail'
                }
            })
            .when('/settings', {
                templateUrl: 'components/user/views/settings.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Settings'
                }
            })
            .otherwise('/profile');
    }

    function init() {
        localStorage.setItem('hasLoggedIn', true);
    }

    function pagetitle($rootScope, $timeout, footerService) {
        return {
            link: function (scope, element) {
                var listener = function (event, current) {
                    var pageTitle = 'It\'s You Online';
                    var routeData = current.$$route && current.$$route.data || {};
                    if (routeData.pageTitle) {
                        pageTitle = current.$$route.data.pageTitle + ' - ' + pageTitle;
                    }
                    footerService.setFooter(routeData.showFooter !== undefined ? routeData.showFooter : true);
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
        if (window.location.hostname === 'dev.itsyou.online') {
            setTimeout(function () {
                window.location.reload();
            }, 9 * 60 * 1000);
        }
        initializePolyfills();
    }

})();
