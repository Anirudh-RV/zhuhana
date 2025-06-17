import Fastify from "fastify";
import fastifyStatic from "@fastify/static";
import { fileURLToPath } from "url";
import path from "path";

// ESM-friendly __dirname
const fastify = Fastify();

fastify.register(fastifyStatic, {
  root: path.join(process.cwd(), "frontend", "build"),
});

fastify.setNotFoundHandler((request, reply) => {
  reply.sendFile("index.html");
});

fastify.listen({ port: 3000, host: "0.0.0.0" }).then(() => {
  console.log("✅ Server running on http://localhost:3000");
});
