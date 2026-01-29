<script lang="ts">
	import { Status, type ToDoItem } from '$lib/gen/greeter_pb';
	import type { PageProps } from './$types';
	import { enhance } from '$app/forms';

	let { data }: PageProps = $props();

	// Local state for managing todos
	let todos = $derived(data.todos || []);

	let editingId = $state<bigint | null>(null);
	let editingTitle = $state('');
	let filter = $state<'all' | 'active' | 'completed'>('all');

	// Computed filtered todos
	let filteredTodos = $derived(() => {
		switch (filter) {
			case 'active':
				return todos.filter((todo) => todo.item?.status === Status['PENDING']);
			case 'completed':
				return todos.filter((todo) => todo.item?.status === Status['COMPLETED']);
			default:
				return todos;
		}
	});

	// Start editing
	function startEdit(id: bigint, currentTitle: string) {
		editingId = id;
		editingTitle = currentTitle;
	}

	// Save edit
	async function saveEdit(id: bigint) {
		if (!editingTitle.trim()) return;

		// TODO: Make actual gRPC call to backend
		todos = todos.map((todo) =>
			todo.id === BigInt(id) ? { ...todo, item: { ...todo.item, title: editingTitle } } : todo
		) as ToDoItem[];

		editingId = null;
		editingTitle = '';
	}

	// Cancel edit
	function cancelEdit() {
		editingId = null;
		editingTitle = '';
	}

	// Stats
	let activeCount = $derived(todos.filter((t) => t.item?.status === Status['PENDING']).length);
	let completedCount = $derived(todos.filter((t) => t.item?.status === Status['COMPLETED']).length);
</script>

