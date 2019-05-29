FROM cloudprober/cloudprober:latest

ENV PORT 8080
ENV CLOUDPROBER_PORT $PORT
EXPOSE $PORT

COPY ./probe.cfg /etc/cloudprober.cfg
RUN chmod a+r /etc/cloudprober.cfg
