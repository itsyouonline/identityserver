describe('Organization Detail Controller test', function () {

    beforeEach(module('itsyouonlineApp'));

    var scope;

    beforeEach(inject(function ($injector, $rootScope, $controller) {
        scope = $rootScope.$new();
        OrganizationDetailController = $controller('OrganizationDetailController', {
            $scope: scope
        });
    }));

    it('Organization Detail Controller should be defined', function () {
        expect(OrganizationDetailController).toBeDefined();
    })
});
