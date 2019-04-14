(function() {
    'use strict';


    angular.module("itsyouonlineApp").service("OrganizationService", OrganizationService);


    OrganizationService.$inject = ['$http','$q'];

    function OrganizationService($http, $q) {
        var apiURL =  'api/organizations';
        var GET = $http.get;
        var POST = $http.post;
        var PUT = $http.put;
        var DELETE = $http.delete;

        return {
            create: create,
            get: get,
            invite: invite,
            removeInvitation: removeInvitation,
            addOrganization: addOrganization,
            getUserOrganizations: getUserOrganizations,
            getInvitations: getInvitations,
            createAPIKey: createAPIKey,
            deleteAPIKey: deleteAPIKey,
            updateAPIKey: updateAPIKey,
            getAPIKeyLabels: getAPIKeyLabels,
            getAPIKey: getAPIKey,
            getOrganizationTree: getOrganizationTree,
            getUsers: getUsers,
            createDNS: createDNS,
            updateDNS: updateDNS,
            deleteDNS: deleteDNS,
            deleteOrganization: deleteOrganization,
            updateMembership: updateMembership,
            updateOrgMembership: updateOrgMembership,
            removeMember: removeMember,
            removeOrgMember: removeOrgMember,
            getLogo: getLogo,
            setLogo: setLogo,
            deleteLogo: deleteLogo,
            getValidityDuration: getValidityDuration,
            SetValidityDuration: SetValidityDuration,
            createRequiredScope: createRequiredScope,
            updateRequiredScope: updateRequiredScope,
            deleteRequiredScope: deleteRequiredScope,
            getDescription: getDescription,
            deleteDescription: deleteDescription,
            updateDescription: updateDescription,
            saveDescription: saveDescription,
            addInclude: addInclude,
            removeInclude: removeInclude,
            checkRequiresTwoFA: checkRequiresTwoFA,
            forceTwoFA: forceTwoFA,
            unforceTwoFA: unforceTwoFA
        };

        function genericHttpCall(httpFunction, url, data) {
            return httpFunction(url, data)
                .then(
                    function (response) {
                        return response.data;
                    },
                    function (reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function create(name, dns, owner, parentOrganization) {
            var url = apiURL;
            if (parentOrganization){
                url += '/' + encodeURIComponent(parentOrganization) + '/suborganizations';
                name = parentOrganization + '.' + name;
            }
            var data = {
                globalid: name,
                dns: dns,
                owners: [owner]
            };
            return genericHttpCall(POST, url, data);
        }

        function get(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid);

            return genericHttpCall(GET, url);
        }

        function invite(globalid, searchString, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/' + encodeURIComponent(role);
            var data = {searchstring: searchString};
            return genericHttpCall(POST, url, data);
        }

        function removeInvitation(globalid, searchString) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/invitations/' + encodeURIComponent(searchString);
            return genericHttpCall(DELETE, url);
        }

        function addOrganization(globalid, searchString, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/' + 'org' + encodeURIComponent(role);

            var data;
            if (role === "members") {
                data = {orgmember: searchString};
            } else {
                data = {orgowner: searchString};
            }
            return genericHttpCall(POST, url, data);
        }

        function getUserOrganizations(username) {
            var url = '/api/users/' + encodeURIComponent(username) + '/organizations';
            return genericHttpCall(GET, url);
        }

        function getInvitations(globalid){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/invitations';
            return genericHttpCall(GET, url);
        }

        function getAPIKeyLabels(globalid){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys';
            return genericHttpCall(GET, url);
        }

        function createAPIKey(globalid, apiKey){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys';
            return genericHttpCall(POST, url, apiKey);
        }

        function updateAPIKey(globalid, oldLabel, newLabel, apikey){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys/' + encodeURIComponent(oldLabel);
            apikey.label = newLabel;
            return genericHttpCall(PUT, url, apikey);
        }

        function deleteAPIKey(globalid, label){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys/' + encodeURIComponent(label);
            return genericHttpCall(DELETE, url);
        }

        function getAPIKey(globalid, label){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys/' + encodeURIComponent(label);
            return genericHttpCall(GET, url);
        }

        function getOrganizationTree(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/tree';
            return genericHttpCall(GET, url);
        }

        function getUsers(globalId) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/users';
            return genericHttpCall(GET, url);
        }

        function createDNS(globalid, dnsName) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/dns';
            var data = {
                name: dnsName
            };
            return genericHttpCall(POST, url, data);
        }

        function updateDNS(globalid, oldDnsName, newDnsName) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/dns/' + encodeURIComponent(oldDnsName);
            var data = {
                name: newDnsName
            };
            return genericHttpCall(PUT, url, data);
        }

        function deleteDNS(globalid, dnsName) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/dns/' + encodeURIComponent(dnsName);
            return genericHttpCall(DELETE, url);
        }

        function deleteOrganization(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid);
            return genericHttpCall(DELETE, url);
        }

        function updateMembership(globalid, username, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/members';
            var data = {
                username: username,
                role: role
            };
            return genericHttpCall(PUT, url, data);
        }

        function updateOrgMembership(globalid, org, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/orgmembers';
            var data = {
                org: org,
                role: role
            };
            return genericHttpCall(PUT, url, data);
        }

        function removeMember(globalid, username, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/' + role + '/' + username;
            return genericHttpCall(DELETE, url);
        }

        function removeOrgMember(globalid, org, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/' + role + '/' + encodeURIComponent(org);
            return genericHttpCall(DELETE, url);
        }

        function getLogo(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/logo';
            return genericHttpCall(GET, url);
        }

        function setLogo(globalid, logo) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/logo';
            var data = {
                globalid: globalid,
                logo: logo
            };
            return genericHttpCall(PUT, url, data);
        }

        function deleteLogo(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/logo';
            return genericHttpCall(DELETE, url);
        }

        function getValidityDuration(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/2fa/validity';
            return genericHttpCall(GET, url);
        }

        function SetValidityDuration(globalid, secondsduration) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/2fa/validity';
            var data = {
                secondsvalidity: secondsduration
            };
            return genericHttpCall(PUT, url, data);
        }

        function createRequiredScope(globalId, requiredScope) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/requiredscopes';
            return genericHttpCall(POST, url, requiredScope);
        }

        function updateRequiredScope(globalId, oldRequiredScope, newRequiredScope) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/requiredscopes/' + encodeURIComponent(oldRequiredScope);
            return genericHttpCall(PUT, url, newRequiredScope);
        }

        function deleteRequiredScope(globalId, requiredScope) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/requiredscopes/' + encodeURIComponent(requiredScope);
            return genericHttpCall(DELETE, url);
        }

        function getDescription(globalId, langKey) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/description/' + encodeURIComponent(langKey);
            return genericHttpCall(GET, url);
        }

        function deleteDescription(globalId, langKey) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/description/' + encodeURIComponent(langKey);
            return genericHttpCall(DELETE, url);
        }

        function updateDescription(globalId, langKey, description) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/description';
            var data = {
                langkey: langKey,
                text: description
            };
            return genericHttpCall(PUT, url, data);
        }

        function saveDescription(globalId, langKey, description) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/description';
            var data = {
                langkey: langKey,
                text: description
            };
            return genericHttpCall(POST, url, data);
        }

        function addInclude(globalId, include) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/orgmembers/includesuborgs';
            var data = {
                globalid: include
            };
            return genericHttpCall(POST, url, data);
        }

        function removeInclude(globalId, include) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/orgmembers/includesuborgs/' + encodeURIComponent(include);
            return genericHttpCall(DELETE, url);
        }

        function checkRequiresTwoFA(globalId) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/2fa/isrequired';
            return genericHttpCall(GET, url);
        }

        function forceTwoFA(globalId) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/2fa/setrequired';
            return genericHttpCall(PUT, url);
        }

        function unforceTwoFA(globalId) {
            var url = apiURL + '/' + encodeURIComponent(globalId) + '/2fa/setoptional';
            return genericHttpCall(PUT, url);
        }
    }
})();
