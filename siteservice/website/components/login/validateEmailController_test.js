describe('Validate Email Controller', function () {

    beforeEach(module('loginApp'));

    var scope;

    beforeEach(inject(function ($injector, $rootScope, $controller) {
        scope = $rootScope.$new();
        validateEmailController = $controller('validateEmailController', {
            $scope: scope
        });
    }));

    it('Validate Email Controller should be defined', function () {
        expect(validateEmailController).toBeDefined();
    });
});
