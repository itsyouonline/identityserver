(function () {
    'use strict';
    angular
        .module('itsyouonlineApp', ['ngCookies', 'ngMaterial', 'ngRoute', 'ngMessages', 'ngSanitize',
            'monospaced.qrcode', 'ui.router',
            'itsyouonline.shared', 'itsyouonline.header', 'itsyouonline.footer', 'itsyouonline.user',
            'itsyouonline.validation', 'itsyouonline.telinput', 'pascalprecht.translate', 'btford.markdown'])
        .config(['$mdThemingProvider', themingConfig])
        .config(['$httpProvider', httpConfig])
        .config(['$stateProvider', '$urlRouterProvider', stateConfig])
        .config(['$translateProvider', translateConfig])
        .config([init])
        .factory('authenticationInterceptor', ['$q', '$window', authenticationInterceptor])
        .directive('pagetitle', ['$rootScope', '$timeout', 'footerService', pagetitle])
        .run(['$rootScope', '$cookies', '$window', 'UserService', runFunction]);

    function stateConfig($stateProvider, $urlRouterProvider) {
        $urlRouterProvider.otherwise('/profile');
        $stateProvider.state('authorize', {
            url: '/authorize',
            templateUrl: 'components/user/views/authorize.html',
            controller: 'AuthorizeController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Authorize',
                showFooter: false
            }
        })
        .state('/company/new', {
            url: '/company/new',
            templateUrl: 'components/company/views/new.html',
            controller: 'CompanyController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'New company'
            }
        })
        .state('organization', {
            url: '/organization/:globalid',
            templateUrl: 'components/organization/views/detailTabs.html',
            controller: 'OrganizationDetailController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Organization detail'
            }
        })
        .state('organization.people', {
            url: '/people',
            templateUrl: 'components/organization/views/detailTabsPeople.html',
            params: {
                pageTitle: 'Organization people'
            }
        })
        .state('organization.structure', {
            url: '/structure',
            templateUrl: 'components/organization/views/detailTabsStructure.html',
            params: {
                pageTitle: 'Organization stucture'
            }
        })
        .state('organization.settings', {
            url: '/settings',
            templateUrl: 'components/organization/views/detailTabsSettings.html',
            params: {
                pageTitle: 'Organization settings'
            }
        })
        .state('profile', {
            url: '/profile',
            templateUrl: 'components/user/views/profile.html',
            controller: 'UserHomeController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Profile'
            }
        })
        .state('notifications', {
            url: '/notifications',
            templateUrl: 'components/user/views/notifications.html',
            controller: 'UserHomeController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Notifications'
            }
        })
        .state('organizations', {
            url: '/organizations',
            templateUrl: 'components/user/views/organizations.html',
            controller: 'UserHomeController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Organizations'
            }
        })
        .state('authorizations', {
            url: '/authorizations',
            templateUrl: 'components/user/views/authorizations.html',
            controller: 'UserHomeController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Authorizations'
            }
        })
        .state('settings', {
            url: '/settings',
            templateUrl: 'components/user/views/settings.html',
            controller: 'UserHomeController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'Settings'
            }
        })
        .state('see', {
            url: '/see',
            templateUrl: 'components/see/views/see-list.html',
            controller: 'SeeListController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'See'
            }
        })
        .state('seeDetail', {
            url: '/see/:uniqueid',
            templateUrl: 'components/see/views/see-detail.html',
            controller: 'SeeDetailController',
            controllerAs: 'vm',
            params: {
                pageTitle: 'See detail'
            }
        });
    }

    function init() {
        localStorage.setItem('hasLoggedIn', true);
    }

    function pagetitle($rootScope, $timeout, footerService) {
        return {
            link: function (scope, element) {
                var listener = function (event, current) {
                    var pageTitle = 'It\'s You Online';
                    var routeData = current && current.params || {};
                    if (routeData.pageTitle) {
                        pageTitle = routeData.pageTitle + ' - ' + pageTitle;
                    }
                    footerService.setFooter(routeData.showFooter !== undefined ? routeData.showFooter : true);
                    $timeout(function () {
                        element.text(pageTitle);
                    }, 0, false);
                };

                $rootScope.$on('$stateChangeSuccess', listener);
            }
        };
    }

    function runFunction($rootScope, $cookies, $window, UserService) {
        $rootScope.user = $cookies.get('itsyou.online.user');
        UserService.setUsername($rootScope.user);
        if ($window.location.hostname === 'dev.itsyou.online') {
            setTimeout(function () {
                $window.location.reload();
            }, 9 * 60 * 1000);
        }
        initializePolyfills();
    }

})();
