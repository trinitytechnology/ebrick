# eBrick application examples

In this example, we demonstrate the versatility of the eBrick framework by showcasing its seamless support for transitioning between monolithic and microservices architectures. You can start by building a monolithic application, where all modules are integrated into a single, unified system. As your application grows, the eBrick framework allows you to effortlessly transition to a microservices architecture, where each service contains one or more modules, thanks to its plug-and-play feature.

Furthermore, if the need arises, you can easily switch back to a monolithic architecture, reintegrating the services into a single application. This flexibility ensures that your application architecture can evolve with your business needs, providing the ability to scale up, break down, or consolidate as necessary, without the need for extensive rewrites or disruptions.

## Bootstrap infrastructure

```bash
docker compose -f docker-compose-infra.yml up -d
```

## Bootstrap Observability

```bash
docker compose -f docker-compose-observability.yml up -d
```

## Running the monolithic application

```bash
cd cmd/mono
go run main.go
```

## Running the tenant microservice application

```bash
cd cmd/service-tenant
go run main.go
```

## Running the environment microservice application

```bash
cd cmd/service-env
go run main.go
```

## Testing Tenant creation

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "linkify"}' http://localhost:8080/api/tenants
```

## Testing with K6

```bash
k6 run k6-scripts.js
```
