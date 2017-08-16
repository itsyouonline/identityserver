describe('Resend Sms Controller test', function() {

    beforeEach(module('loginApp'));

    var scope;

    beforeEach(inject(function ($injector, $rootScope, $controller) {
        scope = $rootScope.$new();
        resendSmsController = $controller('resendSmsController', {
            $scope: scope
        });
    }));

    it('Resend Sms Controller should be defined', function () {
       expect(resendSmsController).toBeDefined();
    });
});
