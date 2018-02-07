(function() {
    'use strict';
    angular.module('itsyouonline.telinput', ['pascalprecht.translate'])
        .directive('telinput', ['countries', '$translate', '$http', '$q', function(countries, $translate, $http, $q) {
            return {

                restrict: 'E',
                replace: true,
                scope: {
                  number: '=phonenumber',
                  validationerrors: '=?error'
                },
                templateUrl: 'components/shared/directives/telinput.html',
                link: function (scope, element, attr) {
                    scope.countries = countries;
                    scope.country = {};
                    scope.sms = "";
                    scope.prefCountry = "";

                    scope.updateSMS = updateSMS;
                    scope.isNumeric = isNumeric;

                    init();

                    function init() {
                        if (!scope.validationerrors) {
                            scope.validationerrors = {};
                        }
                        // check if a valid international number is prefilled
                        if (scope.number) {
                            for (var i = 0; i < countries.length; i++) {
                                if (scope.number.startsWith(countries[i].dial_code)){
                                    scope.country = countries[i];
                                    scope.sms = scope.number.slice(countries[i].dial_code.length);
                                    return;
                                }
                            }
                        }
                        // if there is no valid number get the location
                        getLocation().then(
                            function(data) {
                                scope.prefCountry = data.location;
                                fillPrefCountry();
                            }
                        );
                    }

                    // try to fill in te preffered country
                    function fillPrefCountry() {
                        if (scope.prefCountry && !scope.number) {
                            for (var i = 0; i < countries.length; i++) {
                                if (countries[i].code == scope.prefCountry){
                                    scope.country = countries[i];
                                    break;
                                }
                            }
                        }
                    }

                    function updateSMS() {
                        scope.validationerrors.pattern = false;
                        scope.validationerrors['phone_already_used'] = false;
                        scope.validationerrors['invalid_phonenumber'] = false;
                        scope.validationerrors['illegalactionforcurrentstate'] = false;
                        var phone = scope.sms;
                        if (phone.startsWith("0")) {
                            phone = phone.substring(1);
                        }
                        scope.number = scope.country.dial_code + phone;
                        if (!isNumeric(scope.number)) {
                            scope.validationerrors.pattern = true;
                        }
                    }

                    function isNumeric(n) {
                        return !isNaN(parseFloat(n)) && isFinite(n);
                    }

                    function getLocation() {
                        var url = '/location';
                        return $http.get(url).then(
                            function (response) {
                                return response.data;
                            },
                            function (reason) {
                                return $q.reject(reason);
                            }
                        );
                    }
                }
            }
        }])
})();
