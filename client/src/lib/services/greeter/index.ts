import { createClient } from "@connectrpc/connect";
import { createGrpcTransport } from "@connectrpc/connect-node"; // Use node for SSR
import { Greeter } from "$lib/gen/greeter_pb"; // Generated service definition

// This configuration runs once when the app starts (Singleton pattern)
const transport = createGrpcTransport({
  baseUrl: "http://localhost:50051",
});

// Export the initialized client
export const client = createClient(Greeter, transport);