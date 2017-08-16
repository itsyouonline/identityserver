describe('Login Controller test', function() {

    beforeEach(module('loginApp'));

    var scope;

    beforeEach(inject(function ($http, $window, $rootScope, $interval, LoginService, $controller) {

        scope = $rootScope.$new();

        loginController = $controller('loginController', {
          $http: $http,
          $window: $window,
          $scope: scope,
          $rootScope: $rootScope,
          $interval: $interval,
          LoginService: LoginService
        });
    }));

    it('loginController should be defined', function() {
        expect(loginController).toBeDefined();
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
