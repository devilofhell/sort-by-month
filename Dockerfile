FROM archlinux:latest
WORKDIR /opt
COPY sort-by-month /opt/sort-by-month
CMD [ "mkdir data" ]
ENV WATCH="/opt/data/"
ENV RUNNING_INTERVAL="10s"
ENTRYPOINT ["/opt/sort-by-month"]