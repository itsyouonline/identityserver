// Karma configuration

module.exports = function(config) {
    var configuration = {

        // base path that will be used to resolve all patterns (eg. files, exclude)
        basePath: 'siteservice/website',


        // frameworks to use
        // available frameworks: https://npmjs.org/browse/keyword/karma-adapter
        frameworks: ['jasmine'],


        // list of files / patterns to load in the browser
        files: [

          // load in all required js files from the project.
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular.js',
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular-route.js',
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular-messages.min.js',
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular-animate.min.js',
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular-aria.min.js',
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular-cookies.min.js',
          'https://ajax.googleapis.com/ajax/libs/angularjs/1.5.9/angular-sanitize.min.js',
          'https://ajax.googleapis.com/ajax/libs/angular_material/1.1.1/angular-material.min.js',
          'https://ajax.googleapis.com/ajax/libs/angular-ui-router/0.4.2/angular-ui-router.min.js',
          '../../node_modules/angular-mocks/angular-mocks.js',
          'thirdpartyassets/angular-translate/angular-translate.min.js',
          'thirdpartyassets/angular-translate-loader-static-files/angular-translate-loader-static-files.min.js',
          'thirdpartyassets/angular-translate-handler-log/angular-translate-handler-log.min.js',
          'components/shared/directives/header.js',
          'components/shared/directives/footer.js',
          'components/shared/directives/validation.js',
          'thirdpartyassets/URI.js',
          'thirdpartyassets/qrcode-generator/js/qrcode.js',
          'thirdpartyassets/qrcode-generator/js/qrcode_UTF8.js',
          'thirdpartyassets/angular-qrcode/angular-qrcode.js',
          'components/shared/shared.js',
          'components/shared/configService.js',
          'components/shared/customOnChange.js',
          'components/patches.js',
          'components/common.js',
          'components/shared/country-info.js',
          'components/shared/directives/telinput.js',
          'thirdpartyassets/showdown/compressed/Showdown.min.js',
          'thirdpartyassets/angular-markdown-directive/markdown.js',

          // load main app and dependancies
          'components/app.js',
          'components/user/directives/authorizationDetailsDirective.js',
          'components/user/directives/treeDirective.js',
          'components/user/UserDialogService.js',
          'components/user/authorizeController.js',
          'components/user/controller.js',
          'components/user/service.js',
          'components/organization/controller.js',
          'components/organization/service.js',
          'components/organization/scope.service.js',
          'components/company/controller.js',
          'components/company/service.js',

          // load login app and dependancies
          'components/login/loginApp.js',
          'components/login/loginController.js',
          'components/login/loginService.js',
          'components/login/resendSMSController.js',
          'components/login/smsConfirmationController.js',
          'components/login/twoFactorAuthenticationController.js',
          'components/login/forgotPasswordController.js',
          'components/login/recoverAccountController.js',
          'components/login/organizationInviteController.js',
          'components/login/validateEmailController.js',

          // load register app and dependancies
          'components/registration/registrationApp.js',
          'components/registration/registrationService.js',
          'components/registration/registrationController.js',
          'components/registration/registrationSmsController.js',
          'components/registration/registrationResendSmsController.js',

          // now that all dependancies are loaded, load the test files
          '**/*_test.js'
        ],


        // list of files to exclude
        exclude: [
        ],


        // preprocess matching files before serving them to the browser
        // available preprocessors: https://npmjs.org/browse/keyword/karma-preprocessor
        preprocessors: {
        },


        // test results reporter to use
        // possible values: 'dots', 'progress'
        // available reporters: https://npmjs.org/browse/keyword/karma-reporter
        reporters: ['progress'],


        // web server port
        port: 9876,


        // enable / disable colors in the output (reporters and logs)
        colors: true,


        // level of logging
        // possible values: config.LOG_DISABLE || config.LOG_ERROR || config.LOG_WARN || config.LOG_INFO || config.LOG_DEBUG
        logLevel: config.LOG_DEBUG,


        // enable / disable watching file and executing tests whenever any file changes
        autoWatch: false,


        // start these browsers
        // available browser launchers: https://npmjs.org/browse/keyword/karma-launcher
        browsers: ['ChromeHeadless', 'Firefox'],


        // Continuous Integration mode
        // if true, Karma captures browsers, runs the tests and exits
        singleRun: true,

        // Concurrency level
        // how many browser should be started simultaneous
        concurrency: Infinity,

    };


    config.set(configuration);
}
