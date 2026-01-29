<script lang="ts">
	import { Status, type ToDoDetails, type ToDoItem } from '$lib/gen/greeter_pb';
	import { timestampNow } from '@bufbuild/protobuf/wkt';
    import type { PageProps } from './$types';
    import { enhance } from '$app/forms';
    
    let { data, form }: PageProps = $props();
    
    // Local state for managing todos
    let todos = $state(data.todos || []);

    $effect(() => {
        todos = data.todos || [];
    });
    
    let newTodoTitle = $state('');
    let editingId = $state<bigint | null>(null);
    let editingTitle = $state('');
    let filter = $state<'all' | 'active' | 'completed'>('all');

    // Computed filtered todos
    let filteredTodos = $derived(() => {
        switch (filter) {
            case 'active':
                return todos.filter(todo => todo.item?.status === Status['PENDING']);
            case 'completed':
                return todos.filter(todo => todo.item?.status === Status['COMPLETED']);
            default:
                return todos;
        }
    });

    // Add new todo
    async function addTodo() {
        if (!newTodoTitle.trim()) return;
        
        // TODO: Make actual gRPC call to backend
        const newTodo: ToDoItem = {
            id: BigInt(1),
            item: {
                title: newTodoTitle,
                status: Status['STATUS_UNSPECIFIED'],
                description: "New description"
            } as ToDoDetails,
            createdAt: timestampNow()
        } as ToDoItem;
        
        todos = [...todos, newTodo];
        newTodoTitle = '';
    }

    // Toggle todo completion
    async function toggleTodo(id: bigint) {
        // TODO: Make actual gRPC call to backend
        todos = todos.map(todo => 
            todo.id === BigInt(id) 
                ? { ...todo, item: { ...todo.item, completed: !todo.item?.status } }
                : todo
        ) as ToDoItem[];
    }

    // Start editing
    function startEdit(id: bigint, currentTitle: string) {
        editingId = id;
        editingTitle = currentTitle;
    }

    // Save edit
    async function saveEdit(id: bigint) {
        if (!editingTitle.trim()) return;
        
        // TODO: Make actual gRPC call to backend
        todos = todos.map(todo =>
            todo.id === BigInt(id)
                ? { ...todo, item: { ...todo.item, title: editingTitle } }
                : todo
        ) as ToDoItem[];
        
        editingId = null;
        editingTitle = '';
    }

    // Cancel edit
    function cancelEdit() {
        editingId = null;
        editingTitle = '';
    }

    // Delete todo
    async function deleteTodo(id: bigint) {
        // TODO: Make actual gRPC call to backend
        todos = todos.filter(todo => todo.id !== BigInt(id));
    }

    // Stats
    let activeCount = $derived(todos.filter(t => t.item?.status === Status['PENDING']).length);
    let completedCount = $derived(todos.filter(t => t.item?.status === Status['COMPLETED']).length);
</script>

<div class="min-h-screen bg-gradient-to-br from-purple-50 via-pink-50 to-blue-50">
    <div class="container px-4 py-12 mx-auto max-w-4xl">
        <!-- Header -->
        <div class="text-center mb-8">
            <h1 class="text-5xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-600 to-pink-600 mb-2">
                My Todo List
            </h1>
            <p class="text-gray-600">Stay organized and get things done</p>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-2 gap-4 mb-8">
            <div class="bg-white rounded-lg shadow-md p-4 text-center">
                <div class="text-3xl font-bold text-purple-600">{activeCount}</div>
                <div class="text-sm text-gray-600">Active Tasks</div>
            </div>
            <div class="bg-white rounded-lg shadow-md p-4 text-center">
                <div class="text-3xl font-bold text-green-600">{completedCount}</div>
                <div class="text-sm text-gray-600">Completed</div>
            </div>
        </div>

        <!-- Add Todo Form -->
        <div class="bg-white rounded-lg shadow-lg p-6 mb-6">
    <form method="POST" action="?/create" use:enhance class="flex flex-col gap-4">
        <div class="flex gap-3">
            <input
                type="text"
                required
                name="title"
                placeholder="What needs to be done?"
                class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none"
            />
            <button
                type="submit"
                class="px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white font-semibold rounded-lg hover:opacity-90 transition-all shadow-md"
            >
                Add Task
            </button>
        </div>

        <div class="flex flex-col md:flex-row gap-4 items-start md:items-center text-sm text-gray-600">
            <textarea
                name="description"
                placeholder="Add a description (optional)..."
                rows="1"
                class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 outline-none resize-none"
            ></textarea>

            <div class="flex gap-4 p-2 bg-gray-50 rounded-lg border border-gray-200">
                <span class="font-medium px-1">Status:</span>
                <label class="flex items-center gap-2 cursor-pointer">
                    <input type="radio" name="status" value={Status.PENDING} checked class="text-purple-600" />
                    Pending
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                    <input type="radio" name="status" value={Status.COMPLETED} class="text-purple-600" />
                    Completed
                </label>
            </div>
        </div>
    </form>
