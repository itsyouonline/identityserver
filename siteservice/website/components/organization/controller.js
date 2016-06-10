(function() {
    'use strict';
    angular
        .module("itsyouonlineApp")
        .controller("OrganizationController", OrganizationController)
        .controller("OrganizationDetailController", OrganizationDetailController)
        .controller("InvitationDialogController", InvitationDialogController);


    OrganizationController.$inject = ['$rootScope', '$location', '$routeParams', 'OrganizationService', '$window', '$scope'];
    OrganizationDetailController.$inject = ['$routeParams', '$window', 'OrganizationService', '$mdDialog', '$mdMedia', '$rootScope'];

    function OrganizationController($rootScope, $location, $routeParams, OrganizationService, $window, $scope) {
        var vm = this;
        vm.create = create;
        vm.childOrganizationNames = [];
        vm.parentOrganization = $routeParams.globalid;

        vm.username = $rootScope.user;
        vm.clearErrors = clearErrors;
        vm.getOrganizationDisplayname = getOrganizationDisplayname;
        activate();

        function activate() {
            vm.childOrganizationNames = getChildOrganizations(vm.parentOrganization);
        }

        function create(){
            if (!$scope.newOrganizationForm.$valid) {
                return;
            }
            var dns = [];

            if (vm.dns) {
                dns.push(vm.dns);
            }

            OrganizationService
                .create(vm.name, dns, vm.username, vm.parentOrganization)
                .then(
                    function(data){
                        $location.path("/organizations/" + data.globalid);
                    },
                    function(reason){
                        if (reason.status == 409) {
                            $scope.newOrganizationForm.name.$setValidity('duplicate', false);
                        }
                        else{
                            $window.location.href = "error" + reason.status;
                        }
                    }
                );
        }

        function clearErrors() {
            $scope.newOrganizationForm.name.$setValidity('duplicate', true);
        }
    }

    function OrganizationDetailController($routeParams, $window, OrganizationService, $mdDialog, $mdMedia, $rootScope) {
        var vm = this,
            globalid = $routeParams.globalid;
        vm.invitations = [];
        vm.apikeylabels = [];
        vm.organization = {};
        vm.organizationRoot = {};
        vm.userDetails = {};
        vm.hasEditPermission = false;
        vm.childOrganizationNames = [];

        vm.showInvitationDialog = showInvitationDialog;
        vm.showAPIKeyCreationDialog = showAPIKeyCreationDialog;
        vm.showAPIKeyDialog = showAPIKeyDialog;
        vm.showDNSDialog = showDNSDialog;
        vm.getOrganizationDisplayname = getOrganizationDisplayname;
        vm.fetchInvitations = fetchInvitations;
        vm.fetchAPIKeyLabels = fetchAPIKeyLabels;
        activate();

        function activate() {
            fetch();
        }

        function fetch(){
            OrganizationService
                .get(globalid)
                .then(
                    function(data) {
                        vm.organization = data;
                        vm.childOrganizationNames = getChildOrganizations(vm.organization.globalid);
                        vm.hasEditPermission = vm.organization.owners.indexOf($rootScope.user) !== -1;
                        fetchInvitations();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );

            OrganizationService.getOrganizationTree(globalid)
                .then(function (data) {
                    vm.organizationRoot.children = [];
                    vm.organizationRoot.children.push(data);
                }, function (error) {
                    $window.location.href = "error" + error.status;
                });
        }

        function fetchInvitations() {
            if (!vm.hasEditPermission || vm.invitations.length) {
                return;
            }
            OrganizationService
                .getInvitations(globalid)
                .then(
                    function (data) {
                        vm.invitations = data;
                    },
                    function (reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }


        function fetchAPIKeyLabels(){
            if (!vm.hasEditPermission || vm.apikeylabels.length) {
                return;
            }
            OrganizationService
                .getAPIKeyLabels(globalid)
                .then(
                    function(data) {
                        vm.apikeylabels = data;
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function showInvitationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: InvitationDialogController,
                templateUrl: 'components/organization/views/invitationdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window
                    }
            })
            .then(
                function(invitation) {
                    vm.invitations.push(invitation);
                });
        }




        function showAPIKeyCreationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: APIKeyDialogController,
                templateUrl: 'components/organization/views/apikeydialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window,
                        label: ""
                    }
            })
            .then(
                function(data) {
                    if (data.newLabel) {
                        vm.apikeylabels.push(data.newLabel);
                    }
                });
        }


        function showAPIKeyDialog(ev, label) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: APIKeyDialogController,
                templateUrl: 'components/organization/views/apikeydialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window,
                        label: label
                    }
            })
            .then(
                function(data) {
                    if (data.originalLabel != data.newLabel){
                        if (data.originalLabel) {
                            if (data.newLabel){
                                //replace
                                vm.apikeylabels[vm.apikeylabels.indexOf(data.originalLabel)] = data.newLabel;
                            }
                            else {
                                //remove
                                vm.apikeylabels.splice(vm.apikeylabels.indexOf(data.originalLabel),1);
                            }
                        }
                        else {
                            //add
                            vm.apikeylabels.push(data.newLabel);
                        }
                    }
                });
        }


        function showDNSDialog(ev, dnsName) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: DNSDialogController,
                templateUrl: 'components/organization/views/dnsDialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals: {
                    OrganizationService: OrganizationService,
                    organization: vm.organization.globalid,
                    $window: $window,
                    dnsName: dnsName
                }
            })
                .then(
                    function (data) {
                        if (data.originalDns) {
                            vm.organization.dns.splice(vm.organization.dns.indexOf(data.originalDns), 1);
                        }
                        if (data.newDns) {
                            vm.organization.dns.push(data.newDns);
                        }
                    });
        }
    }

    function getOrganizationDisplayname(globalid) {
        if (globalid) {
            var splitted = globalid.split('.');
            return splitted[splitted.length - 1];
        }
    }

    function InvitationDialogController($scope, $mdDialog, organization, OrganizationService, $window) {

        $scope.role = "member";

        $scope.cancel = cancel;
        $scope.invite = invite;
        $scope.validationerrors = {};


        function cancel(){
            $mdDialog.cancel();
        }

        function invite(username, role){
            $scope.validationerrors = {};
            OrganizationService.invite(organization, username, role).then(
                function(data){
                    $mdDialog.hide(data);
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else if (reason.status == 404){
                        $scope.validationerrors.nosuchuser = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );

        }
    }

    function APIKeyDialogController($scope, $mdDialog, organization, OrganizationService, $window, label) {
        //If there is a key, it is already saved, if not, this means that a new secret is being created.

        $scope.apikey = {secret: ""};

        if (label) {
            $scope.secret = "-- Loading --";
            OrganizationService.getAPIKey(organization, label).then(
                function(data){
                    $scope.apikey = data;
                },
                function(reason){
                    $window.location.href = "error" + reason.status;
                }
            );
        }

        $scope.originalLabel = label;
        $scope.savedLabel = label;
        $scope.label = label;
        $scope.organization = organization;

        $scope.cancel = cancel;
        $scope.validationerrors = {};
        $scope.create = create;
        $scope.update = update;
        $scope.deleteAPIKey = deleteAPIKey;

        $scope.modified = false;


        function cancel(){
            if ($scope.modified) {
                $mdDialog.hide({originalLabel: label, newLabel: $scope.label});
            }
            else {
                $mdDialog.cancel();
            }
        }

        function create(label, apiKey){
            $scope.validationerrors = {};
            apiKey.label = label;
            OrganizationService.createAPIKey(organization, apiKey).then(
                function(data){
                    $scope.modified = true;
                    $scope.apikey = data;
                    $scope.savedLabel = data.label;
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }

        function update(oldLabel, newLabel){
            $scope.validationerrors = {};
            OrganizationService.updateAPIKey(organization, oldLabel, newLabel, $scope.apikey).then(
                function(data){
                    $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel});
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }


        function deleteAPIKey(label){
            $scope.validationerrors = {};
            OrganizationService.deleteAPIKey(organization, label).then(
                function(data){
                    $mdDialog.hide({originalLabel: label, newLabel: ""});
                },
                function(reason){
                    $window.location.href = "error" + reason.status;
                }
            );
        }

    }

    function DNSDialogController($scope, $mdDialog, organization, OrganizationService, $window, dnsName) {
        $scope.organization = organization;
        $scope.dnsName = dnsName;
        $scope.newDnsName = dnsName;

        $scope.cancel = cancel;
        $scope.validationerrors = {};
        $scope.create = create;
        $scope.update = update;
        $scope.remove = remove;

        function cancel() {
            $mdDialog.cancel();
        }

        function create(dnsName) {
            if (!$scope.form.$valid) {
                return;
            }
            $scope.validationerrors = {};
            OrganizationService.createDNS(organization, dnsName).then(
                function (data) {
                    $mdDialog.hide({originalDns: "", newDns: data.name});
                },
                function (reason) {
                    if (reason.status == 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                    else {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }

        function update(oldDns, newDns) {
            if (!$scope.form.$valid) {
                return;
            }
            $scope.validationerrors = {};
            OrganizationService.updateDNS(organization, oldDns, newDns).then(
                function (data) {
                    $mdDialog.hide({originalDns: oldDns, newDns: data.name});
                },
                function (reason) {
                    if (reason.status == 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                    else {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }


        function remove(dnsName) {
            $scope.validationerrors = {};
            OrganizationService.deleteDNS(organization, dnsName).then(
                function () {
                    $mdDialog.hide({originalDns: dnsName, newDns: ""});
                },
                function (reason) {
                    $window.location.href = "error" + reason.status;
                }
            );
        }
    }

    function getChildOrganizations(organization) {
        var children = [];
        if (organization) {
            for (var i = 0, splitted = organization.split('.'); i < splitted.length; i++) {
                var parents = splitted.slice(0, i + 1);
                children.push({
                    name: splitted[i],
                    url: '#/organizations/' + parents.join('.')
                });
            }
        }
        return children;

    }
})();
