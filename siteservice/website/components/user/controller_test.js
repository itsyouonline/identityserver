describe('User Home Controller test', function() {

    const TEST_USER = 'test';

    beforeEach(module('itsyouonlineApp'));

    var userService;

    beforeEach(inject(function ($injector, $rootScope, _UserService_, $controller) {

        //get reference to the service we want to track the method of
        userService = _UserService_;

        // This is set in the app bootstrap so lets just accept that this works for now
        userService.setUsername(TEST_USER);
        UserHomeController = $controller('UserHomeController', {
            'UserService': userService
        });

        spyOn(UserHomeController, 'loadUser');

    }));

    it('User Home Controller should be defined', function () {
       expect(UserHomeController).toBeDefined();
    });

    it('User Home Controller username property should be set', function () {
        expect(UserHomeController.username).toBe(TEST_USER);
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
