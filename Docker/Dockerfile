FROM mongo

RUN mkdir /home/ppov
COPY ./pku_ppov /home/ppov/
COPY ./begin.sh /home/ppov/

RUN mkdir /data/mongo
WORKDIR /home/ppov/
CMD ["sh","begin.sh"]
