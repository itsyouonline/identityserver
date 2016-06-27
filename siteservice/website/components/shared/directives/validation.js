(function () {
    'use strict';
    angular
        .module('itsyouonline.validation', [])
        .directive('passwordValidation', passwordValidation);

    function passwordValidation() {
        return {
            require: 'ngModel',
            scope: {
                firstValue: '=passwordValidation'
            },
            link: function (scope, element, attributes, ngModel) {
                ngModel.$validators.passwordIdentical = function (modelValue) {
                    return modelValue === scope.firstValue;
                };

                scope.$watch('firstValue', function () {
                    ngModel.$validate();
                });
            }
        };
    }

})();