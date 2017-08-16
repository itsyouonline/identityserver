describe('Two Factor Authentication Controller test', function () {

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
        twoFactorAuthenticationController = $controller('twoFactorAuthenticationController', {
            $scope: scope
        });
    }));

    it('Two Factor Authentication Controller should be defined', function () {
        expect(twoFactorAuthenticationController).toBeDefined();
    });
});
