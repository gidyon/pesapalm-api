docker run -it --rm \
    -e KONG_ENDPOINT=http://kong-admin:8001 \
    -e EMAIL=gideonhacer@gmail.com \
    -e DOMAINS=pesapalm.com \
    phpdockerio/kong-certbot-agent

curl -i -X PATCH http://localhost:8001/services/pesapalm_service \
  --data url='http://nginx_pesapalm'


curl -i -X POST http://localhost:8001/services/KongCertbot/routes \
  --data "strip_path=false" \
  --data "preserve_host=false" \
  --data "regex_priority=0" \
  --data "paths[]=/well-known/acme-challenge" \
  --data "hosts[]=your.list" \
  --data "hosts[]=of.domains" \
  --data "hosts[]=for.the" \
  --data "hosts[]=same.certificate" \
  --data "methods[]=GET" \
  --data "protocols[]=http"

KONG_ADMIN_URL=http://localhost:8001
DOMAINS=pesapalm.com
EMAIL=gideonhacer@gmail.com
certbot run --domains $DOMAINS --email $EMAIL \
  -a certbot-kong:kong -i certbot-kong:kong \
  --certbot-kong:kong-admin-url $KONG_ADMIN_URL