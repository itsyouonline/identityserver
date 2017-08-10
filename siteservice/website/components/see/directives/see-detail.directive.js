(function () {
    'use strict';
    angular.module('itsyouonlineApp')
        .directive('seeDetail', [function () {
            return {
                restrict: 'AE',
                templateUrl: 'components/see/directives/see-detail.directive.html',
                scope: {
                    see: '=',
                    detailed: '=',
                    showFullHistory: '=',
                    toggleFullHistory: '&'
                }
            };
        }]);
})();
