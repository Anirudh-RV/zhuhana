FROM python:3.13-slim

WORKDIR /app
COPY . /app

COPY requirements.txt .
RUN pip install --upgrade pip && pip install -r requirements.txt

CMD ["gunicorn", "script_runner.wsgi:application", "--bind", "0.0.0.0:8000"]
