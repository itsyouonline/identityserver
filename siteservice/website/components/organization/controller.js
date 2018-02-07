(function() {
    'use strict';
    angular
        .module('itsyouonlineApp')
        .controller('OrganizationDetailController', OrganizationDetailController)
        .controller('InvitationDialogController', InvitationDialogController);

    InvitationDialogController.$inject = ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'UserDialogService'];
    OrganizationDetailController.$inject = ['$state', '$stateParams', '$window', '$translate', 'OrganizationService', '$mdDialog', '$mdMedia',
        'UserDialogService', 'UserService', 'ScopeService'];

    function OrganizationDetailController($state, $stateParams, $window, $translate, OrganizationService, $mdDialog, $mdMedia,
                                          UserDialogService, UserService, ScopeService) {
        var vm = this,
            globalid = $stateParams.globalid;
        vm.username = UserService.getUsername();
        vm.invitations = [];
        vm.apikeylabels = [];
        vm.organization = {};
        vm.organizationRoot = {};
        vm.childOrganizationNames = [];
        vm.logo = '';
        vm.includemap = {};
        vm.loading = {
            invitations: true,
            users: true,
            userIdentifiers: true
        };
        var pages = {
            'organization.people': 0,
            'organization.structure': 1,
            'organization.see': 2,
            'organization.settings': 3
        };
        vm.selectedTab = pages[$state.current.name];
        vm.userIdentifiers = undefined;

        vm.initSettings = initSettings;
        vm.showInvitationDialog = showInvitationDialog;
        vm.showAddOrganizationDialog = showAddOrganizationDialog;
        vm.showAPIKeyCreationDialog = showAPIKeyCreationDialog;
        vm.showAPIKeyDialog = showAPIKeyDialog;
        vm.showDNSDialog = showDNSDialog;
        vm.getOrganizationDisplayname = getOrganizationDisplayname;
        vm.fetchInvitations = fetchInvitations;
        vm.fetchAPIKeyLabels = fetchAPIKeyLabels;
        vm.showCreateOrganizationDialog = UserDialogService.createOrganization;
        vm.showDeleteOrganizationDialog = showDeleteOrganizationDialog;
        vm.editMember = editMember;
        vm.editOrgMember = editOrgMember;
        vm.canEditRole = canEditRole;
        vm.showLeaveOrganization = showLeaveOrganization;
        vm.showLogoDialog = showLogoDialog;
        vm.getSharedScopeString = getSharedScopeString;
        vm.showRequiredScopeDialog = showRequiredScopeDialog;
        vm.showMissingScopesDialog = showMissingScopesDialog;
        vm.showDescriptionDialog = showDescriptionDialog;
        vm.getScopeTranslation = getScopeTranslation;
        vm.removeInvitation = removeInvitation;
        vm.includeChanged = includeChanged;
        vm.trackinvite = trackinvite;

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
                        if (!vm.organization.orgowners) {
                          vm.organization.orgowners = [];
                        }
                        if (!vm.organization.orgmembers) {
                          vm.organization.orgmembers = [];
                        }
                        vm.childOrganizationNames = getChildOrganizations(vm.organization.globalid);
                        getUsers();
                        fillIncludeMap();
                    }
                );

            OrganizationService.getOrganizationTree(globalid)
                .then(function (data) {
                    vm.organizationRoot.children = [];
                    vm.organizationRoot.children.push(data);
                    var pixelWidth = 200 + getBranchWidth(vm.organizationRoot.children[0]);
                    vm.treeGraphStyle = {
                        'width': pixelWidth + 'px'
                    };
                });

            OrganizationService.getLogo(globalid).then(
                function(data) {
                    vm.logo = data.logo;
                }
            );
            UserService.GetAllUserIdentifiers().then(function (identifiers) {
                vm.userIdentifiers = identifiers;
                vm.loading.userIdentifiers = false;
            });
        }

        function fillIncludeMap() {
            if (vm.organization.includesuborgsof) {
                for (var i = 0; i < vm.organization.includesuborgsof.length; i++) {
                    vm.includemap[vm.organization.includesuborgsof[i]] = true;
                }
            }
        }

        function getBranchWidth(branch) {
            var splitted = branch.globalid.split('.');
            var length = splitted[splitted.length - 1].length * 6;
            var spacing = 0;
            if (branch.children.length > 1) {
                spacing = (branch.children.length - 1) * 80;
            }
            if (branch.children.length === 0) {
                return length;
            }
            else {
                var childWidth = spacing;
                for (var i = 0; i < branch.children.length; i++) {
                    childWidth += getBranchWidth(branch.children[i]);
                }
                return childWidth > length ? childWidth : length;
            }
        }

        function fetchInvitations() {
            if (!vm.hasEditPermission || vm.invitations.length) {
                return;
            }
            vm.loading.invitations = true;
            OrganizationService.getInvitations(globalid).then(function (data) {
                vm.loading.invitations = false;
                vm.invitations = data;
            });
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
                    }
                );
        }

        function initSettings() {
            fetchAPIKeyLabels();
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
                    var invite = {
                        created: invitation.created,
                        user: invitation.user || invitation.emailaddress || invitation.phonenumber,
                        role: invitation.role
                    };
                    vm.invitations.push(invite);
                });
        }

        function showAddOrganizationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: AddOrganizationDialogController,
                templateUrl: 'components/organization/views/addOrganizationMemberDialog.html',
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
                function(data) {
                    if (data.status) {
                        // this is an invite
                        var invite = {
                            created: data.created,
                            user: data.user,
                            role: data.role
                        };
                        vm.invitations.push(invite);
                    }
                    else if (data.role === 'members') {
                        if (!vm.organization.orgmembers) {
                            vm.organization.orgmembers = [];
                        }
                        vm.organization.orgmembers.push(data.organization);
                    } else {
                        if (!vm.organization.orgowners) {
                            vm.organization.orgowners = [];
                        }
                        vm.organization.orgowners.push(data.organization);
                    }
                });
        }

        function showAPIKeyCreationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', '$translate', 'organization', 'OrganizationService', 'label', APIKeyDialogController],
                templateUrl: 'components/organization/views/apikeydialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window,
                        label: ''
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
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', '$translate', 'organization', 'OrganizationService', 'label', APIKeyDialogController],
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
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'dnsName', DNSDialogController],
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

        function showLogoDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$document', '$mdDialog', 'organization', 'OrganizationService', logoDialogController],
                controllerAs: 'vm',
                templateUrl: 'components/organization/views/logoDialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals: {
                    OrganizationService: OrganizationService,
                    organization: vm.organization.globalid,
                    $window: $window
                }
            }).then(function () {
                OrganizationService.getLogo(vm.organization.globalid).then(function (data) {
                    vm.logo = data.logo;
                });
            });
        }

        function showDescriptionDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', descriptionDialogController],
                templateUrl: 'components/organization/views/descriptionDialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals: {
                    OrganizationService: OrganizationService,
                    organization: vm.organization.globalid,
                    $window: $window
                }
            });
        }

        function showDeleteOrganizationDialog(event) {
            $translate(['organization.controller.confirmdelete', 'organization.controller.deleteorg', 'organization.controller.deleteorganization',
                'organization.controller.haschildren', 'organization.controller.yes', 'organization.controller.no', 'delete_org_with_children'], {organization: globalid}).then(function(translations){
                    var text = translations['organization.controller.confirmdelete'];
                    if (hasActualChildren(vm.organization.globalid, vm.organizationRoot.children[0])) {
                       text = translations['delete_org_with_children'];
                    }
                    var confirm = $mdDialog.confirm()
                        .title(translations['organization.controller.deleteorg'])
                        .textContent(text)
                        .ariaLabel(translations['organization.controller.deleteorganization'])
                        .targetEvent(event)
                        .ok(translations['organization.controller.yes'])
                        .cancel(translations['organization.controller.no']);
                    $mdDialog.show(confirm).then(function () {
                        OrganizationService.deleteOrganization(globalid).then(function() {
                            // Check if there is a parent organization. If there is, redirect there, else go to the users profile page
                            var orgTree = $window.location.hash;
                            var url = '#/';
                            if (orgTree.indexOf('.') > -1) {
                                url = orgTree.slice(0, orgTree.lastIndexOf('.'));
                            }
                            $window.location.hash = url;
                        });
                    });
            });
        }

        function hasActualChildren(globalid, branch) {
            if (globalid === branch.globalid) {
                return branch.children.length > 0;
            }
            return hasActualChildren(globalid, branch.children[0]);
        }

        function canEditRole(member) {
            return vm.hasEditPermission && !vm.userIdentifiers.includes(member.username);
        }

        function editMember(event, user) {
            var username = user.username;
            var changeRoleDialog = {
                controller: ['$mdDialog', '$translate', 'OrganizationService', 'UserDialogService', 'organization', 'user', 'initialRole', EditOrganizationMemberController],
                controllerAs: 'ctrl',
                templateUrl: 'components/organization/views/changeRoleDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    organization: vm.organization,
                    user: username,
                    initialRole: user.role
                }
            };

            $mdDialog
                .show(changeRoleDialog)
                .then(function (data) {
                    if (data.action === 'edit') {
                        vm.organization = data.data;
                        var u = vm.users.filter(function (user) {
                            return user.username === username;
                        })[0];

                        u.role = data.newRole;
                    } else if (data.action === 'remove') {
                        var people = vm.organization[data.data.role];
                        people.splice(people.indexOf(data.data), 1);
                        vm.users = vm.users.filter(function (user) {
                            return user.username !== username;
                        });
                    }
                    setUsers();
                });
        }

        function editOrgMember(event, org) {
            var role = 'orgmembers';
            angular.forEach(['orgmembers', 'orgowners'], function (r) {
                if (vm.organization[r].indexOf(org) !== -1) {
                    role = r;
                }
            });
            var changeOrgRoleDialog = {
                controller: ['$mdDialog', '$translate', 'OrganizationService', 'UserDialogService', 'organization', 'org', 'initialRole', EditOrganizationMemberOrgController],
                controllerAs: 'ctrl',
                templateUrl: 'components/organization/views/changeOrganizationRoleDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    organization: vm.organization,
                    org: org,
                    initialRole: role
                }
            };

            $mdDialog
                .show(changeOrgRoleDialog)
                .then(function (data) {
                    if (data.action === 'edit') {
                        vm.organization = data.data;
                    } else if (data.action === 'remove') {
                        var collection;
                        if (data.data.role === 'orgmembers') {
                            collection = vm.organization.orgmembers;
                        } else {
                            collection = vm.organization.orgowners;
                        }
                        collection.splice(collection.indexOf(data.data.org), 1);
                    }
                });
        }

        function showLeaveOrganization(event) {
            $translate(['organization.controller.confirmleave', 'organization.controller.leaveorg', 'organization.controller.leaveorganization', 'organization.controller.yes',
                'organization.controller.no', 'organization.controller.notfound', 'last_owner', 'error'], {organization: globalid}).then(function(translations){
                    var text = translations['organization.controller.confirmleave'];
                    var confirm = $mdDialog.confirm()
                        .title(translations['organization.controller.leaveorg'])
                        .textContent(text)
                        .ariaLabel(translations['organization.controller.leaveorganization'])
                        .targetEvent(event)
                        .ok(translations['organization.controller.yes'])
                        .cancel(translations['organization.controller.no']);
                    $mdDialog
                        .show(confirm)
                        .then(function () {
                            UserService
                                .leaveOrganization(vm.username, globalid)
                                .then(function () {
                                    $window.location.hash = '#/';
                                }, function (response) {
                                    if (response.status === 404) {
                                        UserDialogService.showSimpleDialog(translations['organization.controller.notfound'], translations['error'], null, event);
                                    }
                                    if (response.status === 409) {
                                        UserDialogService.showSimpleDialog(translations['last_owner'], translations['error'], null, event);
                                    }
                                });
                        });
            });
        }

        function getUsers() {
            vm.loading.users = true;
            OrganizationService.getUsers(globalid).then(function (response) {
                vm.loading.users = false;
                vm.users = response.users;
                vm.hasEditPermission = response.haseditpermissions;
                setUsers();
                if (vm.hasEditPermission) {
                    fetchInvitations();
                }
            });
        }

        function setUsers() {
            vm.members = vm.users.filter(function (user) {
                return user.role === 'members';
            });
            vm.owners = vm.users.filter(function (user) {
                return user.role === 'owners';
            });
        }

        function getSharedScopeString(scopes) {
            return scopes.map(function (scope) {
                return scope.replace('organization:', '') + 's';
            }).join(', ');
        }

        function showRequiredScopeDialog(event, requiredScope) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            var dialog = {
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'requiredScope', 'ScopeService', RequiredScopesDialogController],
                controllerAs: 'vm',
                templateUrl: 'components/organization/views/requiredScopesDialog.html',
                targetEvent: event,
                fullscreen: useFullScreen,
                locals: {
                    OrganizationService: OrganizationService,
                    ScopeService: ScopeService,
                    organization: vm.organization.globalid,
                    $window: $window,
                    requiredScope: requiredScope
                }
            };
            $mdDialog.show(dialog).then(function (data) {
                if (data.action === 'create') {
                    vm.organization.requiredscopes.push(data.requiredScope);
                }
                if (data.action === 'update') {
                    var scope = vm.organization.requiredscopes[vm.organization.requiredscopes.indexOf(requiredScope)];
                    angular.copy(data.requiredScope, scope);
                }
                if (data.action === 'delete') {
                    vm.organization.requiredscopes.splice(vm.organization.requiredscopes.indexOf(requiredScope), 1);
                }
            });
        }

        function showMissingScopesDialog(event, user) {
            var title = 'Missing information';
            var msg = 'This user hasn\'t shared some required information:<br />';
            $translate(['organization.controller.missinginfo', 'organization.controller.missinguserinfo']).then(function(translations){
                title = translations['organization.controller.missinginfo'];
                msg = translations['organization.controller.missinguserinfo'];
                msg += user.missingscopes.join('<br />');
                var closeText = 'Close';
                UserDialogService.showSimpleDialog(msg, title, closeText, event);
            });
        }

        function getScopeTranslation(scope) {
            return ScopeService.parseScope(scope);
        }

        function removeInvitation(invite) {
            var searchString = invite.user || invite.phonenumber || invite.emailaddress;
            OrganizationService.removeInvitation(globalid, searchString).then(removeFromView, function (response) {
                if (response.status === 404) {
                    removeFromView();
                } else {
                    $window.location.href = 'error' + response.status;
                }
            });
            function removeFromView() {
                vm.invitations.splice(vm.invitations.indexOf(invite), 1);
            }
        }

        function includeChanged(org) {
            if (vm.includemap[org]) {
                OrganizationService.addInclude(vm.organization.globalid, org);
                return;
            }
            OrganizationService.removeInclude(vm.organization.globalid, org);

        }

        function trackinvite(invite) {
            return invite.user + invite.role;
        }
    }

    function getOrganizationDisplayname(globalid) {
        if (globalid) {
            var split = globalid.split('.');
            return split[split.length - 1];
        }
    }

    function InvitationDialogController($scope, $mdDialog, organization, OrganizationService, UserDialogService) {

        $scope.role = 'members';

        $scope.cancel = cancel;
        $scope.invite = invite;
        $scope.validationerrors = {};


        function cancel(){
            $mdDialog.cancel();
        }

        function invite(searchString, role){
            $scope.validationerrors = {};
            OrganizationService.invite(organization, searchString, role).then(
                function(data){
                    $mdDialog.hide(data);
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else if (reason.status == 404){
                        $scope.validationerrors.nosuchuser = true;
                    } else if (reason.status === 422) {
                        cancel();
                        var msg = 'Organization ' + organization + ' has reached the maximum amount of invitations.';
                        UserDialogService.showSimpleDialog(msg, 'Error');
                    }
                }
            );

        }
    }

    function AddOrganizationDialogController($scope, $mdDialog, organization, OrganizationService) {

        $scope.role = 'members';

        $scope.organization = organization;

        $scope.cancel = cancel;
        $scope.addOrganization = addOrganization;
        $scope.validationerrors = {};


        function cancel(){
            $mdDialog.cancel();
        }

        function addOrganization(searchString, role){
            $scope.validationerrors = {};
            OrganizationService.addOrganization(organization, searchString, role).then(
                function(data){
                    if (!data) {
                        $mdDialog.hide({organization: searchString,role: role});
                    }
                    else {
                        $mdDialog.hide(data);
                    }
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else if (reason.status == 404){
                        $scope.validationerrors.nosuchorganization = true;
                    }
                    else if (reason.status === 422) {
                        cancel();
                        var msg = 'Organization ' + organization + ' has reached the maximum amount of invitations.';
                        UserDialogService.showSimpleDialog(msg, 'Error');
                    }
                }
            );

        }
    }

    function APIKeyDialogController($scope, $mdDialog, $translate, organization, OrganizationService, label) {
        //If there is a key, it is already saved, if not, this means that a new secret is being created.

        $scope.apikey = {secret: ''};

        if (label) {
            $translate(['organization.controller.loadingkey']).then(function(translations){
                $scope.secret = translations['organization.controller.loadingkey'];
                OrganizationService.getAPIKey(organization, label).then(
                    function(data){
                        $scope.apikey = data;
                    }
                );
            });
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
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }

        function update(oldLabel, newLabel){
            $scope.validationerrors = {};
            OrganizationService.updateAPIKey(organization, oldLabel, newLabel, $scope.apikey).then(
                function () {
                    $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel});
                },
                function(reason){
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }


        function deleteAPIKey(label){
            $scope.validationerrors = {};
            OrganizationService.deleteAPIKey(organization, label).then(
                function () {
                    $mdDialog.hide({originalLabel: label, newLabel: ''});
                }
            );
        }

    }

    function DNSDialogController($scope, $mdDialog, organization, OrganizationService, dnsName) {
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
                    $mdDialog.hide({originalDns: '', newDns: data.name});
                },
                function (reason) {
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
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
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }


        function remove(dnsName) {
            $scope.validationerrors = {};
            OrganizationService.deleteDNS(organization, dnsName)
                .then(function () {
                    $mdDialog.hide({originalDns: dnsName, newDns: ''});
                });
        }
    }

    function getChildOrganizations(organization) {
        var children = [];
        if (organization) {
            for (var i = 0, splitted = organization.split('.'); i < splitted.length; i++) {
                var parents = splitted.slice(0, i + 1);
                children.push({
                    name: splitted[i],
                    url: '#/organization/' + parents.join('.') + '/people'
                });
            }
        }
        return children;

    }

    function EditOrganizationMemberController($mdDialog, $translate, OrganizationService, UserDialogService, organization, user, initialRole) {
        var ctrl = this;
        ctrl.role = initialRole;
        ctrl.user = user;
        ctrl.organization = organization;
        ctrl.cancel = cancel;
        ctrl.submit = submit;
        ctrl.remove = remove;

        function cancel() {
            $mdDialog.cancel();
        }

        function submit() {
            OrganizationService
                .updateMembership(organization.globalid, ctrl.user, ctrl.role)
                .then(function (data) {
                    $mdDialog.hide({action: 'edit', data: data, newRole: ctrl.role});
                }, function () {
                    $translate(['organization.controller.cantchangerole']).then(function(translations){
                        UserDialogService.showSimpleDialog(translations['organization.controller.cantchangerole'], 'Error', 'ok', event);
                    });
                });
        }

        function remove() {
            OrganizationService
                .removeMember(organization.globalid, user, initialRole)
                .then(function () {
                    $mdDialog.hide({action: 'remove', data: {role: initialRole, user: user}});
                }, function (response) {
                    $mdDialog.cancel(response);
                });
        }
    }

    function EditOrganizationMemberOrgController($mdDialog, $translate, OrganizationService, UserDialogService, organization, org, initialRole) {
        var ctrl = this;
        ctrl.role = initialRole;
        ctrl.org = org;
        ctrl.organization = organization;
        ctrl.cancel = cancel;
        ctrl.submit = submit;
        ctrl.remove = remove;

        function cancel() {
            $mdDialog.cancel();
        }

        function submit() {
            OrganizationService
                .updateOrgMembership(organization.globalid, ctrl.org, ctrl.role)
                .then(function (data) {
                    $mdDialog.hide({action: 'edit', data: data});
                }, function () {
                    $translate(['organization.controller.cantchangerole']).then(function(translations){
                        UserDialogService.showSimpleDialog(translations['organization.controller.cantchangerole'], 'Error', 'ok', event);
                    });
                });
        }

        function remove() {
            OrganizationService
                .removeOrgMember(organization.globalid, org, initialRole)
                .then(function () {
                    $mdDialog.hide({action: 'remove', data: {role: initialRole, org: org}});
                }, function (response) {
                    $mdDialog.cancel(response);
                });
        }
    }

    function logoDialogController($scope, $document, $mdDialog, organization, OrganizationService) {
        var doc = $document[0];
        $scope.organization = organization;
        $scope.setFile = setFile;
        $scope.cancel = cancel;
        $scope.validationerrors = {};
        $scope.update = update;
        $scope.remove = remove;
        $scope.logoChanged = false;
        OrganizationService.getLogo(organization).then(
            function(data) {
                $scope.logo = data.logo;
                $scope.logoChanged = !$scope.logo;
                makeFileDrop();
            }
        );

        function makeFileDrop() {
            var target = doc.getElementById('logo-upload-preview');
            target.addEventListener('dragover', function (e) {
                e.preventDefault();
            }, true);
            target.addEventListener('drop', function (src) {
	              src.preventDefault();
                //only allow image files, ignore others
                if (!src.dataTransfer.files[0] || !src.dataTransfer.files[0].type.match(/image.*/)) {
                    return;
                }
                var reader = new FileReader();
	              reader.onload = function(e){
		                setFile(e.target.result);
	              };
	              reader.readAsDataURL(src.dataTransfer.files[0]);
            }, true);
        }

        $scope.uploadFile = function(event){
            var files = event.target.files;
            var url = URL.createObjectURL(files[0]);
            setFile(url);
        };

        function setFile(url) {
            var img = new Image();
            var WIDTH = 500;
            var HEIGHT = 240;
            img.src = url;
            $scope.logoChanged = true;

            img.onload = function() {
                var hRatio = 1;
                var wRatio = 1;
                // take the smallest ratio
                if (img.width > WIDTH)
                    wRatio = WIDTH / img.width;
                if (img.height > HEIGHT)
                    hRatio = HEIGHT / img.height;
                var ratio = wRatio > hRatio ? hRatio : wRatio;
                var canvas = doc.createElement('canvas');
                var copyContext = canvas.getContext('2d');

                canvas.width = img.width * ratio;
                canvas.height = img.height * ratio;
                copyContext.drawImage(img, 0, 0, canvas.width, canvas.height);
                $scope.logo = canvas.toDataURL();
                // forces the update button after a file drop, might fix safari issues?
                $scope.$digest();
            };

        }

        function cancel() {
            $mdDialog.cancel();
        }

        function update(logo) {
            if (!$scope.form.$valid) {
                return;
            }
            $scope.validationerrors = {};
            OrganizationService.setLogo(organization, logo).then(
                function () {
                    $mdDialog.hide({logo: logo});
                },
                function (reason) {
                    if (reason.status === 413) {
                        $scope.validationerrors.filesize = true;
                    }
                }
            );
        }


        function remove() {
            $scope.validationerrors = {};
            OrganizationService.deleteLogo(organization)
                .then(function () {
                    $mdDialog.hide({logo: ''});
                });
        }
    }

    function descriptionDialogController($scope, $mdDialog, organization, OrganizationService) {
        $scope.organization = organization;
        $scope.selectedLangKey = 'en';
        $scope.descriptionExists = false;
        $scope.cancel = cancel;
        $scope.remove = remove;
        $scope.update = update;
        $scope.save = save;
        $scope.loadDescription = loadDescription;
        loadDescription();

        function cancel() {
            $mdDialog.cancel();
        }

        function remove() {
            OrganizationService.deleteDescription(organization, $scope.selectedLangKey)
                .then(function() {
                    $mdDialog.hide();
                });
        }

        function update() {
            OrganizationService.updateDescription(organization, $scope.selectedLangKey, $scope.description)
                .then(function() {
                    $mdDialog.hide();
                });
        }

        function save() {
            OrganizationService.saveDescription(organization, $scope.selectedLangKey, $scope.description)
                .then(function() {
                    $mdDialog.hide();
                });
        }

        function loadDescription() {
            OrganizationService.getDescription(organization, $scope.selectedLangKey)
                .then(function(data) {
                    $scope.description = data.text;
                    $scope.descriptionExists = !!data.text;
                });
        }

    }


    function RequiredScopesDialogController($scope, $mdDialog, organization, OrganizationService, requiredScope, ScopeService) {
        var vm = this;
        vm.possibleScopes = ScopeService.getScopes();
        var allAccessScopes = [{
            scope: 'organization:owner',
            translation: 'owners',
            checked: false
        }, {
            scope: 'organization:member',
            translation: 'members',
            checked: false
        }];
        vm.organization = organization;
        var updatedRequiredScope = requiredScope ? JSON.parse(JSON.stringify(requiredScope)) : {
            scope: '',
            accessscopes: []
        };
        vm.originalScope = requiredScope ? requiredScope.scope : null;
        vm.accessScopes = allAccessScopes.map(function (scope) {
            scope.checked = updatedRequiredScope.accessscopes.indexOf(scope.scope) !== -1;
            return scope;
        });
        vm.validationerrors = {};
        vm.cancel = cancel;
        vm.create = create;
        vm.update = update;
        vm.remove = remove;
        vm.allowedScopesChanged = allowedScopesChanged;
        vm.scopeTypedChanged = scopeTypedChanged;
        vm.scopeChanged = scopeChanged;

        init();

        function init() {
            vm.editingScope = ScopeService.parseScope(updatedRequiredScope.scope);
        }

        function scopeTypedChanged() {
            vm.editingScope.scope.vars = ScopeService.parseScope(vm.editingScope.scope.scope).scope.vars;
            scopeChanged();
        }

        function scopeChanged() {
            updatedRequiredScope.scope = ScopeService.makeScope(vm.editingScope);
        }

        function allowedScopesChanged() {
            updatedRequiredScope.accessscopes = vm.accessScopes.filter(function (scope) {
                return scope.checked;
            }).map(function (scope) {
                return scope.scope;
            });
        }

        function cancel() {
            $mdDialog.cancel();
        }

        function create() {
            if (!vm.form.$valid) {
                return;
            }
            vm.showPleaseSelectAccessScope = updatedRequiredScope.accessscopes.length === 0;
            if (vm.showPleaseSelectAccessScope) {
                return;
            }
            vm.validationerrors = {};
            OrganizationService.createRequiredScope(organization, updatedRequiredScope).then(
                function () {
                    $mdDialog.hide({action: 'create', requiredScope: updatedRequiredScope});
                },
                function (reason) {
                    if (reason.status === 409) {
                        vm.validationerrors.duplicate = true;
                    }
                }
            );
        }

        function update() {
            if (!vm.form.$valid) {
                return;
            }
            vm.showPleaseSelectAccessScope = updatedRequiredScope.accessscopes.length === 0;
            if (vm.showPleaseSelectAccessScope) {
                return;
            }
            vm.validationerrors = {};
            OrganizationService.updateRequiredScope(organization, vm.originalScope, updatedRequiredScope).then(
                function () {
                    $mdDialog.hide({action: 'update', requiredScope: updatedRequiredScope});
                },
                function (reason) {
                    if (reason.status === 409) {
                        vm.validationerrors.duplicate = true;
                    }
                }
            );
        }


        function remove() {
            $scope.validationerrors = {};
            OrganizationService.deleteRequiredScope(organization, vm.originalScope)
                .then(function () {
                    $mdDialog.hide({action: 'delete'});
                });
        }
    }

})();
