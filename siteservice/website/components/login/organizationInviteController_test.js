describe('Organization Invite Controller test', function () {

    beforeEach(module('loginApp'));

    beforeEach(function () {
        module(function($provide) {
            $provide.value('$window', {
                location: {href: ''}
            });
        });
    });

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
