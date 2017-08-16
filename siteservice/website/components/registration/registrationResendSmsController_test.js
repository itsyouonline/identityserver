describe('Resend Sms Controller test', function () {

    beforeEach(module('itsyouonline.registration'));

    var scope;

    beforeEach(inject(function ($injector, $rootScope, $controller) {
        scope = $rootScope.$new();
        resendSmsController = $controller('resendSmsController', {
            $scope: scope
        });
    }));

    it('resendSmsController should be defined', function () {
        expect(resendSmsController).toBeDefined();
    });
});
