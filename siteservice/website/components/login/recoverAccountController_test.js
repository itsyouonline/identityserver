describe('Reset Password Controller test', function () {

    beforeEach(module('loginApp'));

    beforeEach(function () {
        module(function($provide) {
            $provide.value('$window', {
                location: {href: ''}
            });
        });
    });

    beforeEach(inject(function ($injector, $controller) {
        resetPasswordController = $controller('resetPasswordController', {

        });
    }));

    it('Reset Password Controller should be defined', function () {
        expect(resetPasswordController).toBeDefined();
    });
});