<div class="min-h-screen bg-gradient-to-br from-purple-50 via-pink-50 to-blue-50">
	<div class="container mx-auto max-w-4xl px-4 py-12">
		<!-- Header -->
		<div class="mb-8 text-center">
			<h1
				class="mb-2 bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-5xl font-bold text-transparent"
			>
				My Todo List
			</h1>
			<p class="text-gray-600">Stay organized and get things done</p>
		</div>

		<!-- Stats -->
		<div class="mb-8 grid grid-cols-2 gap-4">
			<div class="rounded-lg bg-white p-4 text-center shadow-md">
				<div class="text-3xl font-bold text-purple-600">{activeCount}</div>
				<div class="text-sm text-gray-600">Active Tasks</div>
			</div>
			<div class="rounded-lg bg-white p-4 text-center shadow-md">
				<div class="text-3xl font-bold text-green-600">{completedCount}</div>
				<div class="text-sm text-gray-600">Completed</div>
			</div>
		</div>

		<!-- Add Todo Form -->
		<div class="mb-6 rounded-lg bg-white p-6 shadow-lg">
			<form method="POST" action="?/create" use:enhance class="flex flex-col gap-4">
				<div class="flex gap-3">
					<input
						type="text"
						required
						name="title"
						placeholder="What needs to be done?"
						class="flex-1 rounded-lg border border-gray-300 px-4 py-3 outline-none focus:ring-2 focus:ring-purple-500"
					/>
					<button
						type="submit"
						class="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 px-6 py-3 font-semibold text-white shadow-md transition-all hover:opacity-90"
					>
						Add Task
					</button>
				</div>

				<div
					class="flex flex-col items-start gap-4 text-sm text-gray-600 md:flex-row md:items-center"
				>
					<textarea
						name="description"
						placeholder="Add a description (optional)..."
						rows="1"
						class="flex-1 resize-none rounded-lg border border-gray-300 px-4 py-2 outline-none focus:ring-2 focus:ring-purple-500"
					></textarea>

					<div class="flex gap-4 rounded-lg border border-gray-200 bg-gray-50 p-2">
						<span class="px-1 font-medium">Status:</span>
						<label class="flex cursor-pointer items-center gap-2">
							<input
								type="radio"
								name="status"
								value={Status.PENDING}
								checked
								class="text-purple-600"
							/>
							Pending
						</label>
						<label class="flex cursor-pointer items-center gap-2">
							<input type="radio" name="status" value={Status.COMPLETED} class="text-purple-600" />
							Completed
						</label>
					</div>
				</div>
			</form>
		</div>

		<!-- Filter Buttons -->
		<div class="mb-6 flex gap-2 rounded-lg bg-white p-2 shadow-md">
			<button
				onclick={() => (filter = 'all')}
				class={`flex-1 rounded-md px-4 py-2 font-medium transition-all ${
					filter === 'all'
						? 'bg-purple-600 text-white shadow-md'
						: 'text-gray-600 hover:bg-gray-100'
				}`}
			>
				All ({todos.length})
			</button>
			<button
				onclick={() => (filter = 'active')}
				class={`flex-1 rounded-md px-4 py-2 font-medium transition-all ${
					filter === 'active'
						? 'bg-purple-600 text-white shadow-md'
						: 'text-gray-600 hover:bg-gray-100'
				}`}
			>
				Active ({activeCount})
			</button>
			<button
				onclick={() => (filter = 'completed')}
				class={`flex-1 rounded-md px-4 py-2 font-medium transition-all ${
					filter === 'completed'
						? 'bg-purple-600 text-white shadow-md'
						: 'text-gray-600 hover:bg-gray-100'
				}`}
			>
				Completed ({completedCount})
			</button>
		</div>

		<!-- Todo List -->
		<div class="overflow-hidden rounded-lg bg-white shadow-lg">
			{#if filteredTodos().length > 0}
				<ul class="divide-y divide-gray-200">
					{#each filteredTodos() as todo (todo.id)}
						<li class="p-4 transition-colors hover:bg-gray-50">
							{#if editingId && BigInt(editingId) === todo.id}
								<!-- Edit Mode -->
								<div class="flex gap-3">
									<input
										type="text"
										bind:value={editingTitle}
										class="flex-1 rounded-lg border border-purple-300 px-3 py-2 focus:ring-2 focus:ring-purple-500 focus:outline-none"
									/>
									<button
										onclick={() => saveEdit(todo.id)}
										class="rounded-lg bg-green-600 px-4 py-2 text-white transition-colors hover:bg-green-700"
									>
										Save
									</button>
									<button
										onclick={cancelEdit}
										class="rounded-lg bg-gray-400 px-4 py-2 text-white transition-colors hover:bg-gray-500"
									>
										Cancel
									</button>
								</div>
							{:else}
								<!-- View Mode -->
								<div class="flex items-center gap-3">
									<form
										method="POST"
										action="?/updateStatus"
										use:enhance={() => {
											return async ({ update }) => {
												await update();
											};
										}}
									>
										<input type="hidden" name="id" value={todo.id.toString()} />
										<input
											type="hidden"
											name="status"
											value={todo.item?.status === Status.COMPLETED
												? Status.PENDING
												: Status.COMPLETED}
										/>
										<input
											type="checkbox"
											checked={todo.item?.status === Status['COMPLETED']}
											onchange={(e) => e.currentTarget.form?.requestSubmit()}
											class="h-5 w-5 cursor-pointer rounded border-gray-300 text-purple-600 focus:ring-purple-500"
										/>
									</form>
									<span
										class={`flex-1 text-lg ${
											todo.item?.status === Status['COMPLETED']
												? 'text-gray-400 line-through'
												: 'text-gray-800'
										}`}
									>
										{todo.item?.title}
									</span>
									<div class="flex gap-2">
										<button
											onclick={() => startEdit(todo.id, todo.item?.title || '')}
											class="rounded-lg p-2 text-blue-600 transition-colors hover:bg-blue-50"
											title="Edit"
										>
											<svg
												xmlns="http://www.w3.org/2000/svg"
												class="h-5 w-5"
												viewBox="0 0 20 20"
												fill="currentColor"
											>
												<path
													d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z"
												/>
											</svg>
										</button>
										<form
											method="POST"
											action="?/delete"
											use:enhance={() => {
												return async ({ update }) => {
													await update();
												};
											}}
										>
											<input type="hidden" name="id" value={todo.id.toString()} />
											<button
												type="submit"
												class="rounded-lg p-2 text-red-600 transition-colors hover:bg-red-50"
												title="Delete"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="h-5 w-5"
													viewBox="0 0 20 20"
													fill="currentColor"
												>
													<path
														fill-rule="evenodd"
														d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z"
														clip-rule="evenodd"
													/>
												</svg>
											</button>
										</form>
									</div>
								</div>
							{/if}
						</li>
					{/each}
				</ul>
			{:else}
				<div class="p-12 text-center">
					<div class="mb-4 text-6xl">📝</div>
					<p class="text-lg text-gray-500">
						{filter === 'all'
							? 'No todos yet. Add one to get started!'
							: filter === 'active'
								? 'No active tasks. Great job!'
								: 'No completed tasks yet. Keep going!'}
					</p>
				</div>
			{/if}
		</div>

		{#if todos.length > 0 && completedCount > 0}
			<div class="mt-6 text-center">
				<button
					onclick={() => (todos = todos.filter((t) => t.item?.status === Status['COMPLETED']))}
					class="rounded-lg px-6 py-2 font-medium text-red-600 transition-colors hover:bg-red-50"
				>
					Clear Completed Tasks
				</button>
			</div>
		{/if}
	</div>
</div>
