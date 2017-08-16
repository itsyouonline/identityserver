describe('Reset Password Controller test', function () {

    beforeEach(module('loginApp'));

    beforeEach(inject(function ($injector, $controller) {
        resetPasswordController = $controller('resetPasswordController', {

        });
    }));

    it('Reset Password Controller should be defined', function () {
        expect(resetPasswordController).toBeDefined();
    });
});
