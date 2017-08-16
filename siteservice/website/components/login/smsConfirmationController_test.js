describe('Sms Confirmation Controller test', function () {

    beforeEach(module('loginApp'));

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