</div>

        <!-- Filter Buttons -->
        <div class="flex gap-2 mb-6 bg-white rounded-lg shadow-md p-2">
            <button
                onclick={() => filter = 'all'}
                class={`flex-1 px-4 py-2 rounded-md font-medium transition-all ${
                    filter === 'all'
                        ? 'bg-purple-600 text-white shadow-md'
                        : 'text-gray-600 hover:bg-gray-100'
                }`}
            >
                All ({todos.length})
            </button>
            <button
                onclick={() => filter = 'active'}
                class={`flex-1 px-4 py-2 rounded-md font-medium transition-all ${
                    filter === 'active'
                        ? 'bg-purple-600 text-white shadow-md'
                        : 'text-gray-600 hover:bg-gray-100'
                }`}
            >
                Active ({activeCount})
            </button>
            <button
                onclick={() => filter = 'completed'}
                class={`flex-1 px-4 py-2 rounded-md font-medium transition-all ${
                    filter === 'completed'
                        ? 'bg-purple-600 text-white shadow-md'
                        : 'text-gray-600 hover:bg-gray-100'
                }`}
            >
                Completed ({completedCount})
            </button>
        </div>

        <!-- Todo List -->
        <div class="bg-white rounded-lg shadow-lg overflow-hidden">
            {#if filteredTodos().length > 0}
                <ul class="divide-y divide-gray-200">
                    {#each filteredTodos() as todo (todo.id)}
                        <li class="p-4 hover:bg-gray-50 transition-colors">
                            {#if editingId && BigInt(editingId) === todo.id}
                                <!-- Edit Mode -->
                                <div class="flex gap-3">
                                    <input
                                        type="text"
                                        bind:value={editingTitle}
                                        class="flex-1 px-3 py-2 border border-purple-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                                    />
                                    <button
                                        onclick={() => saveEdit(todo.id)}
                                        class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
                                    >
                                        Save
                                    </button>
                                    <button
                                        onclick={cancelEdit}
                                        class="px-4 py-2 bg-gray-400 text-white rounded-lg hover:bg-gray-500 transition-colors"
                                    >
                                        Cancel
                                    </button>
                                </div>
                            {:else}
                                <!-- View Mode -->
                                <div class="flex items-center gap-3">
                                    <input
                                        type="checkbox"
                                        checked={todo.item?.status === Status['COMPLETED']}
                                        onchange={() => toggleTodo(todo.id)}
                                        class="w-5 h-5 text-purple-600 border-gray-300 rounded focus:ring-purple-500 cursor-pointer"
                                    />
                                    <span
                                        class={`flex-1 text-lg ${
                                            todo.item?.status === Status['COMPLETED']
                                                ? 'line-through text-gray-400'
                                                : 'text-gray-800'
                                        }`}
                                    >
                                        {todo.item?.title}
                                    </span>
                                    <div class="flex gap-2">
                                        <button
                                            onclick={() => startEdit(todo.id, todo.item?.title || '')}
                                            class="p-2 text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                                            title="Edit"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                                            </svg>
                                        </button>
                                        <button
                                            onclick={() => deleteTodo(todo.id)}
                                            class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                                            title="Delete"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                                                <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                                            </svg>
                                        </button>
                                    </div>
                                </div>
                            {/if}
                        </li>
                    {/each}
                </ul>
            {:else}
                <div class="p-12 text-center">
                    <div class="text-6xl mb-4">📝</div>
                    <p class="text-gray-500 text-lg">
                        {filter === 'all' ? 'No todos yet. Add one to get started!' : 
                         filter === 'active' ? 'No active tasks. Great job!' :
                         'No completed tasks yet. Keep going!'}
                    </p>
                </div>
            {/if}
        </div>

        {#if todos.length > 0 && completedCount > 0}
            <div class="mt-6 text-center">
                <button
                    onclick={() => todos = todos.filter(t => t.item?.status === Status['COMPLETED'])}
                    class="px-6 py-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors font-medium"
                >
                    Clear Completed Tasks
                </button>
            </div>
        {/if}
    </div>
</div>