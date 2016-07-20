from flask import Blueprint, jsonify, request, Response, json

deliveries_api = Blueprint('deliveries_api', __name__)


@deliveries_api.route('/deliveries', methods=['GET'])
def deliveries_get():
    '''
    Get a list of deliveries
    It is handler for GET /deliveries
    '''
    data = {
        "id": "4",
        "at": "Tue, 08 Jul 2014 13:00:00 GMT",
        "toAddressId": "gi6w4fgi",
        "orderItemId": "6782798",
        "status": "completed",
        "droneId": "f"
    }
    return Response(json.dumps(data), mimetype='application/json')


@deliveries_api.route('/deliveries', methods=['POST'])
def deliveries_post():
    '''
    Create/request a new delivery
    It is handler for POST /deliveries
    '''

    return jsonify()


@deliveries_api.route('/deliveries/<deliveryId>', methods=['GET'])
def deliveries_byDeliveryId_get(deliveryId):
    '''
    Get information on a specific delivery
    It is handler for GET /deliveries/<deliveryId>
    '''

    return jsonify()


@deliveries_api.route('/deliveries/<deliveryId>', methods=['PATCH'])
def deliveries_byDeliveryId_patch(deliveryId):
    '''
    Update the information on a specific delivery
    It is handler for PATCH /deliveries/<deliveryId>
    '''

    return jsonify()


@deliveries_api.route('/deliveries/<deliveryId>', methods=['DELETE'])
def deliveries_byDeliveryId_delete(deliveryId):
    '''
    Cancel a specific delivery
    It is handler for DELETE /deliveries/<deliveryId>
    '''

    return jsonify()
