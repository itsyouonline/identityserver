from flask import Blueprint, jsonify, request



drones_api = Blueprint('drones_api', __name__)


@drones_api.route('/drones', methods=['GET'])
def drones_get():
    '''
    Get a list of drones
    It is handler for GET /drones
    '''
    
    return jsonify()


@drones_api.route('/drones', methods=['POST'])
def drones_post():
    '''
    Add a new drone to the fleet
    It is handler for POST /drones
    '''
    
    return jsonify()


@drones_api.route('/drones/<droneId>', methods=['GET'])
def drones_byDroneId_get(droneId):
    '''
    Get information on a specific drone
    It is handler for GET /drones/<droneId>
    '''
    
    return jsonify()


@drones_api.route('/drones/<droneId>', methods=['PATCH'])
def drones_byDroneId_patch(droneId):
    '''
    Update the information on a specific drone
    It is handler for PATCH /drones/<droneId>
    '''
    
    return jsonify()


@drones_api.route('/drones/<droneId>', methods=['DELETE'])
def drones_byDroneId_delete(droneId):
    '''
    Remove a drone from the fleet
    It is handler for DELETE /drones/<droneId>
    '''
    
    return jsonify()


@drones_api.route('/drones/<droneId>/deliveries', methods=['GET'])
def drones_byDroneId_deliveries_get(droneId):
    '''
    The deliveries scheduled for the current drone
    It is handler for GET /drones/<droneId>/deliveries
    '''
    
    return jsonify()
