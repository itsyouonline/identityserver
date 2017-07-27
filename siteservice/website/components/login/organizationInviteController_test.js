describe('Organization Invite Controller test', function () {

    beforeEach(module('loginApp'));

    beforeEach(inject(function ($http, $window, $routeParams, $mdDialog, $controller) {

        organizationInviteController = $controller('organizationInviteController', {
            $http: $http,
            $window: $window,
            $routeParams: $routeParams,
            $mdDialog: $mdDialog
        });
    }));

    it('organizationInviteController should be defined', function() {
        expect(organizationInviteController).toBeDefined();
    });
})
