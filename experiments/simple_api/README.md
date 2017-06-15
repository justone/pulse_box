Steps:

1. Make a Python3 virtual env
2. `pip install -r requirements.txt`
3. `FLASK_APP=app.py SERIAL=/dev/ttyUSB0 flask run`
4. `curl -v localhost:5000/blink/8`
