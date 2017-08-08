(function () {
    'use strict';
    angular
        .module('itsyouonlineApp')
        .controller('SeeListController', ['UserService', SeeListController]);

    function SeeListController(UserService) {
        var vm = this;
        vm.documents = [];
        vm.loaded = {
            documents: false
        };
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
            vm.loaded.documents = true;
            UserService.getSeeObjects().then(function (documents) {
                vm.documents = documents;
                vm.loaded.documents = true;
            });
        }
    }

})();
