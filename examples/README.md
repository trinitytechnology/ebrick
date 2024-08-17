# eBrick application examples

In this example, we highlight the elasticity of the eBrick framework by showcasing its ability to support seamless transitions between monolithic and microservices architectures. You can begin by building a monolithic application, where all modules are integrated into a single, unified system. As your application grows, the eBrick framework allows you to effortlessly transition to a microservices architecture, where each service contains one or more modules, thanks to its plug-and-play feature.

Moreover, if the need arises, you can just as easily transition back to a monolithic architecture, reintegrating the services into a single application. This flexibility ensures that your application architecture can evolve with your business needs, providing the ability to scale up, break down, or consolidate as necessary—all without extensive rewrites or disruptions.

```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "linkify"}' http://localhost:8080/api/tenants
```
