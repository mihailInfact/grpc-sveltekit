import { client } from '$lib/services/greeter/index';
import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { CreateRequestSchema, Status } from '$lib/gen/greeter_pb';

export const load: PageServerLoad = async () => {
  const response = await client.getAll({});

  return {
    // SvelteKit "sees" the type of response.items here 
    // and passes it to the PageProps interface.
    todos: response.items 
  };
};

export const actions: Actions = {

  create: async ({ request }) => {
    const data = await request.formData();

    const title = data.get('title')?.toString();
    const description = data.get('description')?.toString() || "";
    const statusRaw = data.get('status')?.toString();

    const status: Status = statusRaw ? parseInt(statusRaw) : Status['STATUS_UNSPECIFIED'];

    if (!title) {
      return fail(400, { message: 'Title is required' });
    }

    try {
      await client.create({ 
        item: {
          title,
          description,
          status
        }
      }); 
      return { success: true };
    } catch (e) {
      console.error('gRPC Error:', e);
      return fail(500, { message: 'Failed to communicate with Go backend' });
    }
  }
}