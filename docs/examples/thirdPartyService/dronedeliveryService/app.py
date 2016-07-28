from flask import Flask, send_from_directory, send_file

import wtforms_json
from deliveries import deliveries_api
from drones import drones_api

app = Flask(__name__)

app.config["WTF_CSRF_ENABLED"] = False
wtforms_json.init()

app.register_blueprint(deliveries_api)
app.register_blueprint(drones_api)


@app.route('/apidocs/<path:path>')
def send_js(path):
    return send_from_directory('apidocs', path)


@app.route('/', methods=['GET'])
def home():
    return send_file('index.html')


if __name__ == "__main__":
    app.run(debug=True)
