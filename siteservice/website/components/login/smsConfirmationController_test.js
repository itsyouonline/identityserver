describe('Sms Confirmation Controller test', function () {

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
        smsConfirmationController = $controller('smsConfirmationController', {
            $scope: scope
        });
    }));

    it('Sms Confirmation Controller should be defined', function () {
        expect(smsConfirmationController).toBeDefined();
    });
});
