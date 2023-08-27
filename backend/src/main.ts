import { Application, Router } from "https://deno.land/x/oak@v11.1.0/mod.ts";
import { oakCors } from "https://deno.land/x/cors@v1.2.2/mod.ts";
import { Client } from "https://deno.land/x/mysql@v2.12.0/mod.ts";

const client = await new Client().connect({
  hostname: "mysql",
  username: "mysql",
  password: "mysql",
  db: "app",
});

const router = new Router();
router
  .get("/", (context) => {
    context.response.body = "Welcome to dinosaur API!";
    console.log("root");
  })
  .get("/mysql", async (context) => {
    const users = await client.execute("SELECT * FROM user");
    context.response.body = JSON.stringify(users);
  });

const app = new Application();
app.use(oakCors()); // Enable CORS for All Routes
app.use(router.routes());
app.use(router.allowedMethods());

await app.listen({ port: 8000 });
