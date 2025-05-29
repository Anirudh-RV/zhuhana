FROM python:3.13-slim

WORKDIR /app
COPY . /app

COPY requirements.txt .
RUN pip install --upgrade pip && pip install -r requirements.txt

# Run only the script, nothing else
CMD ["python", "script_runner.py"]
