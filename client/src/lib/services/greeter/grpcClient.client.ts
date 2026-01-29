import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ToDoService } from "$lib/gen/greeter_pb";

const transport = createConnectTransport({
  baseUrl: "http://localhost:50051",
  useBinaryFormat: true // TEST ME
});

export const client = createClient(ToDoService, transport);