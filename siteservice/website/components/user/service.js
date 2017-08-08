(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .service("UserService", UserService)
        .service("NotificationService", NotificationService);

    UserService.$inject = ['$http', '$q', '$rootScope'];
    NotificationService.$inject = ['$http','$q'];

    function NotificationService($http, $q) {
        var apiURL = 'api/users';

        return {
            get: get,
            accept: accept,
            reject: reject,
            acceptorganizationinvite: acceptorganizationinvite,
            rejectorganizationinvite: rejectorganizationinvite
        };

        function get(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/notifications';

            return $http
                .get(url)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function accept(invitation, user) {
            var url = apiURL + '/' + encodeURIComponent(user) + '/organizations/' + encodeURIComponent(invitation.organization) + '/roles/' + encodeURIComponent(invitation.role) ;
            return $http
                .post(url, invitation)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function reject(invitation, user) {
            var url = apiURL + '/' + encodeURIComponent(user) + '/organizations/' + encodeURIComponent(invitation.organization) + '/roles/' + encodeURIComponent(invitation.role) ;

            return $http
                .delete(url)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function acceptorganizationinvite(invitation) {
            var url = 'api/organizations/' + encodeURIComponent(invitation.user) + '/organizations/' + encodeURIComponent(invitation.organization) + '/roles/' + encodeURIComponent(invitation.role);

            return $http
                .post(url, invitation)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
            );
        }

        function rejectorganizationinvite(invitation) {
            var url = 'api/organizations/' + encodeURIComponent(invitation.user) + '/organizations/' + encodeURIComponent(invitation.organization) + '/roles/' + encodeURIComponent(invitation.role);

            return $http
                .delete(url)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
            );
        }
    }

    function UserService($http, $q) {
        var apiURL = 'api/users';
        var GET = $http.get,
            POST = $http.post,
            PUT = $http.put,
            DELETE = $http.delete;
        var username = null;

        return {
            get: get,
            getUserIdentifier: getUserIdentifier,
            registerNewEmailAddress: registerNewEmailAddress,
            updateEmailAddress: updateEmailAddress,
            deleteEmailAddress: deleteEmailAddress,
            registerNewPhonenumber: registerNewPhonenumber,
            updatePhonenumber: updatePhonenumber,
            deletePhonenumber: deletePhonenumber,
            registerNewAddress: registerNewAddress,
            updateAddress: updateAddress,
            deleteAddress: deleteAddress,
            getAuthorization: getAuthorization,
            getAuthorizations: getAuthorizations,
            saveAuthorization: saveAuthorization,
            deleteAuthorization: deleteAuthorization,
            getSeeObjects: getSeeObjects,
            getSeeObject: getSeeObject,
            registerNewBankAccount: registerNewBankAccount,
            updateBankAccount: updateBankAccount,
            deleteBankAccount: deleteBankAccount,
            deleteFacebookAccount: deleteFacebookAccount,
            deleteGithubAccount: deleteGithubAccount,
            updatePassword: updatePassword,
            updateName: updateName,
            getVerifiedPhones: getVerifiedPhones,
            sendPhoneVerificationCode: sendPhoneVerificationCode,
            verifyPhone: verifyPhone,
            getVerifiedEmailAddresses: getVerifiedEmailAddresses,
            sendEmailAddressVerification: sendEmailAddressVerification,
            getAPIKeys: getAPIKeys,
            createAPIKey: createAPIKey,
            updateAPIKey: updateAPIKey,
            deleteAPIKey: deleteAPIKey,
            createPublicKey: createPublicKey,
            updatePublicKey: updatePublicKey,
            deletePublicKey: deletePublicKey,
            getTwoFAMethods: getTwoFAMethods,
            getAuthenticatorSecret: getAuthenticatorSecret,
            setAuthenticator: setAuthenticator,
            removeAuthenticator: removeAuthenticator,
            createDigitalWalletAddress: createDigitalWalletAddress,
            updateDigitalWalletAddress: updateDigitalWalletAddress,
            deleteDigitalWalletAddress: deleteDigitalWalletAddress,
            leaveOrganization: leaveOrganization,
            getAvatars: getAvatars,
            createAvatarFromLink: createAvatarFromLink,
            createAvatarFromFile: createAvatarFromFile,
            updateAvatarLink: updateAvatarLink,
            updateAvatarFile: updateAvatarFile,
            deleteAvatar: deleteAvatar,
            getUsername: getUsername,
            setUsername: setUsername
        };

        function genericHttpCall(httpFunction, url, data, config) {
            if (data){
                return httpFunction(url, data, config)
                    .then(
                        function(response) {
                            return response.data;
                        },
                        function(reason) {
                            return $q.reject(reason);
                        }
                    );
            }
            else {
                return httpFunction(url, config)
                    .then(
                        function(response) {
                            return response.data;
                        },
                        function(reason) {
                            return $q.reject(reason);
                        }
                    );
            }
        }

        function get(refresh) {
            var url = apiURL + '/' + encodeURIComponent(username);
            // cache: true prevents from doing the same request twice and just returns the value of the previous request to this url
            return genericHttpCall(GET, url, null, {cache: !refresh});
        }

        function getUserIdentifier() {
            var deferred = $q.defer();
            this.get().then(function (user) {
                if (user.firstname || user.lastname) {
                    deferred.resolve(user.firstname + ' ' + user.lastname);
                    return;
                }
                var email = user.emailaddresses.filter(function (email) {
                    return email.verified;
                })[0];
                if (email) {
                    return deferred.resolve(email.emailaddress);
                }
                var phone = user.phonenumbers.filter(function (phone) {
                    return phone.verified;
                })[0];
                if (phone) {
                    return deferred.resolve(phone.phonenumber);
                }
            });
            return deferred.promise;

        }

        function getUsername() {
            return username;
        }

        function setUsername(name) {
            username = name;
        }

        function registerNewEmailAddress(username, emailaddress) {
            var lang = localStorage.getItem('langKey');
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses?lang=' + lang;
            return genericHttpCall($http.post, url, emailaddress);
        }

        function updateEmailAddress(username, oldlabel, emailaddress) {
            var lang = localStorage.getItem('langKey');
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses/' + encodeURIComponent(oldlabel) + '?lang=' + lang;
            return genericHttpCall($http.put, url, emailaddress);
        }

        function deleteEmailAddress(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses/' + encodeURIComponent(label) ;
            return genericHttpCall($http.delete, url);
        }

        function registerNewPhonenumber(username, phonenumber) {
            phonenumber.phonenumber = phonenumber.phonenumber.replace(/ /g, '');
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers';
            return genericHttpCall($http.post, url, phonenumber);
        }

        function updatePhonenumber(username, oldlabel, phonenumber) {
            phonenumber.phonenumber = phonenumber.phonenumber.replace(/ /g, '');
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(oldlabel) ;
            return genericHttpCall($http.put, url, phonenumber);
        }

        function deletePhonenumber(username, label, force) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(label) + '?force=' + !!force;
            return genericHttpCall($http.delete, url);
        }

        function registerNewAddress(username, address) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/addresses';
            return genericHttpCall($http.post, url, address);
        }

        function updateAddress(username, oldlabel, address) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/addresses/' + encodeURIComponent(oldlabel) ;
            return genericHttpCall($http.put, url, address);
        }

        function deleteAddress(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/addresses/' + encodeURIComponent(label) ;
            return genericHttpCall($http.delete, url);
        }

        function registerNewBankAccount(username, bankAccount) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/banks';
            return genericHttpCall($http.post, url, bankAccount);
        }

        function updateBankAccount(username, oldLabel, bankAccount) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/banks/' + encodeURIComponent(oldLabel);
            return genericHttpCall($http.put, url, bankAccount);
        }

        function deleteBankAccount(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/banks/' + encodeURIComponent(label);
            return genericHttpCall($http.delete, url);
        }

        function getAuthorizations(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/authorizations';
            return genericHttpCall($http.get, url);
        }

        function getAuthorization(username, organization) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/authorizations/' + encodeURIComponent(organization);
            return genericHttpCall($http.get, url);
        }

        function saveAuthorization(authorization) {
            var url = apiURL + '/' + encodeURIComponent(authorization.username) + '/authorizations/' + encodeURIComponent(authorization.grantedTo);
            return genericHttpCall($http.put, url, authorization);
        }

        function deleteAuthorization(authorization) {
            var url = apiURL + '/' + encodeURIComponent(authorization.username) + '/authorizations/' + encodeURIComponent(authorization.grantedTo);
            return genericHttpCall($http.delete, url);
        }

        function getSeeObjects() {
            var url = apiURL + '/' + encodeURIComponent(username) + '/see';
            return genericHttpCall($http.get, url, null, {cache: true});
        }

        function getSeeObject(uniqueId, all) {
            var queryString = '';
            if (all) {
                queryString = '?version=0';
            }
            var url = apiURL + '/' + encodeURIComponent(username) + '/see/' + encodeURIComponent(uniqueId) + queryString;
            return genericHttpCall($http.get, url, null, {cache: true});
        }

        function deleteFacebookAccount(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/facebook';
            return genericHttpCall($http.delete, url);
        }

        function deleteGithubAccount(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/github';
            return genericHttpCall($http.delete, url);
        }

        function updatePassword(username, currentPassword, newPassword) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/password';
            var data = {
                currentpassword: currentPassword,
                newpassword: newPassword
            };
            return genericHttpCall($http.put, url, data);
        }

        function updateName(username, firstname, lastname) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/name';
            var data = {
                firstname: firstname,
                lastname: lastname
            };
            return genericHttpCall($http.put, url, data);
        }

        function getVerifiedPhones(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers?validated=true';
            return genericHttpCall($http.get, url);
        }

        function sendPhoneVerificationCode(username, label) {
            var lang = localStorage.getItem('langKey');
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(label) + '/validate?lang=' + lang;
            return genericHttpCall($http.post, url);
        }

        function verifyPhone(username, label, validationKey, confirmationCode) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(label) + '/validate';
            var data = {
                smscode: confirmationCode,
                validationkey: validationKey
            };
            return genericHttpCall($http.put, url, data);
        }

        function getVerifiedEmailAddresses(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses?validated=true';
            return genericHttpCall($http.get, url);
        }

        function sendEmailAddressVerification(username, label){
            var lang = localStorage.getItem('langKey');
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses/' + encodeURIComponent(label) + '/validate?lang=' + lang;
            return genericHttpCall($http.post, url);
        }

        function getAPIKeys(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/apikeys';
            return genericHttpCall($http.get, url);
        }

        function createAPIKey(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/apikeys';
            var data = {
                label: label
            };
            return genericHttpCall($http.post, url, data);
        }

        function updateAPIKey(username, oldLabel, newLabel) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/apikeys/' + encodeURIComponent(oldLabel);
            var data = {
                label: newLabel
            };
            return genericHttpCall($http.put, url, data);
        }

        function deleteAPIKey(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/apikeys/' + encodeURIComponent(label);
            return genericHttpCall($http.delete, url);
        }

        function createPublicKey(username, label, publicKey) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/publickeys';
            var data = {
                label: label,
                publicKey: publicKey
            };
            return genericHttpCall(POST, url, data);
        }

        function updatePublicKey(username, oldLabel, newLabel, publicKey) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/publickeys/' + encodeURIComponent(oldLabel);
            var data = {
                label: newLabel,
                publicKey: publicKey
            };
            return genericHttpCall(PUT, url, data);
        }

        function deletePublicKey(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/publickeys/' + encodeURIComponent(label);
            return genericHttpCall(DELETE, url);
        }

        function getTwoFAMethods(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/twofamethods';
            return genericHttpCall($http.get, url);
        }

        function getAuthenticatorSecret(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/totp';
            return genericHttpCall($http.get, url);
        }

        function setAuthenticator(username, secret, code) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/totp';
            var data = {
                totpsecret: secret,
                totpcode: code
            };
            return genericHttpCall($http.post, url, data);
        }

        function removeAuthenticator(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/totp';
            return genericHttpCall($http.delete, url);
        }

        function createDigitalWalletAddress(username, walletAddress) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/digitalwallet';
            return genericHttpCall(POST, url, walletAddress);
        }

        function updateDigitalWalletAddress(username, oldLabel, walletAddress) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/digitalwallet/' + encodeURIComponent(oldLabel);
            return genericHttpCall(PUT, url, walletAddress);
        }

        function deleteDigitalWalletAddress(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/digitalwallet/' + encodeURIComponent(label);
            return genericHttpCall(DELETE, url);
        }

        function leaveOrganization(username, globalid) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/organizations/' + encodeURIComponent(globalid) + '/leave';
            return genericHttpCall(DELETE, url);
        }

        function getAvatars(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/avatars';
            return genericHttpCall(GET, url);
        }

        function createAvatarFromLink(username, label, link) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/avatar';
            var data = {
                label: label,
                source: link
            };
            return $http.post(url, data);
        }

        function createAvatarFromFile(username, label, file) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/avatar/img/' + encodeURIComponent(label);
            var formData = new FormData();
            formData.append('file', file);
            return $http.post(url, formData, {
              transform: angular.identity,
              headers: {'Content-Type': undefined}
            });
        }

        function updateAvatarLink(username, oldLabel, newLabel, link) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/avatar/' + encodeURIComponent(oldLabel);
            var data = {
                label: newLabel,
                source: link
            };
            return $http.put(url, data);
        }

        function updateAvatarFile(username, oldLabel, newLabel, file) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/avatar/' +
                encodeURIComponent(oldLabel) + '/to/' + encodeURIComponent(newLabel);
            var formData = new FormData();
            formData.append('file', file);
            return $http.put(url, formData, {
                transform: angular.identity,
                headers: {'Content-Type': undefined}
            });
        }

        function deleteAvatar(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/avatar/' + encodeURIComponent(label);
            return genericHttpCall(DELETE, url)
        }
    }
})();
