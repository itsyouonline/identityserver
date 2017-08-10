(function () {
    'use strict';
    angular
        .module('itsyouonlineApp')
        .controller('SeeListController', ['$stateParams', 'UserService', SeeListController]);

    function SeeListController($stateParams, UserService) {
        var vm = this;
        var organization = $stateParams.organization;
        vm.documents = [];
        vm.loaded = false;
        vm.userIdentifier = null;

        init();

        function init() {
            getUserIdentifier();
            getDocuments();
        }

        function getUserIdentifier() {
            UserService.getUserIdentifier().then(function (userIdentifier) {
                vm.userIdentifier = userIdentifier;
            });
        }

        function getDocuments() {
            vm.loaded = false;
            UserService.getSeeObjects(organization).then(function (documents) {
                vm.documents = documents;
                vm.loaded = true;
            });
        }
    }

})();
