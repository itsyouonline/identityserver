(function () {
    'use strict';
    angular
        .module('itsyouonlineApp')
        .service('ScopeService', ScopeService);

    function ScopeService() {
        var scopes = [{
            scope: 'user:name',
            translation: 'full_name'
        }, {
            scope: 'user:facebook',
            translation: 'facebook_account'
        }, {
            scope: 'user:github',
            translation: 'github_account'
        }, {
            scope: 'user:memberof:{organization}',
            translation: 'member_of_organization',
            vars: {
                organization: true
            }
        }, {
            scope: 'user:address:{label}',
            translation: 'address',
            vars: {
                label: true
            }
        }, {
            scope: 'user:bankaccount:{label}',
            translation: 'bank_account',
            vars: {
                label: true
            }
        }, {
            scope: 'user:digitalwalletaddress:{label}:{currency}',
            translation: 'digital_wallet_address',
            vars: {
                label: true,
                currency: true
            }
        }, {
            scope: 'user:email:{label}',
            translation: 'email_address',
            vars: {
                label: true
            }
        }, {
            scope: 'user:validated:email:{label}',
            translation: 'validated_email_address',
            vars: {
                label: true
            }
        }, {
            scope: 'user:phone:{label}:{write}',
            translation: 'phone_number',
            vars: {
                label: true,
                write: true
            }
          }, {
              scope: 'user:validated:phone:{label}',
              translation: 'validated_phone_number',
              vars: {
                  label: true
              }
        }, {
            scope: 'user:publickey:{label}',
            translation: 'public_key',
            vars: {
                label: true
            }
        }, {
            scope: 'user:avatar:{label}',
            translation: 'avatar',
            vars: {
                label: true
            }
        }, {
            scope: 'user:keystore',
            translation: 'keystore'
        }, {
            scope: 'user:see',
            translation: 'see'
        }];

        function getScopes() {
            return scopes;
        }

        function parseScope(scope) {
            var label, organization, currency, base, splitPermission, splitScope, write;
            splitPermission = scope.split(':');
            base = splitPermission[0] + ':' + splitPermission[1];
            var scopeObject = scopes.find(function (s) {
                return s.scope.startsWith(base);
            });
            if (!scopeObject) {
                return;
            }
            splitScope = scopeObject.scope.split(':');
            angular.forEach(splitScope, function (scopePart, i) {
                switch (scopePart) {
                    case '{label}':
                        label = splitPermission[i];
                        break;
                    case '{organization}':
                        organization = splitPermission[i];
                        break;
                    case '{write}':
                        write = splitPermission[i] === 'write';
                        break;
                }
            });
            return {
                scope: scopeObject,
                label: label || '',
                organization: organization || '',
                currency: currency || '',
                write: write || false
            };
        }

        function makeScope(scopeObject) {
            var scopeTemplate = scopeObject.scope.scope,
                label = scopeObject.label,
                currency = scopeObject.currency,
                organization = scopeObject.organization,
                write = scopeObject.write;
            if (scopeTemplate.includes('{label}')) {
                scopeTemplate = scopeTemplate.replace(label ? '{label}' : ':{label}', label || '');
            }
            if (scopeTemplate.includes('{currency}')) {
                scopeTemplate = scopeTemplate.replace(currency ? '{currency}' : ':{currency}', currency || '');
            }
            if (scopeTemplate.includes('{organization}')) {
                scopeTemplate = scopeTemplate.replace(organization ? '{organization}' : ':{organization}', organization || '');
            }
            if (scopeTemplate.includes('{write}')) {
                scopeTemplate = scopeTemplate.replace(':{write}', write === true ? ':write' : '');
            }
            return scopeTemplate;
        }

        return {
            getScopes: getScopes,
            makeScope: makeScope,
            parseScope: parseScope
        };
    }

})();
