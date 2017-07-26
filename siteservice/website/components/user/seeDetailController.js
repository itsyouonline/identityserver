(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .controller("SeeDetailController", SeeDetailController);

    SeeDetailController.$inject =  ['$scope', '$rootScope', '$routeParams', '$location', '$window', '$q', '$translate',
        'UserService'];

    function SeeDetailController($scope, $rootScope, $routeParams, $location, $window, $q, $translate,
                                 UserService) {
        var vm = this,
            globalid = $routeParams.globalid,
            uniqueid = $routeParams.uniqueid;
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
