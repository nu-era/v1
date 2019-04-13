FROM python:3-alpine
WORKDIR /app
COPY . /app
RUN pip3 install -r requirements.txt
EXPOSE 5000-5110
ENTRYPOINT ["python3", "app.py"]