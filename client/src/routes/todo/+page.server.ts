import { client } from '$lib/services/greeter/index';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
  // Full type safety here - no 'as HelloRequest' cast needed
  const response = await client.sayHello({ 
    name: "Svelte 5 Singleton" 
  });

  return {
    greeting: response.message
  };
};