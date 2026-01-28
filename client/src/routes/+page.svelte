<script>
  const features = [
    {
      title: "Setting Up gRPC-Web",
      description: "Learn how to configure gRPC-Web in your SvelteKit project and set up the necessary build tools.",
      icon: "⚙️"
    },
    {
      title: "Creating Service Clients",
      description: "Generate TypeScript clients from your .proto files and integrate them seamlessly with SvelteKit.",
      icon: "🔌"
    },
    {
      title: "Making RPC Calls",
      description: "Discover how to make unary, server streaming, and bidirectional streaming calls from your Svelte components.",
      icon: "📡"
    },
    {
      title: "Error Handling",
      description: "Implement robust error handling and retry logic for production-ready applications.",
      icon: "🛡️"
    }
  ];
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
  <div class="container mx-auto px-4 py-16">
    <!-- Header -->
    <div class="text-center mb-16">
      <div class="inline-block mb-6">
        <div class="text-6xl mb-4">🚀</div>
      </div>
      <h1 class="text-5xl font-bold text-gray-900 mb-4">
        Welcome to the gRPC-Web Integration
      </h1>
      <h2 class="text-3xl font-semibold text-indigo-600 mb-6">
        with SvelteKit Tutorial
      </h2>
      <p class="text-xl text-gray-600 max-w-2xl mx-auto">
        Learn how to build high-performance, type-safe web applications by connecting your SvelteKit frontend to gRPC services.
      </p>
    </div>

    <!-- CTA -->
    <div class="text-center">
      <a href="/todo" class="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-8 py-3 rounded-lg shadow-md transition-colors duration-300">
        Get Started →
      </a>
      <p class="text-gray-600 mt-4">Ready to dive in? Let's build something amazing!</p>
    </div>
  </div>
</div>
