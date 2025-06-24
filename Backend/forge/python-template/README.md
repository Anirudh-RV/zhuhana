# python template

- Build:

  > docker build -t anirudhrv1234/template-django-app:latest .

- Push:

  > docker push anirudhrv1234/django-app:latest

- Run:
  > docker run --name tester --env-file env.env -e PYTHONUNBUFFERED=1 -p 8050:8000 anirudhrv1234/template-django-app:latest
