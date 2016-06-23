(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("AuthorizeController", AuthorizeController);


    AuthorizeController.$inject = ['$scope', '$rootScope', '$location', '$window', 'UserService', 'UserDialogService'];

    function AuthorizeController($scope, $rootScope, $location, $window, UserService, UserDialogService) {
        var vm = this;

        var queryParams = $location.search();
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requestedorganizations = [];
        vm.username = $rootScope.user;

        vm.user = {};

        UserDialogService.init(vm);
        vm.showEmailDialog = addEmail;
        vm.showPhonenumberDialog = addPhone;
        vm.showAddressDialog = addAddress;
        vm.showBankAccountDialog = bank;
        vm.showDigitalWalletAddressDialog = digitalWalletAddress;

        var properties = ['addresses', 'emailaddresses', 'phonenumbers', 'bankaccounts', 'digitalwallet'];
        $scope.requested = {
            organizations: {},
            facebook: false,
            github: false
        };
        $scope.authorizations = {
            organizations: {}
        };
        angular.forEach(properties, function (prop) {
            $scope.requested[prop] = [];
            $scope.authorizations[prop] = [];
        });

        $scope.update = update;
        $scope.isNew = true;

        activate();

        function activate() {
            fetch();
        }

        function fetch() {

            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.user = data;
                        parseScopes();
                    },
                    function(reason) {
                        $window.location.href = 'error' + reason.status;
                    }
                );
        }

        function parseScopes() {
            if (vm.requestedScopes) {
                var scopes = vm.requestedScopes.split(',');
                // Filter duplicated scopes
                scopes = scopes.filter(function (item, pos, self) {
                    return self.indexOf(item) === pos;
                });
                angular.forEach(scopes, function (scope, i) {
                    var splitPermission = scope.split(':');
                    // Empty label -> 'main'
                    var permissionLabel = splitPermission.length > 2 && splitPermission[2] ? splitPermission[2] : 'main';
                    if (scope === 'user:name') {
                        $scope.requested.name = true;
                        $scope.authorizations.name = true;
                    }
                    if (scope.startsWith('user:memberof:')) {
                        $scope.requested.organizations[permissionLabel] = true;
                    }
                    else if (scope.startsWith('user:address:')) {
                        $scope.requested.addresses.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:email:')) {
                        $scope.requested.emailaddresses.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:phone:')) {
                        $scope.requested.phonenumbers.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:bankaccount:')) {
                        $scope.requested.bankaccounts.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:digitalwalletaddress:')) {
                        $scope.requested.digitalwallet.push({
                            label: permissionLabel,
                            currency: splitPermission.length === 4 ? splitPermission[3] : ''
                        });
                    }
                    else if (scope === 'user:github') {
                        $scope.requested.github = true;
                        $scope.authorizations.github = true;
                    }
                    else if (scope === 'user:facebook') {
                        $scope.requested.facebook = true;
                        $scope.authorizations.facebook = true;
                    }
                });
                angular.forEach($scope.requested, function (value, property) {
                    if (properties.indexOf(property) === -1) {
                        return;
                    }
                    // loop over requests
                    angular.forEach(value, function (requestedLabel) {
                        // select first by default, None if the user did not configure this property yet
                        if (typeof requestedLabel === 'object') {
                            requestedLabel = requestedLabel.label;
                        }
                        var authorization = {
                            requestedlabel: requestedLabel,
                            reallabel: vm.user[property].length ? vm.user[property][0].label : ''
                        };
                        $scope.authorizations[property].push(authorization);
                    });
                });
            }
        }

        function update() {
            // called by the authorizationDetailsDirective
            $scope.authorizations.username = vm.username;
            $scope.authorizations.grantedTo = vm.requestingorganization;
            UserService
                .saveAuthorization($scope.authorizations)
                .then(
                    function (data) {
                        var u = URI($location.absUrl());
                        var endpoint = queryParams["endpoint"];
                        delete queryParams.endpoint;
                        u.pathname(endpoint);
                        u.search(queryParams);
                        $window.location.href = u.toString();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function addEmail(event, label) {
            selectDefault(UserDialogService.emailDetail, event, label, 'emailaddresses');
        }

        function addPhone(event, label) {
            selectDefault(UserDialogService.phonenumberDetail, event, label, 'phonenumbers');
        }

        function addAddress(event, label) {
            selectDefault(UserDialogService.addressDetail, event, label, 'addresses');
        }

        function bank(event, label) {
            selectDefault(UserDialogService.bankAccount, event, label, 'bankaccounts');
        }

        function digitalWalletAddress(event, label) {
            selectDefault(UserDialogService.digitalWalletAddressDetail, event, label, 'digitalwallet');
        }

        function selectDefault(fx, event, label, property) {
            fx(event).then(function (data) {
                $scope.getAuthorizationByLabel(property, label).reallabel = data.data.label;

            }, function () {
                // Select first possible value, else 'None'
                $scope.getAuthorizationByLabel(property, label).reallabel = vm.user[property][0] ? vm.user[property][0].label : '';
            });
        }
    }
})();
