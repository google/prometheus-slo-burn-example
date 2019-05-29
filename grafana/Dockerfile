FROM grafana/grafana:6.0.1

# We do this so we can have sqlite later
USER root
RUN apt update && apt install -y sqlite && rm -rf /var/lib/apt/lists/*

USER grafana
RUN mkdir -p /var/log/grafana /var/lib/grafana
COPY --chown=grafana:grafana grafana.ini /etc/grafana/
COPY --chown=grafana:grafana datasource.yml /etc/grafana/provisioning/datasources/
COPY --chown=grafana:grafana dashboard.yml /etc/grafana/provisioning/dashboards/
COPY --chown=grafana:grafana *.json /var/lib/grafana/dashboards/

# We need to run this command to make sure the database exists. We then modify
# the database to set the default dashboard. We have to do this because we turn
# off user accounts and the API in our settings.
RUN grafana-cli admin reset-admin-password $(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

# This magic incantation to set the preferences for an anonomyous user is
# undocumented, but can be found in code at
# https://github.com/grafana/grafana/blob/v6.0.0/pkg/services/sqlstore/preferences.go#L29
RUN echo "insert into preferences (org_id, user_id, version, home_dashboard_id, team_id, timezone, theme, created, updated) VALUES (1,0,0,1,0,'','', datetime('now'), datetime('now'));" | sqlite3 /var/lib/grafana/grafana.db
