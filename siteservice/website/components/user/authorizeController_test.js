describe('Authorize Conctroller test', function() {

    beforeEach(module('itsyouonlineApp'));

    var scope;

    beforeEach(inject(function ($injector, $rootScope, $controller) {
        scope = $rootScope.$new();
        AuthorizeController = $controller('AuthorizeController', {
            $scope: scope
        });
    }));

    it('AuthorizeController should be defined', function () {
        expect(AuthorizeController).toBeDefined();
    });
});
