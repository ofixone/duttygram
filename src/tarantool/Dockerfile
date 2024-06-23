FROM tarantool/tarantool:2.11

COPY . /opt/tarantool

RUN mkdir -p /var/log/tarantool/ && \
   chown -R tarantool:tarantool /var/log/tarantool