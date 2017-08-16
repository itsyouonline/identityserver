describe('Resend Sms Controller test', function() {

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
        resendSmsController = $controller('resendSmsController', {
            $scope: scope
        });
    }));

    it('Resend Sms Controller should be defined', function () {
       expect(resendSmsController).toBeDefined();
    });
});
