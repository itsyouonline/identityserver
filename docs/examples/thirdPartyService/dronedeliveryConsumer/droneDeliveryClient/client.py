import requests
from .client_utils import build_query_string

BASE_URI = ""


class Client:
    def __init__(self):
        self.url = BASE_URI
        self.session = requests.Session()
        self.auth_header = ''
    
    def set_auth_header(val):
        ''' set authorization header value'''
        self.auth_header = val

    def deliveries_get(self, headers=None, query_params=None):
        """
        Get a list of deliveries
        It is method for GET /deliveries
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/deliveries"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri, headers=headers)


    def deliveries_post(self, data, headers=None, query_params=None):
        """
        Create/request a new delivery
        It is method for POST /deliveries
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/deliveries"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, data, headers=headers)


    def deliveries_byDeliveryId_get(self, deliveryId, headers=None, query_params=None):
        """
        Get information on a specific delivery
        It is method for GET /deliveries/{deliveryId}
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/deliveries/"+deliveryId
        uri = uri + build_query_string(query_params)
        return self.session.get(uri, headers=headers)


    def deliveries_byDeliveryId_patch(self, data, deliveryId, headers=None, query_params=None):
        """
        Update the information on a specific delivery
        It is method for PATCH /deliveries/{deliveryId}
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/deliveries/"+deliveryId
        uri = uri + build_query_string(query_params)
        return self.session.patch(uri, data, headers=headers)


    def deliveries_byDeliveryId_delete(self, deliveryId, headers=None, query_params=None):
        """
        Cancel a specific delivery
        It is method for DELETE /deliveries/{deliveryId}
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/deliveries/"+deliveryId
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri, headers=headers)


    def drones_get(self, headers=None, query_params=None):
        """
        Get a list of drones
        It is method for GET /drones
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/drones"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri, headers=headers)


    def drones_post(self, data, headers=None, query_params=None):
        """
        Add a new drone to the fleet
        It is method for POST /drones
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/drones"
        uri = uri + build_query_string(query_params)
        return self.session.post(uri, data, headers=headers)


    def drones_byDroneId_get(self, droneId, headers=None, query_params=None):
        """
        Get information on a specific drone
        It is method for GET /drones/{droneId}
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/drones/"+droneId
        uri = uri + build_query_string(query_params)
        return self.session.get(uri, headers=headers)


    def drones_byDroneId_patch(self, data, droneId, headers=None, query_params=None):
        """
        Update the information on a specific drone
        It is method for PATCH /drones/{droneId}
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/drones/"+droneId
        uri = uri + build_query_string(query_params)
        return self.session.patch(uri, data, headers=headers)


    def drones_byDroneId_delete(self, droneId, headers=None, query_params=None):
        """
        Remove a drone from the fleet
        It is method for DELETE /drones/{droneId}
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/drones/"+droneId
        uri = uri + build_query_string(query_params)
        return self.session.delete(uri, headers=headers)


    def drones_byDroneId_deliveries_get(self, droneId, headers=None, query_params=None):
        """
        The deliveries scheduled for the current drone
        It is method for GET /drones/{droneId}/deliveries
        """
        if self.auth_header:
            if not headers:
                headers = {'Authorization': self.auth_header}
            else:
                headers['Authorization'] = self.auth_header

        uri = self.url + "/drones/"+droneId+"/deliveries"
        uri = uri + build_query_string(query_params)
        return self.session.get(uri, headers=headers)
