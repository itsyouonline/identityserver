describe('Registration Controller test', function() {

    beforeEach(module('itsyouonline.registration'));

    var scope;

    beforeEach(inject(function ( $window, $cookies, $mdUtil, $rootScope, configService, registrationService, $controller) {

        scope = $rootScope.$new();

        registrationController = $controller('registrationController', {
            $scope: scope,
            $window: $window,
            $cookies: $cookies,
            $mdUtil: $mdUtil,
            $rootScope: $rootScope,
            configService: configService,
            registrationService: registrationService
        });
    }));

    it('Registration Controller should be defined', function () {
        expect(registrationController).toBeDefined();
    });

    it('String startsWith method should be definded', function () {
        expect(String.prototype.startsWith).toBeDefined();
    });

    it('String includes method should be defined', function () {
      expect(String.prototype.includes).toBeDefined();
    });

    it('Array find method should be defined', function () {
      expect(Array.prototype.find).toBeDefined();
    });
    it('Array includes method should be defined', function () {
      expect(Array.prototype.includes).toBeDefined();
    });
});
