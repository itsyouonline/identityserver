(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .service("UserService", UserService)
        .service("NotificationService", NotificationService);

    UserService.$inject = ['$http','$q'];
    NotificationService.$inject = ['$http','$q'];

    function NotificationService($http, $q) {
        var apiURL = 'api/users';

        var service = {
            get: get,
            accept: accept,
            reject: reject
        };

        return service;

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

        function accept(invitation) {
            var url = apiURL + '/' + encodeURIComponent(invitation.user) + '/organizations/' + encodeURIComponent(invitation.organization) + '/roles/' + encodeURIComponent(invitation.role) ;

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

        function reject(invitation) {
            var url = apiURL + '/' + encodeURIComponent(invitation.user) + '/organizations/' + encodeURIComponent(invitation.organization) + '/roles/' + encodeURIComponent(invitation.role) ;

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

        return {
            get: get,
            registerNewEmailAddress: registerNewEmailAddress,
            updateEmailAddress: updateEmailAddress,
            deleteEmailAddress: deleteEmailAddress,
            registerNewPhonenumber: registerNewPhonenumber,
            updatePhonenumber: updatePhonenumber,
            deletePhonenumber: deletePhonenumber,
            registerNewAddress: registerNewAddress,
            updateAddress: updateAddress,
            deleteAddress: deleteAddress,
            getAuthorizations: getAuthorizations,
            saveAuthorization: saveAuthorization,
            deleteAuthorization: deleteAuthorization,
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
            deleteAPIKey: deleteAPIKey
        };

        function genericHttpCall(httpFunction, url, data) {
            if (data){
                return httpFunction(url, data)
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
                return httpFunction(url)
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

        function get(username) {
            var url = apiURL + '/' + encodeURIComponent(username);
            return genericHttpCall($http.get, url);
        }

        function update(username, user) {
            var url = apiURL + '/' + encodeURIComponent(username);
            return genericHttpCall($http.put, url, user);
        }

        function registerNewEmailAddress(username, label, emailaddress) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses';
            return genericHttpCall($http.post, url, {label: label, emailaddress: emailaddress});
        }

        function updateEmailAddress(username, oldlabel, newlabel, emailaddress) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses/' + encodeURIComponent(oldlabel) ;
            return genericHttpCall($http.put, url, {label: newlabel, emailaddress: emailaddress});
        }

        function deleteEmailAddress(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses/' + encodeURIComponent(label) ;
            return genericHttpCall($http.delete, url);
        }

        function registerNewPhonenumber(username, label, phonenumber) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers';
            return genericHttpCall($http.post, url, {label: label, phonenumber: phonenumber});
        }

        function updatePhonenumber(username, oldlabel, newlabel, phonenumber) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(oldlabel) ;
            return genericHttpCall($http.put, url, {label: newlabel, phonenumber: phonenumber});
        }

        function deletePhonenumber(username, label, force) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(label) + '?force=' + !!force;
            return genericHttpCall($http.delete, url);
        }

        function registerNewAddress(username, label, address) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/addresses';
            return genericHttpCall($http.post, url, {label: label, address: address});
        }

        function updateAddress(username, oldlabel, newlabel, address) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/addresses/' + encodeURIComponent(oldlabel) ;
            return genericHttpCall($http.put, url, {label: newlabel, address: address});
        }

        function deleteAddress(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/addresses/' + encodeURIComponent(label) ;
            return genericHttpCall($http.delete, url);
        }

        function registerNewBankAccount(username, label, bankAccount) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/banks';
            return genericHttpCall($http.post, url, {label: label, bank: bankAccount});
        }

        function updateBankAccount(username, oldLabel, label, bankAccount) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/banks/' + encodeURIComponent(oldLabel);
            return genericHttpCall($http.put, url, {label: label, bank: bankAccount});
        }

        function deleteBankAccount(username, label) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/banks/' + encodeURIComponent(label);
            return genericHttpCall($http.delete, url);
        }

        function getAuthorizations(username) {
            var url = apiURL + '/' + encodeURIComponent(username) + '/authorizations/';
            return genericHttpCall($http.get, url);
        }

        function saveAuthorization(authorization) {
            var url = apiURL + '/' +  encodeURIComponent(authorization.username) + '/authorizations/' + encodeURIComponent(authorization.grantedTo);
            return genericHttpCall($http.put, url, authorization);
        }

        function deleteAuthorization(authorization) {
            var url = apiURL + '/' + encodeURIComponent(authorization.username) + '/authorizations/' + encodeURIComponent(authorization.grantedTo);
            return genericHttpCall($http.delete, url);
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
            var url = apiURL + '/' + encodeURIComponent(username) + '/phonenumbers/' + encodeURIComponent(label) + '/validate';
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
            var url = apiURL + '/' + encodeURIComponent(username) + '/emailaddresses/' + encodeURIComponent(label) + '/validate';
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
    }
})();
