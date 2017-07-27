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
            globalid = $stateParams.globalid,
            uniqueid = $stateParams.uniqueid;
        vm.username = $rootScope.user;
        vm.loading = true;

        activate();

        function activate() {
            fetch();
        }

        function fetch(){
            UserService
                .getSeeObject(vm.username, globalid, uniqueid)
                .then(
                    function(data) {
                        vm.seeObject = data;
                        vm.loading = false;
                    }
                );
        }
    }

})();
