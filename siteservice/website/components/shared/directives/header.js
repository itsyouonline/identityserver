(function () {
    'use strict';
    angular.module('itsyouonline.header', ['pascalprecht.translate'])
        .directive('itsYouOnlineHeader', ['$location', '$window', '$translate', function ($location, $window, $translate) {
            return {
                restrict: 'E',
                replace: true,
                templateUrl: 'components/shared/directives/header.html',
                link: function (scope, element, attr) {
                    scope.header_login = attr.register !== undefined;
                    scope.header_registration = attr.login !== undefined;
                    scope.showCookieWarning = !localStorage.getItem('cookiewarning-dismissed');
                    scope.hideCookieWarning  = hideCookieWarning;
                    scope.updateLanguage = updateLanguage;
                    scope.setLanguage = setLanguage;
                    scope.pushClick = pushClick;
                    scope.toggleMenu = toggleMenu;
                    init();

                    function init() {
                        scope.langKey = localStorage.getItem('langKey');
                    }

                    function hideCookieWarning(){
                        localStorage.setItem('cookiewarning-dismissed', true);
                        scope.showCookieWarning = false;
                    }

                    function updateLanguage(){
                        localStorage.setItem('langKey', scope.langKey);
                        localStorage.setItem('selectedLangKey', scope.langKey);
                        $translate.use(scope.langKey);
                    }

                    function setLanguage(lang) {
                        scope.langKey = lang;
                        scope.langSelect = !scope.langSelect;
                        scope.toggleNavMenu = !scope.toggleNavMenu;
                        updateLanguage();
                    }

                    function pushClick(url) {
                        $location.path(url);
                        scope.toggleNavMenu = !scope.toggleNavMenu;
                    }

                    function toggleMenu() {
                        scope.toggleNavMenu = !scope.toggleNavMenu;
                        scope.langSelect = false;
                    }
                }
            };
        }]);
})();
