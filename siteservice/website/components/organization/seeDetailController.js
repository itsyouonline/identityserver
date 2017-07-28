(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .controller("SeeDetailController", SeeDetailController);

    SeeDetailController.$inject =  ['$scope', '$rootScope', '$stateParams', '$location', '$window', '$q', '$translate',
        'UserService'];

    function SeeDetailController($scope, $rootScope, $stateParams, $location, $window, $q, $translate,
                                 UserService) {
        var vm = this,
          uniqueid = $stateParams.uniqueid;
        vm.globalid = $stateParams.globalid;
        vm.username = $rootScope.user;
        vm.loading = true;
        vm.isShowingFullHistory = false;
        vm.toggleFullHistory = toggleFullHistory

        activate();

        function activate() {
            fetch();
        }

        function toggleFullHistory(event) {
          vm.isShowingFullHistory = !vm.isShowingFullHistory;
          fetch();
        }

        function fetch(){
            UserService
                .getSeeObject(vm.username, vm.globalid, uniqueid, vm.isShowingFullHistory)
                .then(
                    function(data) {
                        vm.seeObject = data;
                        vm.seeObject.versions.sort(function(a, b) {
                            return b.version - a.version;
                        })
                        vm.loading = false;
                    }
                );
        }


    }

})();
