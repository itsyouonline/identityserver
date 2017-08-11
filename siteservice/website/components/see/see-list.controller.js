(function () {
    'use strict';
    angular
        .module('itsyouonlineApp')
        .controller('SeeListController', ['$stateParams', 'UserService', SeeListController]);

    function SeeListController($stateParams, UserService) {
        var vm = this;
        vm.organization = $stateParams.globalid;
        vm.documents = [];
        vm.loaded = false;
        vm.userIdentifier = undefined;
        vm.noDocsTranslation = vm.organization ? 'no_see_documents_for_organization' : 'no_see_documents';

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
            UserService.getSeeObjects(vm.organization).then(function (documents) {
                vm.documents = documents;
                vm.loaded = true;
            });
        }
    }

})();
