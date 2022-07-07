FROM python:3.10

WORKDIR /app
COPY . /app/

RUN pip install poetry
RUN poetry export > requirements.txt

RUN pip install -r requirements.txt
EXPOSE 8000
ENTRYPOINT [ "./start.py" ]
CMD [ "./start.py", "--generate", "-u='--host 0.0.0.0'" ]
