<script lang="ts">
	import { onMount } from 'svelte';
	import { user } from './../store';
	import { goto } from '$app/navigation';
	let userData: any;
	let msg: string;
	let messages: any = [];
	$: messages;
	const unsubscribe = user.subscribe((value) => {
		userData = value;
	});
	let roomName = '';
	let responseRoom = null;
	let roomID: number;
	let rooms: any = [];

	let socket: WebSocket;
	const fetchRooms = async () => {
		try {
			const response = await fetch('http://localhost:8080/ws/getRooms', {
				method: 'GET',
				headers: { 'Content-Type': 'application/json' }
			});

			if (response.ok) {
				rooms = await response.json();
				console.log(rooms);
			}
		} catch (err) {
			console.log(err);
			rooms = [];
		}
	};

	function generateRandomString(length: number) {
		const characters = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
		let result = '';
		for (let i = 0; i < length; i++) {
			result += characters.charAt(Math.floor(Math.random() * characters.length));
		}
		return result;
	}

	const createRoom = async () => {
		try {
			const res = await fetch('http://localhost:8080/ws/createRoom', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({
					id: generateRandomString(10),
					name: roomName
				})
			});

			if (res.ok) {
				responseRoom = await res.json();
				console.log(responseRoom);
				fetchRooms();
			}
		} catch (err) {
			console.log(err);
		}
	};
	const checkAuth = () => {
		if (!userData) {
			goto('/login');
		}
	};
	onMount(() => {
		fetchRooms();
		checkAuth();
	});
	const sendMessage = () => {
		socket.send(msg);
		msg = '';
	};
	function joinRoom(ID: number) {
		roomID = ID;
		socket = new WebSocket(
			`ws://localhost:8080/ws/${ID}?userId=${userData.id}&username=${userData.name}`
		);

		socket.onopen = (event) => {
			console.log('opened');
		};
		socket.onmessage = (event) => {
			const message = JSON.parse(event.data);
			messages = [...messages, message];
			console.log(messages);
		};

		socket.onclose = (event) => {
			console.log('Socket closed:', event);
		};
	}
</script>

<div class=" px-8 py-4 w-full h-screen-5">
	<p>Welcome, {userData ? userData.name : 'Guest'} (ID: {userData ? userData.id : 'No ID'})!</p>
	<div class="flex justify-center gap-4">
		<input type="text" bind:value={roomName} class="input input-bordered" placeholder="room name" />
		<button class=" btn btn-neutral" on:click={createRoom}> create room </button>
	</div>
	<div class="flex h-full gap-20 pt-5">
		<div class=" w-1/4">
			<div class="font-bold">Rooms</div>
			<div class=" flex flex-col gap-4 pt-6">
				{#each rooms as room}
					<div class="p-4 flex border-sky-700 border-2 items-center rounded-md w-full">
						<div class="w-full">
							<div class="text-sm">room</div>
							<div class="text-blue font-bold text-lg">{room.name}</div>
						</div>
						<div class="">
							<button
								class="px-4 btn btn-ghost text-white bg-blue rounded-md"
								on:click={() => joinRoom(room.id)}>join</button
							>
						</div>
					</div>
				{/each}
			</div>
		</div>
		<div class="w-full h-full flex flex-col justify-between">
			<div>
				{#each messages as m}
					{#if m.username != userData.name}
						<div class="chat chat-start">
							<div class="chat-header">{m.username}</div>
							<div class="chat-bubble">{m.content}</div>
						</div>
					{:else}
						<div class="chat chat-end">
							<div class="chat-header">{m.username}</div>
							<div class="chat-bubble">{m.content}</div>
						</div>
					{/if}
				{/each}
			</div>

			<div class="flex justify-center gap-x-2 mt-3 p-5">
				<input bind:value={msg} type="text" class="input w-3/4" placeholder="Type your message" />
				<button on:click={sendMessage} class=" btn"> Send </button>
			</div>
		</div>
	</div>
</div>

<style>
	.h-screen-5 {
		height: 95vh;
	}
</style>
