import { createClient } from '@connectrpc/connect';
import { createGrpcTransport } from '@connectrpc/connect-node'; // Use node for SSR
import { ToDoService } from '$lib/gen/greeter_pb'; // Generated service definition

const protocolLoggingInterceptor = (next) => async (req) => {
    console.log("Request Headers:", req.header);
    const res = await next(req);
    // You can also inspect response headers here
    return res;
};

// This configuration runs once when the app starts (Singleton pattern)
const transport = createGrpcTransport({
	baseUrl: 'http://localhost:50051',
	interceptors: [protocolLoggingInterceptor],
});

// Export the initialized client
export const client = createClient(ToDoService, transport);
