FROM prom/prometheus

COPY --chown=nobody:nogroup *.yml /etc/prometheus/

# Tests your rules before starting
RUN promtool check rules /etc/prometheus/*.rules.yml
