(function() {
    'use strict';


    angular.module("itsyouonlineApp").service("CompanyService",CompanyService);


    CompanyService.$inject = ['$http','$q'];

    function CompanyService($http, $q) {
        var apiURL =  'api/companies';

        var service = {
            create: create
        };
        return service;

        function create(name, taxnr){
            var url = apiURL;
            return $http.post(url, {globalid:name,taxnr:taxnr}).then(
                function(response) {
                    return response.data;
                },
                function(reason){
                    return $q.reject(reason);
                }
            );
        }

    }


})();
