/**
 * Created by lucas on 27/05/16.
 */
(function () {
    'use strict';
    angular.module('itsyouonlineApp')
        .directive('authorizationDetails', function () {
            return {
                restrict: 'AE',
                templateUrl: 'components/user/directives/authorizationDetails.html',
                link: function (scope, element, attr) {
                    scope.save = save;
                    scope.getAuthorizationByLabel = getAuthorizationByLabel;
                    function getAuthorizationByLabel(property, requestedLabel) {
                        return scope.authorizations[property].filter(function (val) {
                            return val.requestedlabel === requestedLabel;
                        })[0];
                    }

                    function save() {
                        scope.authorizations.organizations = [];
                        angular.forEach(scope.requested.organizations, function (allowed, organization) {
                            if (allowed) {
                                scope.authorizations.organizations.push(organization);
                            }
                        });
                        // Filter unauthorized permission labels
                        angular.forEach(scope.authorizations, function (value, key) {
                            if (Array.isArray(value)) {
                                angular.forEach(value, function (val, i) {
                                    if (!val || val.reallabel === '') {
                                        value.splice(value.indexOf(val), 1);
                                    }
                                });
                            }
                        });
                        scope.update();
                    }
                }
            };
        });
})();