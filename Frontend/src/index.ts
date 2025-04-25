import Fastify from "fastify";

const fastify = Fastify();

fastify.get("/", async (request, reply) => {
  return "Hello there! 👋";
});

const start = async () => {
  try {
    await fastify.listen({ port: 3000, host: "0.0.0.0" });
    console.log("Server is running on http://0.0.0.0:3000");
  } catch (err) {
    console.error(err);
    process.exit(1);
  }
};

start();
