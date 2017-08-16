describe('Forgot Password Controller', function () {

    beforeEach(module('loginApp'));

    beforeEach(function () {
        module(function($provide) {
            $provide.value('$window', {
                location: {href: ''}
            });
        });
    });

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
