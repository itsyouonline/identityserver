describe('Organization Invite Controller test', function () {

    beforeEach(module('loginApp'));

    beforeEach(inject(function ($http, $window, $stateParams, $mdDialog, $controller) {

        organizationInviteController = $controller('organizationInviteController', {
            $http: $http,
            $window: $window,
            $stateParams: $stateParams,
            $mdDialog: $mdDialog
        });
    }));

    it('organizationInviteController should be defined', function() {
        expect(organizationInviteController).toBeDefined();
    });
})
