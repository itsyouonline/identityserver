(function () {
    'use strict';

    angular
        .module('itsyouonlineApp')
        .controller('SeeDetailController', ['$stateParams', 'UserService', SeeDetailController]);

    function SeeDetailController($stateParams, UserService) {
        var vm = this,
            uniqueid = $stateParams.uniqueid,
            organization = $stateParams.globalid;

        vm.userIdentifier = undefined;
        vm.loading = true;
        vm.isShowingFullHistory = false;
        vm.toggleFullHistory = toggleFullHistory;

        init();

        function init() {
            getSee();
            UserService.getUserIdentifier().then(function (userIdentifier) {
                vm.userIdentifier = userIdentifier;
            })
        }

        function toggleFullHistory() {
            vm.isShowingFullHistory = !vm.isShowingFullHistory;
            getSee();
        }

        function getSee() {
            UserService
                .getSeeObject(uniqueid, organization, vm.isShowingFullHistory)
                .then(
                    function (data) {
                        data.versions.sort(function (a, b) {
                            return b.version - a.version;
                        });
                        vm.seeObject = data;
                        vm.loading = false;
                    }
                );
        }
    }

})();
