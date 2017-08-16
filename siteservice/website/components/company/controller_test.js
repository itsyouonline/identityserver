describe('Company controller test ', function() {

  beforeEach(module('itsyouonlineApp'));

  beforeEach(function () {
        module(function($provide) {
            $provide.value('$window', {
                location: {href: ''}
            });
        });
    });

  var CompanyController, location;

  beforeEach(inject(function ($location, CompanyService, $controller) {
    location = $location;
    CompanyController = $controller('CompanyController', {
      $location: location,
      CompanyService: CompanyService
    });
  }));

  it('CompanyController defined', function() {
    expect(CompanyController).toBeDefined();
  });

});
