from flask import Blueprint, jsonify, request, Response, json
import jwt
import requests
import string

deliveries_api = Blueprint('deliveries_api', __name__)


@deliveries_api.route('/deliveries', methods=['GET'])
def deliveries_get():
    pubkey = "-----BEGIN PUBLIC KEY-----\nMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2\n7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6\n6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv\n-----END PUBLIC KEY-----"

    header = request.headers.get("Authorization")
    if not header or not header.startswith("bearer "):
        return 'Unauthorized', 401
    webtoken = header.split()
    decoded_jwt = jwt.decode(webtoken[1], pubkey, algorithms=["ES384"], audience="dronedelivery")
    print(decoded_jwt)
    if decoded_jwt["iss"] != "itsyouonline":
        return 'Unauthorized', 401
    else:
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
