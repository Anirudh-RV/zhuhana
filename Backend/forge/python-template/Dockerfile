# Dockerfile
FROM python:3.13-slim

# Set work directory
WORKDIR /code

# Install dependencies
COPY requirements.txt .
RUN pip install --upgrade pip && pip install -r requirements.txt

# Copy project
COPY . .

# Run the app
CMD ["gunicorn", "useralgorithm.wsgi:application", "--bind", "0.0.0.0:8000"]
