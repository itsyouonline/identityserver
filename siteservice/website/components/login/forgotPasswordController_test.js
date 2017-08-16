describe('Forgot Password Controller', function () {

    beforeEach(module('loginApp'));

    var scope;

    beforeEach(inject(function ($injector, $rootScope, $controller) {
        scope = $rootScope.$new();
        forgotPasswordController = $controller('forgotPasswordController', {
            $scope: scope
        });
    }));

    it('Forgot Password Controller should be defined', function () {
        expect(forgotPasswordController).toBeDefined();
    });
});
