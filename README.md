Go Microservices architecture for building SaaS

membership (Accessed by user-token and backoffice-token)
 - organization (sql)
    - user
    - workspace
      - subscription
 - backoffice (sql)
needs workspace id to be provided in the request

cashier (Accessed by user-token and backoffice-token)
 - payment_method  (nosql)
 - payment_intent (nosql)
 - receipt (nosql)
 - invoice (nosql)
needs subscription id to be provided in the request

notifier (Accessed by backoffice-token)
 - email
 - sms

service (Accessed by user-token and backoffice-token and apikey-token)
 - plan (nosql)
 - activity (nosql)
 - archive (nosql)
 - ... service specific resources


/v1/identity/<resource>
/v1/membership/<workspace-id>/subscriptions
/v1/cashier/<workspace-id>/<resource>

/v1/service/<workspace-id>/<resource>
service will lookup the workspace id from identity and retreive apikey and match it with the provided.
if valid it will lookup the subscription id from membership 
if there is subscription to the service and it is active it will allow the request
and cache key: apikey , value: subscription id. the cache path is <service_name>.apikey.<apikey>


